// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package struct_markdown

import (
	"fmt"
	"os"
	"strings"
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
			[]string{"../test-data/packer-plugin-happycloud/builder/happycloud/config.go"},
			0,
			FileCheck{
				Expected: []string{
					"../test-data/packer-plugin-happycloud/docs-partials/builder/happycloud/Config-not-required.mdx",
					"../test-data/packer-plugin-happycloud/docs-partials/builder/happycloud/Config.mdx",
					"../test-data/packer-plugin-happycloud/docs-partials/builder/happycloud/Config-required.mdx",
					"../test-data/packer-plugin-happycloud/docs-partials/builder/happycloud/CustomerEncryptionKey-not-required.mdx",
					"../test-data/packer-plugin-happycloud/docs-partials/builder/happycloud/CustomerEncryptionKey.mdx",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.args), func(t *testing.T) {
			defer tt.FileCheck.Cleanup(t)
			cmd := &Command{}
			if got := cmd.Run(tt.args); got != tt.want {
				t.Errorf("CMD.Run() = %v, want %v", got, tt.want)
			}
			targetedPath := strings.TrimPrefix(tt.args[0], "../test-data/packer-plugin-happycloud/")
			for _, p := range tt.FileCheck.ExpectedFiles() {
				raw, _ := os.ReadFile(p)
				content := string(raw)
				if !strings.Contains(content, targetedPath) {
					t.Errorf("%s must contain '%s'. Its content is:\n%s", p, targetedPath, content)
				}
			}
		})
	}

}
