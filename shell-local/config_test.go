// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package shell_local

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToLinuxPath(t *testing.T) {
	winPath := "C:/path/to/your/file"
	winBashPath := "/mnt/c/path/to/your/file"
	converted, _ := ConvertToLinuxPath(winPath)
	assert.Equal(t, winBashPath, converted,
		"Should have converted %s to %s -- not %s", winPath, winBashPath, converted)

}

func TestValidateDefaultExecuteCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix execute command behavior")
	}

	tests := []struct {
		name           string
		configure      func(*Config)
		executeCommand []string
	}{
		{
			name: "inline invokes its shebang interpreter",
			configure: func(config *Config) {
				config.Inline = []string{"echo foo"}
			},
			executeCommand: []string{"/bin/sh", "-c", `{{.Vars}} /bin/sh -e "$0"`, "{{.Script}}"},
		},
		{
			name: "command invokes its shebang interpreter",
			configure: func(config *Config) {
				config.Command = "echo foo"
			},
			executeCommand: []string{"/bin/sh", "-c", `{{.Vars}} /bin/sh -e "$0"`, "{{.Script}}"},
		},
		{
			name: "custom inline shebang is preserved",
			configure: func(config *Config) {
				config.Inline = []string{"echo foo"}
				config.InlineShebang = "/bin/bash -eu"
			},
			executeCommand: []string{"/bin/sh", "-c", `{{.Vars}} /bin/bash -eu "$0"`, "{{.Script}}"},
		},
		{
			name: "script retains direct execution",
			configure: func(config *Config) {
				config.Script = os.DevNull
			},
			executeCommand: []string{"/bin/sh", "-c", "{{.Vars}} {{.Script}}"},
		},
		{
			name: "custom execute command is preserved",
			configure: func(config *Config) {
				config.Inline = []string{"echo foo"}
				config.ExecuteCommand = []string{"custom", "{{.Script}}"}
			},
			executeCommand: []string{"custom", "{{.Script}}"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := &Config{}
			test.configure(config)

			if err := Validate(config); err != nil {
				t.Fatalf("validating config: %s", err)
			}

			assert.Equal(t, test.executeCommand, config.ExecuteCommand)
		})
	}
}
