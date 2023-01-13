// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/zclconf/go-cty/cty"
)

type testDatasource struct {
	configCalled bool
	configVal    []interface{}

	outputSpecCalled bool
	outputSpec       hcldec.ObjectSpec

	executeCalled bool
	executeValue  cty.Value
}

func (*testDatasource) ConfigSpec() hcldec.ObjectSpec { return nil }

func (d *testDatasource) Configure(configs ...interface{}) error {
	d.configCalled = true
	d.configVal = configs
	return nil
}

func (d *testDatasource) OutputSpec() hcldec.ObjectSpec {
	d.outputSpecCalled = true
	return d.outputSpec
}

func (d *testDatasource) Execute() (cty.Value, error) {
	d.executeCalled = true
	return d.executeValue, nil
}

func TestDatasource(t *testing.T) {
	d := new(testDatasource)
	client, server := testClientServer(t)
	defer client.Close()
	defer server.Close()
	server.RegisterDatasource(d)
	dsClient := client.Datasource()

	// Test Configure
	config := 42
	if err := dsClient.Configure(config); err != nil {
		t.Fatalf("error: %s", err)
	}
	if !d.configCalled {
		t.Fatal("config should be called")
	}
	expected := []interface{}{int64(42)}
	if !reflect.DeepEqual(d.configVal, expected) {
		t.Fatalf("unknown config value: %#v", d.configVal)
	}

	// Test OutPutSpec
	d.outputSpec = map[string]hcldec.Spec{
		"foo": &hcldec.AttrSpec{Name: "foo", Type: cty.String, Required: false},
	}
	spec := dsClient.OutputSpec()
	if !reflect.DeepEqual(spec, d.outputSpec) {
		t.Fatalf("unknown outputSpec value: %#v", spec)
	}

	// Test Execute
	d.executeValue = cty.StringVal("foo")
	val, err := dsClient.Execute()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if val.Equals(d.executeValue).False() {
		t.Fatalf("unknown value: %#v", val)
	}
}

func TestDatasource_Implements(t *testing.T) {
	var raw interface{} = new(datasource)
	if _, ok := raw.(packer.Datasource); !ok {
		t.Fatal("not a datasource")
	}
}
