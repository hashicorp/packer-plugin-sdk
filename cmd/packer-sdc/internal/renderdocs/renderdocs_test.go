// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package renderdocs

import (
	"fmt"
	"path/filepath"
	"testing"

	. "github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc/internal/cmd"
)

func Test_RenderDocsFolder(t *testing.T) {
	tests := []struct {
		cmd              []string
		outputFolderHash FileCheck
		wantErr          bool
	}{
		{
			[]string{
				"-src", filepath.Join("test-data/docs"),
				"-partials", filepath.Join("test-data/docs-partials"),
				"-dst", filepath.Join("test-data/.docs"),
			},
			FileCheck{
				ExpectedContent: map[string]string{
					"test-data/.docs/builder-docs.mdx": `Hello and welcome to the awesome docs


foo


Bar:

bar


<!-- Code generated from the comments of the Config struct in communicator/config.go; DO NOT EDIT MANUALLY -->

Config is the common configuration a builder uses to define and configure a Packer
communicator. Embed this struct in your builder config to implement
communicator support.

<!-- End of code generated from the comments of the Config struct in communicator/config.go; -->


End of file
`,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.cmd), func(t *testing.T) {
			defer tt.outputFolderHash.Cleanup(t)
			cmd := Command{}
			if err := cmd.run(tt.cmd); (err != nil) != tt.wantErr {
				t.Errorf("renderDocsFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			tt.outputFolderHash.Verify(t, "")
		})
	}
}
