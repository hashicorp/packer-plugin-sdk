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
			name:       "no go-cty module dependency",
			fixtureDir: filepath.Join("testdata", "norequire"),
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
			var dryRun bool
			testFixer := NewGoCtyFixer()
			testFixtureDir := tc.fixtureDir
			expectedFn := filepath.Join(testFixtureDir, "go.mod")
			expected, err := os.ReadFile(expectedFn)
			if err != nil {
				t.Fatalf("failed while reading text fixture: %s", err)
			}

			outFileDir := copyToTempFile(t, expected)
			if err := testFixer.Fix(outFileDir, dryRun); err != nil {
				t.Fatalf("expected dryrun check to not err but it did: %v", err)
			}

			outFileFn := filepath.Join(outFileDir, "go.mod")
			fixed, err := os.ReadFile(outFileFn)
			if err != nil {
				t.Fatalf("failed while reading text fixture: %s", err)
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
			versionStr: "1.13.1",
			fixtureDir: filepath.Join("testdata", "unfixed", "version"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var dryRun bool
			testFixer := NewGoCtyFixer()
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
			if err := testFixer.Fix(outFileDir, dryRun); err != nil {
				t.Fatalf("expected dryrun check to not err but it did: %v", err)
			}

			outFileFn := filepath.Join(outFileDir, "go.mod")
			fixed, err := os.ReadFile(outFileFn)
			if err != nil {
				t.Fatalf("failed while reading text fixture: %s", err)
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
	var dryRun bool
	testFixer := NewGoCtyFixer()
	testFixtureDir := filepath.Join("testdata", "invalid")
	expectedFn := filepath.Join(testFixtureDir, "go.mod")
	expected, err := os.ReadFile(expectedFn)
	if err != nil {
		t.Fatalf("failed while reading text fixture: %s", err)
	}

	outFileDir := copyToTempFile(t, expected)
	if err := testFixer.Fix(outFileDir, dryRun); err == nil {
		t.Fatalf("expected dryrun check to err but it didn't: %v", err)
	}
}
