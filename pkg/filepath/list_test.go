// Package filepath helps deal with configuration paths that come from user
// config files. It supports a 'slash-only' way of interpretting paths, and
// allows transforming those paths into a system agnostic path. See
// https://github.com/hashicorp/packer/issues/6188 for concrete examples.

package filepath

import (
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestListFromString(t *testing.T) {

	tests := []struct {
		path string
		want List
	}{
		{"c:/path/to/file.txt",
			List{"c:", "path", "to", "file.txt"}},
		{"~/go/src/github.com/hashicorp/packer-plugin-sdk/cmd",
			List{"~", "go", "src", "github.com", "hashicorp", "packer-plugin-sdk", "cmd"}},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := ListFromString(tt.path)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ListFromString() --got ++want %s", diff)
			}
		})
	}
}

func TestList_String(t *testing.T) {
	tests := []struct {
		list              *List
		wantUnixString    string
		wantWindowsString string
	}{
		{
			&List{"~", "go", "src", "github.com", "hashicorp", "packer-plugin-sdk", "cmd"},
			`~/go/src/github.com/hashicorp/packer-plugin-sdk/cmd`,
			`~\go\src\github.com\hashicorp\packer-plugin-sdk\cmd`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.list.String(), func(t *testing.T) {

			unixString := tt.list.UnixString()
			if unixString != tt.wantUnixString {
				t.Fatalf("bad unix string: expected: %s, got %s", tt.wantUnixString, unixString)
			}

			windowsString := tt.list.WindowsString()
			if windowsString != tt.wantWindowsString {
				t.Fatalf("bad windows string: expected: %s, got %s", tt.wantWindowsString, windowsString)
			}

			var expectedString string
			gotString := tt.list.String()
			if runtime.GOOS == "windows" {
				expectedString = tt.list.WindowsString()
			} else {
				expectedString = tt.list.UnixString()
			}
			if expectedString != gotString {
				t.Fatalf("On %q, bad string: expected: %s, got %s", runtime.GOOS, expectedString, gotString)
			}
		})
	}
}
