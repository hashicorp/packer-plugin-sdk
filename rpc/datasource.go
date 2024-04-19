// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
	"reflect"
	"unsafe"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/set"
)

// An implementation of packer.Datasource where the data source is actually
// executed over an RPC connection.
type datasource struct {
	commonClient
}

type DatasourceConfigureArgs struct {
	Configs []interface{}
}

type DatasourceConfigureResponse struct {
	Error *BasicError
}

func (d *datasource) Configure(configs ...interface{}) error {
	configs, err := encodeCTYValues(configs)
	if err != nil {
		return err
	}
	var resp DatasourceConfigureResponse
	if err := d.client.Call(d.endpoint+".Configure", &DatasourceConfigureArgs{Configs: configs}, &resp); err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
	}
	return err
}

type OutputSpecResponse struct {
	OutputSpec []byte
}

func (d *datasource) OutputSpec() hcldec.ObjectSpec {
	resp := new(OutputSpecResponse)
	if err := d.client.Call(d.endpoint+".OutputSpec", new(interface{}), resp); err != nil {
		err := fmt.Errorf("Datasource.OutputSpec failed: %v", err)
		panic(err.Error())
	}
	res := hcldec.ObjectSpec{}
	err := gob.NewDecoder(bytes.NewReader(resp.OutputSpec)).Decode(&res)
	if err != nil {
		panic("ici:" + err.Error())
	}
	return res
}

type ExecuteResponse struct {
	Value []byte
	Error *BasicError
}

func (d *datasource) Execute() (cty.Value, error) {
	res := new(cty.Value)
	resp := new(ExecuteResponse)
	if err := d.client.Call(d.endpoint+".Execute", new(interface{}), resp); err != nil {
		err := fmt.Errorf("Datasource.Execute failed: %v", err)
		return *res, err
	}
	if err := gob.NewDecoder(bytes.NewReader(resp.Value)).Decode(&res); err != nil {
		return *res, err
	}
	if err := gobDecodeFixNumberPtrVal(res); err != nil {
		return *res, err
	}
	if err := resp.Error; err != nil {
		return *res, err
	}
	return *res, nil
}

// DatasourceServer wraps a packer.Datasource implementation and makes it
// exportable as part of a Golang RPC server.
type DatasourceServer struct {
	contextCancel func()

	commonServer
	d packer.Datasource
}

func (d *DatasourceServer) Configure(args *DatasourceConfigureArgs, reply *DatasourceConfigureResponse) error {
	config, err := decodeCTYValues(args.Configs)
	if err != nil {
		return err
	}
	err = d.d.Configure(config...)
	reply.Error = NewBasicError(err)
	return err
}

func (d *DatasourceServer) OutputSpec(args *DatasourceConfigureArgs, reply *OutputSpecResponse) error {
	spec := d.d.OutputSpec()
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(spec)
	reply.OutputSpec = b.Bytes()
	return err
}

func (d *DatasourceServer) Execute(args *interface{}, reply *ExecuteResponse) error {
	spec, err := d.d.Execute()
	reply.Error = NewBasicError(err)
	b := bytes.NewBuffer(nil)
	err = gob.NewEncoder(b).Encode(spec)
	reply.Value = b.Bytes()
	if reply.Error != nil {
		err = reply.Error
	}
	return err
}

func (d *DatasourceServer) Cancel(args *interface{}, reply *interface{}) error {
	if d.contextCancel != nil {
		d.contextCancel()
	}
	return nil
}

func init() {
	gob.Register(new(cty.Value))
}

// goDecodeFixNumberPtr fixes an unfortunate quirk of round-tripping cty.Number
// values through gob: the big.Float.GobEncode method is implemented on a
// pointer receiver, and so it loses the "pointer-ness" of the value on
// encode, causing the values to emerge the other end as big.Float rather than
// *big.Float as we expect elsewhere in cty.
//
// The implementation of gobDecodeFixNumberPtr mutates the given raw value
// during its work, and may either return the same value mutated or a new
// value. Callers must no longer use whatever value they pass as "raw" after
// this function is called.
func gobDecodeFixNumberPtr(raw interface{}, ty cty.Type) interface{} {
	// Unfortunately we need to work recursively here because number values
	// might be embedded in structural or collection type values.

	switch {
	case ty.Equals(cty.Number):
		if bf, ok := raw.(big.Float); ok {
			return &bf // wrap in pointer
		}
	case ty.IsMapType():
		if m, ok := raw.(map[string]interface{}); ok {
			for k, v := range m {
				m[k] = gobDecodeFixNumberPtr(v, ty.ElementType())
			}
		}
	case ty.IsListType():
		if s, ok := raw.([]interface{}); ok {
			for i, v := range s {
				s[i] = gobDecodeFixNumberPtr(v, ty.ElementType())
			}
		}
	case ty.IsSetType():
		if s, ok := raw.(set.Set); ok {
			newS := set.NewSet(s.Rules())
			for it := s.Iterator(); it.Next(); {
				newV := gobDecodeFixNumberPtr(it.Value(), ty.ElementType())
				newS.Add(newV)
			}
			return newS
		}
	case ty.IsObjectType():
		if m, ok := raw.(map[string]interface{}); ok {
			for k, v := range m {
				aty := ty.AttributeType(k)
				m[k] = gobDecodeFixNumberPtr(v, aty)
			}
		}
	case ty.IsTupleType():
		if s, ok := raw.([]interface{}); ok {
			for i, v := range s {
				ety := ty.TupleElementType(i)
				s[i] = gobDecodeFixNumberPtr(v, ety)
			}
		}
	}

	return raw
}

// gobDecodeFixNumberPtrVal is a helper wrapper around gobDecodeFixNumberPtr
// that works with already-constructed values. This is primarily for testing,
// to fix up intentionally-invalid number values for the parts of the test
// code that need them to be valid, such as calling GoString on them.
func gobDecodeFixNumberPtrVal(val *cty.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	v := (*interface{})(unsafe.Pointer(reflect.Indirect(reflect.ValueOf(val)).FieldByName("v").UnsafeAddr()))
	*v = gobDecodeFixNumberPtr(*v, val.Type())

	return nil
}
