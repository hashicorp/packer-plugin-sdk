// +build windows

package pathing

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPath(t *testing.T) {
	// temporary directories for env vars
	packerConfigTempDir, err := ioutil.TempDir(os.TempDir(), "*")
	if err != nil {
		t.Fatalf("Failed to create temp test directory: failing test: %v", err)
	}
	defer os.RemoveAll(packerConfigTempDir)

	// reset env
	packerConfigDir := os.Getenv("PACKER_CONFIG_DIR")
	os.Setenv("PACKER_CONFIG_DIR", "")
	defer func() {
		os.Setenv("PACKER_CONFIG_DIR", packerConfigDir)
	}()

	homedir, err := homeDir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	tests := []struct {
		name    string
		env     map[string]string
		want    string
		wantErr bool
	}{
		{
			"base",
			nil,
			filepath.Join(homedir, defaultConfigDir),
			false,
		},
		{
			"env PACKER_CONFIG_DIR",
			map[string]string{"PACKER_CONFIG_DIR": packerConfigTempDir},
			filepath.Join(packerConfigTempDir, defaultConfigDir),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			got, err := ConfigDir()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Name: %v, ConfigPath() error = %v, wantErr %v",
					tt.name,
					err,
					tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf(
					"Name: %v, ConfigPath() = %v, want %v",
					tt.name,
					got,
					tt.want)
			}
			resetTestEnv()
		})
	}
}

func resetTestEnv() {
	os.Setenv("PACKER_CONFIG_DIR", "")
}
