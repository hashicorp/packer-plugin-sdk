package rpc

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hcldec"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

var b testing.B

func TestCommonServer_ConfigSpec(t *testing.T) {
	tt := []struct {
		name      string
		component packersdk.HCL2Speccer
	}{
		{
			name:      "Builder Component Server",
			component: new(packersdk.MockBuilder),
		},
		{
			name:      "Datasource Component Server",
			component: new(packersdk.MockDatasource),
		},
		{
			name:      "Provisioner Component Server",
			component: new(packersdk.MockProvisioner),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Start the server
			client, server := testClientServer(t)
			defer client.Close()
			defer server.Close()

			var spec hcldec.ObjectSpec
			switch v := tc.component.(type) {
			case packersdk.Builder:
				server.RegisterBuilder(v)
				remote := client.Builder()
				spec = remote.ConfigSpec()
			case packersdk.Datasource:
				server.RegisterDatasource(v)
				remote := client.Datasource()
				spec = remote.ConfigSpec()
			case packersdk.Provisioner:
				server.RegisterProvisioner(v)
				remote := client.Provisioner()
				spec = remote.ConfigSpec()
			case packersdk.PostProcessor:
				server.RegisterPostProcessor(v)
				remote := client.PostProcessor()
				spec = remote.ConfigSpec()
			default:
				t.Fatalf("Unknown component type %T", v)
			}

			if len(spec) == 0 {
				t.Errorf("expected remote.ConfigSpec for %T to return a valid hcldec.ObjectSpec, but return %v", tc.component, spec)
			}

		})
	}

}
