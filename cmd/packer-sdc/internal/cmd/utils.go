package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type FileCheck struct {
	Expected, NotExpected []string
	ExpectedContent       map[string]string
}

func (fc FileCheck) cleanup(t *testing.T) {
	for _, file := range fc.expectedFiles() {
		t.Logf("removing %v", file)
		if err := os.Remove(file); err != nil {
			t.Errorf("failed to remove file %s: %v", file, err)
		}
	}
}

func (fc FileCheck) expectedFiles() []string {
	expected := fc.Expected
	for file := range fc.ExpectedContent {
		expected = append(expected, file)
	}
	return expected
}

func (fc FileCheck) verify(t *testing.T, dir string) {
	for _, f := range fc.expectedFiles() {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("Expected to find %s: %v", f, err)
		}
	}
	for _, f := range fc.NotExpected {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			t.Errorf("Expected to not find %s", f)
		}
	}
	for file, expectedContent := range fc.ExpectedContent {
		content, err := ioutil.ReadFile(filepath.Join(dir, file))
		if err != nil {
			t.Fatalf("ioutil.ReadFile: %v", err)
		}
		if diff := cmp.Diff(expectedContent, string(content)); diff != "" {
			t.Errorf("content of %s differs: %s", file, diff)
		}
	}
}
