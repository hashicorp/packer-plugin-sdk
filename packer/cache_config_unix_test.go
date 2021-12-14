//go:build darwin || freebsd || linux || netbsd || openbsd || solaris
// +build darwin freebsd linux netbsd openbsd solaris

package packer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCachePath(t *testing.T) {
	// temporary directories for env vars
	xdgCacheHomeTempDir, err := ioutil.TempDir(os.TempDir(), "*")
	if err != nil {
		t.Fatalf("Failed to create temp test directory: failing test: %v", err)
	}
	defer os.RemoveAll(xdgCacheHomeTempDir)
	packerCacheTempDir, err := ioutil.TempDir(os.TempDir(), "*")
	if err != nil {
		t.Fatalf("Failed to create temp test directory: failing test: %v", err)
	}
	defer os.RemoveAll(packerCacheTempDir)

	// reset env
	packerCacheDir := os.Getenv("PACKER_CACHE_DIR")
	os.Setenv("PACKER_CACHE_DIR", "")
	defer func() {
		os.Setenv("PACKER_CACHE_DIR", packerCacheDir)
	}()

	xdgCacheHomeDir := os.Getenv("XDG_CACHE_HOME")
	os.Setenv("XDG_CACHE_HOME", "")
	defer func() {
		os.Setenv("XDG_CACHE_HOME", xdgCacheHomeDir)
	}()

	type args struct {
		paths []string
	}
	tests := []struct {
		name    string
		args    args
		env     map[string]string
		want    string
		wantErr bool
	}{
		{
			"base",
			args{},
			nil,
			filepath.Join(os.Getenv("HOME"), ".cache", "packer"),
			false,
		},
		{
			"base and path",
			args{[]string{"a", "b"}},
			nil,
			filepath.Join(os.Getenv("HOME"), ".cache", "packer", "a", "b"),
			false,
		},
		{
			"env PACKER_CACHE_DIR and path",
			args{[]string{"a", "b"}},
			map[string]string{"PACKER_CACHE_DIR": packerCacheTempDir},
			filepath.Join(packerCacheTempDir, "a", "b"),
			false,
		},
		{
			"env XDG_CACHE_HOME and path",
			args{[]string{"a", "b"}},
			map[string]string{"XDG_CACHE_HOME": xdgCacheHomeTempDir},
			filepath.Join(xdgCacheHomeTempDir, "packer", "a", "b"),
			false,
		},
		{
			"env PACKER_CACHE_DIR, XDG_CACHE_HOME, and path",
			args{[]string{"a", "b"}},
			map[string]string{
				"XDG_CACHE_HOME":   xdgCacheHomeTempDir,
				"PACKER_CACHE_DIR": packerCacheTempDir,
			},
			filepath.Join(packerCacheTempDir, "a", "b"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			got, err := CachePath(tt.args.paths...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CachePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CachePath() = %v, want %v", got, tt.want)
			}
			resetTestEnv()
		})
	}
}

func resetTestEnv() {
	os.Setenv("PACKER_CACHE_DIR", "")
	os.Setenv("XDG_CACHE_HOME", "")
}
