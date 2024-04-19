// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// mapstructure-to-hcl2 fills the gaps between hcl2 and mapstructure for Packer
//
// By generating a struct that the HCL2 ecosystem understands making use of
// mapstructure tags.
//
// Packer heavily uses the mapstructure decoding library to load/parse user
// config files. Packer now needs to move to HCL2.
//
// Here are a few differences/gaps betweens hcl2 and mapstructure:
//
//   - in HCL2 all basic struct fields (string/int/struct) that are not pointers
//     are required ( must be set ). In mapstructure everything is optional.
//
//   - mapstructure allows to 'squash' fields
//     (ex: Field CommonStructType `mapstructure:",squash"`) this allows to
//     decorate structs and reuse configuration code. HCL2 parsing libs don't have
//     anything similar.
//
// mapstructure-to-hcl2 will parse Packer's config files and generate the HCL2
// compliant code that will allow to not change any of the current builders in
// order to softly move to HCL2.
package mapstructure_to_hcl2

import (
	"fmt"
	"testing"

	. "github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/cmd"
)

func TestCommand_Run(t *testing.T) {

	tests := []struct {
		args []string
		want int
		FileCheck
	}{
		{
			[]string{"-type", "Config,CustomerEncryptionKey", "../test-data/packer-plugin-happycloud/builder/happycloud/config.go"},
			0,
			FileCheck{
				Expected: []string{"../test-data/packer-plugin-happycloud/builder/happycloud/config.hcl2spec.go"},
			},
		},
		{
			[]string{"-type", "Config", "../test-data/field-conflict/test_mapstructure_field_conflict.go"},
			1,
			FileCheck{},
		},
		{
			[]string{"-type", "Config", "../test-data/tag-conflict/test_mapstructure_tag_conflict.go"},
			1,
			FileCheck{},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.args), func(t *testing.T) {
			// remove files before actually generating them; because our ci
			// generates files all the time. This fails if the expected files
			// are not present
			tt.FileCheck.Cleanup(t)
			cmd := &Command{}
			if got := cmd.Run(tt.args); got != tt.want {
				t.Errorf("CMD.Run() = %v, want %v", got, tt.want)
			}
			tt.FileCheck.Verify(t, ".")
		})
	}
}
