package renderdocs

import (
	"fmt"
	"path/filepath"
	"testing"
)

func Test_RenderDocsFolder(t *testing.T) {
	tests := []struct {
		cmd              []string
		outputFolderHash map[string]string
		wantErr          bool
	}{
		{
			[]string{
				"-src", filepath.Join("test-data/docs"),
				"-partials", filepath.Join("test-data/docs-partials"),
				"-dst", filepath.Join("test-data/.docs"),
			},
			map[string]string{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.cmd), func(t *testing.T) {
			cmd := Command{}
			if err := cmd.run(tt.cmd); (err != nil) != tt.wantErr {
				t.Errorf("renderDocsFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
