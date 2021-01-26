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
	set.SetVersion(pluginVersion.InitializePluginVersion(
		"1.1.1", ""))

	outputDesc := set.description()

	sdkVersion := pluginVersion.InitializePluginVersion(pluginVersion.Version, pluginVersion.VersionPrerelease)
	if diff := cmp.Diff(SetDescription{
		Version:        "1.1.1",
		SDKVersion:     sdkVersion.String(),
		APIVersion:     "x" + APIVersionMajor + "_" + APIVersionMinor,
		Builders:       []string{"example", "example-2"},
		PostProcessors: []string{"example", "example-2"},
		Provisioners:   []string{"example", "example-2"},
		Datasources:    []string{"example", "example-2"},
	}, outputDesc); diff != "" {
		t.Fatalf("Unexpected description: %s", diff)
	}

	err := set.RunCommand("start", "builder", "example")
	if diff := cmp.Diff(err.Error(), ErrManuallyStartedPlugin.Error()); diff != "" {
		t.Fatalf("Unexpected error: %s", diff)
	}
}
