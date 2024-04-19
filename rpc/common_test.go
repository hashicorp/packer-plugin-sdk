// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hcldec"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

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

			var configSpecTestFn func() hcldec.ObjectSpec
			switch v := tc.component.(type) {
			case packersdk.Builder:
				configSpecTestFn = func() hcldec.ObjectSpec {
					server.RegisterBuilder(v)
					remote := client.Builder()
					spec := remote.ConfigSpec()
					return spec
				}
			case packersdk.Datasource:
				configSpecTestFn = func() hcldec.ObjectSpec {
					server.RegisterDatasource(v)
					remote := client.Datasource()
					spec := remote.ConfigSpec()
					return spec
				}
			case packersdk.Provisioner:
				configSpecTestFn = func() hcldec.ObjectSpec {
					server.RegisterProvisioner(v)
					remote := client.Provisioner()
					spec := remote.ConfigSpec()
					return spec
				}
			case packersdk.PostProcessor:
				configSpecTestFn = func() hcldec.ObjectSpec {
					server.RegisterPostProcessor(v)
					remote := client.PostProcessor()
					spec := remote.ConfigSpec()
					return spec
				}
			default:
				t.Fatalf("Unknown component type %T", v)
			}

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Call to ConfigSpec for %s panicked: %v", tc.name, r)
				}
			}()

			spec := configSpecTestFn()
			if len(spec) == 0 {
				t.Errorf("expected remote.ConfigSpec for %T to return a valid hcldec.ObjectSpec, but return %v", tc.component, spec)
			}

		})
	}

}
