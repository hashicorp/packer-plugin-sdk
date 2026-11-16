// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/zclconf/go-cty/cty"
)

type testDatasource struct {
	configCalled bool
	configVal    []any

	outputSpecCalled bool
	outputSpec       hcldec.ObjectSpec

	executeCalled bool
	executeValue  cty.Value
}

type testSensitiveDatasource struct {
	testDatasource

	sensitiveOutputs []string
}

func (d *testSensitiveDatasource) SensitiveOutputs() []string {
	return d.sensitiveOutputs
}

// legacyDatasourceServer stands in for a plugin built with an SDK that
// predates the SensitiveOutputs method.
type legacyDatasourceServer struct{}

func (*legacyDatasourceServer) Cancel(args *any, reply *any) error {
	return nil
}

// failingDatasourceServer serves SensitiveOutputs but fails the call.
type failingDatasourceServer struct{}

func (*failingDatasourceServer) SensitiveOutputs(args *any, reply *SensitiveOutputsResponse) error {
	return errors.New("metadata unavailable")
}

func (*testDatasource) ConfigSpec() hcldec.ObjectSpec { return nil }

func (d *testDatasource) Configure(configs ...any) error {
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
	expected := []any{int64(42)}
	if !reflect.DeepEqual(d.configVal, expected) {
		t.Fatalf("unknown config value: %#v", d.configVal)
	}

	// Test OutPutSpec
	d.outputSpec = map[string]hcldec.Spec{
		"foo": &hcldec.AttrSpec{Name: "foo", Type: cty.String, Required: false},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Call to ConfigSpec for datasource panicked: %v", r)
		}
	}()
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

func TestDatasource_SensitiveOutputs(t *testing.T) {
	tests := []struct {
		name string
		ds   packer.Datasource
		want []string
	}{
		{
			name: "implemented",
			ds:   &testSensitiveDatasource{sensitiveOutputs: []string{"result", "token"}},
			want: []string{"result", "token"},
		},
		{
			name: "implemented with no outputs",
			ds:   &testSensitiveDatasource{},
		},
		{
			name: "not implemented",
			ds:   new(testDatasource),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := testClientServer(t)
			defer client.Close()
			defer server.Close()
			server.RegisterDatasource(tt.ds)

			got := client.Datasource().(packer.SensitiveOutputer).SensitiveOutputs()
			if len(got) != len(tt.want) || (len(got) > 0 && !reflect.DeepEqual(got, tt.want)) {
				t.Fatalf("got sensitive outputs %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDatasource_SensitiveOutputsLegacyPlugin(t *testing.T) {
	client, server := testClientServer(t)
	defer client.Close()
	defer server.Close()
	if err := server.server.RegisterName(DefaultDatasourceEndpoint, new(legacyDatasourceServer)); err != nil {
		t.Fatalf("register legacy datasource server: %s", err)
	}

	got := client.Datasource().(packer.SensitiveOutputer).SensitiveOutputs()
	if got != nil {
		t.Fatalf("got sensitive outputs %#v from a legacy plugin, want nil", got)
	}
}

func TestDatasource_SensitiveOutputsUnexpectedError(t *testing.T) {
	client, server := testClientServer(t)
	defer client.Close()
	defer server.Close()
	if err := server.server.RegisterName(DefaultDatasourceEndpoint, new(failingDatasourceServer)); err != nil {
		t.Fatalf("register failing datasource server: %s", err)
	}

	defer func() {
		if recover() == nil {
			t.Fatal("expected SensitiveOutputs to panic on an unexpected RPC error")
		}
	}()
	got := client.Datasource().(packer.SensitiveOutputer).SensitiveOutputs()
	t.Fatalf("got sensitive outputs %#v, want a panic", got)
}

func TestDatasource_Implements(t *testing.T) {
	var raw any = new(datasource)
	if _, ok := raw.(packer.Datasource); !ok {
		t.Fatal("not a datasource")
	}
	if _, ok := raw.(packer.SensitiveOutputer); !ok {
		t.Fatal("not a sensitive outputer")
	}
}
