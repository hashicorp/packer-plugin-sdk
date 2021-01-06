package rpc

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/zclconf/go-cty/cty"
)

// An implementation of packer.Datasource where the data source is actually
// executed over an RPC connection.
type datasource struct {
	commonClient
}

type DatasourceConfigureArgs struct {
	Configs []interface{}
}

func (d *datasource) Configure(configs ...interface{}) error {
	configs, err := encodeCTYValues(configs)
	if err != nil {
		return err
	}
	args := &DatasourceConfigureArgs{Configs: configs}
	return d.client.Call(d.endpoint+".Configure", args, new(interface{}))
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
	err := gob.NewDecoder(bytes.NewReader(resp.Value)).Decode(&res)
	if err != nil {
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

func (d *DatasourceServer) Configure(args *DatasourceConfigureArgs, reply *interface{}) error {
	config, err := decodeCTYValues(args.Configs)
	if err != nil {
		return err
	}
	return d.d.Configure(config...)
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
