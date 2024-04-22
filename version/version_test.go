// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package version

import (
	"fmt"
	"testing"
)

func Test_PluginVersionCreate(t *testing.T) {
	tests := []struct {
		name                string
		coreVersion         string
		preVersion          string
		metaVersion         string
		expectError         bool
		expectVersionString string
	}{
		{
			"Valid semver core only version",
			"1.0.0",
			"",
			"",
			false,
			"1.0.0",
		},
		{
			"Valid semver, should get canonical version",
			"01.001.001",
			"",
			"",
			false,
			"1.1.1",
		},
		{
			"Valid semver with prerelease, should get canonical version",
			"1.001.010",
			"dev",
			"",
			false,
			"1.1.10-dev",
		},
		{
			"Valid semver with metadata, should get canonical version",
			"1.001.010",
			"",
			"123abcdef",
			false,
			"1.1.10+123abcdef",
		},
		{
			"Valid semver with prerelease and metadata, should get canonical version",
			"1.001.010",
			"dev",
			"123abcdef",
			false,
			"1.1.10-dev+123abcdef",
		},
		{
			"Invalid version, should fail",
			".1.1",
			"",
			"",
			true,
			"",
		},
		{
			"4-parts version, should not be accepted",
			"1.1.1.1",
			"",
			"",
			true,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				panicMsg := recover()
				if !tt.expectError && panicMsg != nil {
					t.Errorf("creating version panicked, should not have.")
				}

				if tt.expectError && panicMsg == nil {
					t.Errorf("creating version should have panicked, but did not.")
				}

				if panicMsg != nil {
					t.Logf("panic message was: %v", panicMsg)
				}
			}()

			ver := NewPluginVersion(tt.coreVersion, tt.preVersion, tt.metaVersion)
			verStr := ver.String()
			if verStr != tt.expectVersionString {
				t.Errorf("string format mismatch, version created is %q, expected %q", verStr, tt.expectVersionString)
			}
		})
	}
}

func TestFormattedVersionString(t *testing.T) {
	GitCommit = "abcdef12345"
	defer func() {
		GitCommit = ""
	}()

	expectedVersion := fmt.Sprintf("1.0.0-dev (%s)", GitCommit)

	ver := InitializePluginVersion("1.0.0", "dev")
	formatted := ver.FormattedVersion()
	if formatted != expectedVersion {
		t.Errorf("Expected formatted version %q; got %q", expectedVersion, formatted)
	}

	expectedVersion = fmt.Sprintf("1.0.0-dev+meta (%s)", GitCommit)
	ver = NewPluginVersion("1.0.0", "dev", "meta")
	formatted = ver.FormattedVersion()
	if formatted != expectedVersion {
		t.Errorf("Expected formatted version %q; got %q", expectedVersion, formatted)
	}
}
