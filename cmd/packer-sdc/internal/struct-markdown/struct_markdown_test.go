package struct_markdown

import (
	"fmt"
	"testing"
)

func TestCommand_Run(t *testing.T) {

	type args struct {
		args []string
	}
	tests := []struct {
		args []string
		want int
		FileCheck
	}{
		{
			[]string{"-type", "Config,CustomerEncryptionKey", "test-data/packer-plugin-google/builder/happycloud/config.go"},
			0,
			FileCheck{
				Expected: []string{"test-data/packer-plugin-google/builder/happycloud/config.hcl2spec.go"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.args), func(t *testing.T) {
			defer tt.FileCheck.Cleanup(t)
			cmd := &CMD{}
			if got := cmd.Run(tt.args); got != tt.want {
				t.Errorf("CMD.Run() = %v, want %v", got, tt.want)
			}
			tt.FileCheck.Verify(t, ".")
		})
	}

}
