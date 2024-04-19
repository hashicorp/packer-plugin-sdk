// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fix

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func copyToTempFile(t testing.TB, data []byte) string {
	t.Helper()
	dir := t.TempDir()
	fn := filepath.Join(dir, "go.mod")
	err := os.WriteFile(fn, data, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test mod file: %s", err)
	}
	return dir
}

func goctyTestFixer() goCtyFix {
	return goCtyFix{
		OldPath:    oldPath,
		NewPath:    newPath,
		NewVersion: newVersion,
	}
}

func TestFixGoCty_FixNotNeeded(t *testing.T) {
	tt := []struct {
		name       string
		fixtureDir string
	}{
		{
			name:       "empty mod file",
			fixtureDir: filepath.Join("testdata", "empty"),
		},
		{
			name:       "no requires for go-cty or packer-plugin-sdk modules",
			fixtureDir: filepath.Join("testdata", "missing-requires", "both"),
		},
		{
			name:       "no go-cty module dependency",
			fixtureDir: filepath.Join("testdata", "missing-requires", "go-cty"),
		},
		{
			name:       "no packer-plugin-sdk module dependency",
			fixtureDir: filepath.Join("testdata", "missing-requires", "packer-plugin-sdk"),
		},
		{
			name:       "previously fixed mod file",
			fixtureDir: filepath.Join("testdata", "fixed", "basic"),
		},
		{
			name:       "fixed mod file with other replace directives",
			fixtureDir: filepath.Join("testdata", "fixed", "many-replace"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testFixer := goctyTestFixer()
			testFixtureDir := tc.fixtureDir
			expectedFn := filepath.Join(testFixtureDir, modFilename)
			expected, err := os.ReadFile(expectedFn)
			if err != nil {
				t.Fatalf("failed while reading text fixture: %s", err)
			}

			outFileDir := copyToTempFile(t, expected)
			outFileFn := filepath.Join(outFileDir, "go.mod")
			fixed, err := testFixer.fix(outFileFn, expected)
			if err != nil {
				t.Fatalf("expected fix to not err but it did: %v", err)
			}

			if diff := cmp.Diff(expected, fixed); diff != "" {
				t.Errorf("expected no differences but got %q", diff)
			}

		})
	}
}

func TestFixGoCty_Unfixed(t *testing.T) {
	tt := []struct {
		name       string
		versionStr string
		fixtureDir string
	}{
		{
			name:       "basic unfixed mod file",
			fixtureDir: filepath.Join("testdata", "unfixed", "basic"),
		},
		{
			name:       "unfixed mod file with other replace directives",
			fixtureDir: filepath.Join("testdata", "unfixed", "many-replace"),
		},
		{
			name:       "out of date fix",
			fixtureDir: filepath.Join("testdata", "unfixed", "version"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testFixer := goctyTestFixer()
			if tc.versionStr != "" {
				testFixer.NewVersion = tc.versionStr
			}
			testFixtureDir := tc.fixtureDir
			unfixedFn := filepath.Join(testFixtureDir, "go.mod")
			unfixed, err := os.ReadFile(unfixedFn)
			if err != nil {
				t.Fatalf("failed while reading text fixture: %s", err)
			}

			outFileDir := copyToTempFile(t, unfixed)
			outFileFn := filepath.Join(outFileDir, modFilename)
			fixed, err := testFixer.fix(outFileFn, unfixed)
			if err != nil {
				t.Fatalf("expected fix to not err but it did: %v", err)
			}

			expectedFn := filepath.Join(testFixtureDir, "fixed.go.mod")
			expected, err := os.ReadFile(expectedFn)
			if err != nil {
				t.Fatalf("failed while reading text fixture: %s", err)
			}

			if diff := cmp.Diff(expected, fixed); diff != "" {
				t.Errorf("expected differences but got %q", diff)
			}

		})
	}
}

func TestFixGoCty_InvalidReplacePath(t *testing.T) {
	testFixer := goctyTestFixer()
	testFixtureDir := filepath.Join("testdata", "invalid")
	expectedFn := filepath.Join(testFixtureDir, modFilename)
	expected, err := os.ReadFile(expectedFn)
	if err != nil {
		t.Fatalf("failed while reading text fixture: %s", err)
	}

	outFileDir := copyToTempFile(t, expected)
	outFileFn := filepath.Join(outFileDir, modFilename)
	if _, err := testFixer.fix(outFileFn, expected); err == nil {
		t.Fatalf("expected fix to err but it didn't: %v", err)
	}
}
