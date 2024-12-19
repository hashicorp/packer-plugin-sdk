// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	pluginVersion "github.com/hashicorp/packer-plugin-sdk/version"
)

type MockBuilder struct {
	packersdk.Builder
}

var _ packersdk.Builder = new(MockBuilder)

type MockProvisioner struct {
	packersdk.Provisioner
}

var _ packersdk.Provisioner = new(MockProvisioner)

type MockPostProcessor struct {
	packersdk.PostProcessor
}

var _ packersdk.PostProcessor = new(MockPostProcessor)

type MockDatasource struct {
	packersdk.Datasource
}

var _ packersdk.Datasource = new(MockDatasource)

func TestSet(t *testing.T) {
	set := NewSet()
	set.RegisterBuilder("example-2", new(MockBuilder))
	set.RegisterBuilder("example", new(MockBuilder))
	set.RegisterPostProcessor("example", new(MockPostProcessor))
	set.RegisterPostProcessor("example-2", new(MockPostProcessor))
	set.RegisterProvisioner("example", new(MockProvisioner))
	set.RegisterProvisioner("example-2", new(MockProvisioner))
	set.RegisterDatasource("example", new(MockDatasource))
	set.RegisterDatasource("example-2", new(MockDatasource))
	set.SetVersion(pluginVersion.NewPluginVersion(
		"1.1.1", "", ""))

	outputDesc := set.description()

	sdkVersion := pluginVersion.NewPluginVersion(pluginVersion.Version, pluginVersion.VersionPrerelease, "")
	if diff := cmp.Diff(SetDescription{
		Version:         "1.1.1",
		SDKVersion:      sdkVersion.String(),
		APIVersion:      "x" + APIVersionMajor + "." + APIVersionMinor,
		Builders:        []string{"example", "example-2"},
		PostProcessors:  []string{"example", "example-2"},
		Provisioners:    []string{"example", "example-2"},
		Datasources:     []string{"example", "example-2"},
		ProtocolVersion: ProtocolVersion2,
	}, outputDesc); diff != "" {
		t.Fatalf("Unexpected description: %s", diff)
	}

	err := set.RunCommand("start", "builder", "example")
	if diff := cmp.Diff(err.Error(), ErrManuallyStartedPlugin.Error()); diff != "" {
		t.Fatalf("Unexpected error: %s", diff)
	}
}

func TestSetProtobufArgParsing(t *testing.T) {
	testCases := []struct {
		name     string
		useProto bool
		in, out  []string
	}{
		{
			name:     "no --protobuf argument provided",
			in:       []string{"start", "builder", "example"},
			out:      []string{"start", "builder", "example"},
			useProto: false,
		},
		{
			name:     "providing --protobuf as first argument",
			in:       []string{"--protobuf", "start", "builder", "example"},
			out:      []string{"start", "builder", "example"},
			useProto: true,
		},
		{
			name:     "providing --protobuf as last argument",
			in:       []string{"start", "builder", "example", "--protobuf"},
			out:      []string{"start", "builder", "example"},
			useProto: true,
		},
		{
			name:     "providing --protobuf as middle argument",
			in:       []string{"start", "builder", "--protobuf", "example"},
			out:      []string{"start", "builder", "example"},
			useProto: true,
		},
		{
			name:     "providing --protobuf multiple times",
			in:       []string{"--protobuf", "start", "builder", "--protobuf", "example", "--protobuf"},
			out:      []string{"start", "builder", "example"},
			useProto: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := NewSet()
			got := set.parseProtobufFlag(tc.in...)

			if diff := cmp.Diff(got, tc.out); diff != "" {
				t.Errorf("Unexpected args: %s", diff)
			}

			if set.useProto != tc.useProto {
				t.Errorf("expected useProto to be %t when %s but got %t", tc.useProto, tc.name, set.useProto)
			}
		})

	}
}
