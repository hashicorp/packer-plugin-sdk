package packer

import (
	"os"
	"path/filepath"
)

// CachePath returns an absolute path to a cache file or directory
//
// When the directory is not absolute, CachePath will try to make a
// a cache depending on the operating system.
//
// TODO: Update this with better examples
// NOTE: cache directory will change depending on operating system dependent
// For Windows:
//   PACKER_CACHE_DIR=""            CacheDir() => "./packer_cache/
//   PACKER_CACHE_DIR=""            CacheDir("foo") => "./packer_cache/foo
//   PACKER_CACHE_DIR="bar"         CacheDir("foo") => "./bar/foo
//   PACKER_CACHE_DIR="/home/there" CacheDir("foo", "bar") => "/home/there/foo/bar
// For Unix:
//   PACKER_CACHE_DIR="",            XDG_CONFIG_HOME="", Default_config CacheDir() => "$HOME/cache/packer"
//   PACKER_CACHE_DIR="",            XDG_CONFIG_HOME="", Default_config CacheDir("foo") => "$HOME/cache/packer/foo"
//   PACKER_CACHE_DIR="bar",         XDG_CONFIG_HOME="", Default_config CacheDir("foo") => "./bar/foo
//   PACKER_CACHE_DIR="/home/there", XDG_CONFIG_HOME="", Default_config CacheDir("foo", "bar") => "/home/there/foo/bar
func CachePath(paths ...string) (path string, err error) {
	defer func() {
		// create the dir based on return path if it doesn't exist
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}()
	cacheDir := getDefaultCacheDir()
	if cd := os.Getenv("PACKER_CACHE_DIR"); cd != "" {
		cacheDir = cd
	}

	paths = append([]string{cacheDir}, paths...)
	return filepath.Abs(filepath.Join(paths...))
}
