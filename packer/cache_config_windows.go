// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package packer

const (
	defaultConfigFile = "packer_cache"
)

func getDefaultCacheDir() string {
	return defaultConfigFile
}
