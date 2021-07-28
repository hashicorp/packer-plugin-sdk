// +build darwin freebsd linux netbsd openbsd solaris

package pathing

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPath(t *testing.T) {
	// temporary directories for env vars
	xdgConfigHomeTempDir, err := ioutil.TempDir(os.TempDir(), "*")
	if err != nil {
		t.Fatalf("Failed to create temp test directory: failing test: %v", err)
	}
	defer os.RemoveAll(xdgConfigHomeTempDir)

	packerConfigTempDir, err := ioutil.TempDir(os.TempDir(), "*")
	if err != nil {
		t.Fatalf("Failed to create temp test directory: failing test: %v", err)
	}
	defer os.RemoveAll(packerConfigTempDir)

	homeTempDir, err := ioutil.TempDir(os.TempDir(), "*")
	if err != nil {
		t.Fatalf("Failed to create temp test directory: failing test: %v", err)
	}
	defer os.RemoveAll(homeTempDir)

	homeDirDefaultConfigTempDir, err := ioutil.TempDir(os.TempDir(), "*")
	if err != nil {
		t.Fatalf("Failed to create temp test directory: failing test: %v", err)
	}

	err = os.Mkdir(filepath.Join(homeDirDefaultConfigTempDir, defaultConfigDir), 0755)
	if err != nil {
		t.Fatalf("Failed to create temp test file: failing test: %v", err)
	}
	defer os.RemoveAll(homeDirDefaultConfigTempDir)

	// reset env
	packerConfigDir := os.Getenv("PACKER_CONFIG_DIR")
	os.Setenv("PACKER_CONFIG_DIR", "")
	defer func() {
		os.Setenv("PACKER_CONFIG_DIR", packerConfigDir)
	}()

	xdgConfigHomeDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", "")
	defer func() {
		os.Setenv("XDG_CONFIG_HOME", xdgConfigHomeDir)
	}()

	homeDir := os.Getenv("HOME")
	os.Setenv("HOME", "")
	defer func() {
		os.Setenv("HOME", homeDir)
	}()

	tests := []struct {
		name    string
		env     map[string]string
		want    string
		wantErr bool
	}{
		{
			"no HOME env var",
			nil,
			"",
			true,
		},
		{
			"base",
			map[string]string{"HOME": homeTempDir},
			filepath.Join(homeTempDir, ".config", "packer"),
			false,
		},
		{
			"XDG_CONFIG_HOME set without default file",
			map[string]string{
				"XDG_CONFIG_HOME": xdgConfigHomeTempDir,
				"HOME":            homeTempDir,
			},
			filepath.Join(xdgConfigHomeTempDir, "packer"),
			false,
		},
		{
			"env PACKER_CONFIG_DIR",
			map[string]string{"PACKER_CONFIG_DIR": packerConfigTempDir},
			filepath.Join(packerConfigTempDir, defaultConfigDir),
			false,
		},
		{
			"env PACKER_CONFIG_DIR, XDG_CONFIG_HOME",
			map[string]string{
				"XDG_CONFIG_HOME":   xdgConfigHomeTempDir,
				"PACKER_CONFIG_DIR": packerConfigTempDir,
			},
			filepath.Join(packerConfigTempDir, defaultConfigDir),
			false,
		},
		{
			"Old Default Config Found",
			map[string]string{"HOME": homeDirDefaultConfigTempDir},
			filepath.Join(homeDirDefaultConfigTempDir, defaultConfigDir),
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
	os.Setenv("XDG_CONFIG_HOME", "")
	os.Setenv("HOME", "")
}
