// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonsteps

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	urlhelper "github.com/hashicorp/go-getter/v2/helper/url"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/tmp"
)

var _ multistep.Step = new(StepDownload)

func toSha1(in string) string {
	b := sha1.Sum([]byte(in))
	return hex.EncodeToString(b[:])
}

func abs(t *testing.T, path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	u, err := urlhelper.Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	return u.String()
}

func TestStepDownload_Run(t *testing.T) {
	srvr := httptest.NewServer(http.FileServer(http.Dir("test-fixtures")))
	defer srvr.Close()

	cs := map[string]string{
		"/root/basic.txt":   "f572d396fae9206628714fb2ce00f72e94f2258f",
		"/root/another.txt": "7c6e5dd1bacb3b48fdffba2ed096097eb172497d",
	}

	type fields struct {
		Checksum    string
		Description string
		ResultKey   string
		TargetPath  string
		Url         []string
		Extension   string
	}

	tests := []struct {
		name      string
		fields    fields
		want      multistep.StepAction
		wantFiles []string
	}{
		{"Empty URL field passes",
			fields{Url: []string{}},
			multistep.ActionContinue,
			nil,
		},
		{"not passing a checksum passes",
			fields{Url: []string{abs(t, "./test-fixtures/root/another.txt")}},
			multistep.ActionContinue,
			[]string{
				toSha1(abs(t, "./test-fixtures/root/another.txt")) + ".lock",
			},
		},
		{"double slashes on a local filesystem passes",
			fields{Url: []string{abs(t, "./test-fixtures/root//another.txt")}},
			multistep.ActionContinue,
			[]string{
				toSha1(abs(t, "./test-fixtures/root//another.txt")) + ".lock",
			},
		},
		{"none checksum works",
			fields{Url: []string{abs(t, "./test-fixtures/root/another.txt")}, Checksum: "none"},
			multistep.ActionContinue,
			[]string{
				toSha1(abs(t, "./test-fixtures/root/another.txt")) + ".lock",
			},
		},
		{"bad checksum removes file - checksum from string - no Checksum Type",
			fields{Extension: "txt", Url: []string{abs(t, "./test-fixtures/root/another.txt")}, Checksum: cs["/root/basic.txt"]},
			multistep.ActionHalt,
			[]string{
				toSha1(cs["/root/basic.txt"]) + ".txt.lock", // a lock file is created & deleted on mac for each download
			},
		},
		{"bad checksum removes file - checksum from string - Checksum Type",
			fields{Extension: "txt", Url: []string{abs(t, "./test-fixtures/root/another.txt")}, Checksum: "sha1:" + cs["/root/basic.txt"]},
			multistep.ActionHalt,
			[]string{
				toSha1("sha1:"+cs["/root/basic.txt"]) + ".txt.lock",
			},
		},
		{"bad checksum removes file - checksum from url - Checksum Type",
			fields{Extension: "txt", Url: []string{abs(t, "./test-fixtures/root/basic.txt")}, Checksum: "file:" + srvr.URL + "/root/another.txt.sha1sum"},
			multistep.ActionHalt,
			[]string{
				toSha1("file:"+srvr.URL+"/root/another.txt.sha1sum") + ".txt.lock",
			},
		},
		{"successfull http dl - checksum from http file - parameter",
			fields{Extension: "txt", Url: []string{srvr.URL + "/root/another.txt"}, Checksum: "file:" + srvr.URL + "/root/another.txt.sha1sum"},
			multistep.ActionContinue,
			[]string{
				toSha1("file:"+srvr.URL+"/root/another.txt.sha1sum") + ".txt",
				toSha1("file:"+srvr.URL+"/root/another.txt.sha1sum") + ".txt.lock",
			},
		},
		{"successfull http dl - checksum from http file - url",
			fields{Extension: "txt", Url: []string{srvr.URL + "/root/another.txt?checksum=file:" + srvr.URL + "/root/another.txt.sha1sum"}},
			multistep.ActionContinue,
			[]string{
				toSha1("file:"+srvr.URL+"/root/another.txt.sha1sum") + ".txt",
				toSha1("file:"+srvr.URL+"/root/another.txt.sha1sum") + ".txt.lock",
			},
		},
		{"successfull http dl - checksum from url",
			fields{Extension: "txt", Url: []string{srvr.URL + "/root/another.txt?checksum=" + cs["/root/another.txt"]}},
			multistep.ActionContinue,
			[]string{
				toSha1(cs["/root/another.txt"]) + ".txt",
				toSha1(cs["/root/another.txt"]) + ".txt.lock",
			},
		},
		{"successfull http dl - checksum from parameter - no checksum type",
			fields{Extension: "txt", Url: []string{srvr.URL + "/root/another.txt?"}, Checksum: cs["/root/another.txt"]},
			multistep.ActionContinue,
			[]string{
				toSha1(cs["/root/another.txt"]) + ".txt",
				toSha1(cs["/root/another.txt"]) + ".txt.lock",
			},
		},
		{"successfull http dl - checksum from parameter - checksum type",
			fields{Extension: "txt", Url: []string{srvr.URL + "/root/another.txt?"}, Checksum: "sha1:" + cs["/root/another.txt"]},
			multistep.ActionContinue,
			[]string{
				toSha1("sha1:"+cs["/root/another.txt"]) + ".txt",
				toSha1("sha1:"+cs["/root/another.txt"]) + ".txt.lock",
			},
		},
		{"successfull relative symlink - checksum from url",
			fields{Extension: "txt", Url: []string{"./test-fixtures/root/another.txt?checksum=" + cs["/root/another.txt"]}},
			multistep.ActionContinue,
			[]string{
				toSha1(cs["/root/another.txt"]) + ".txt.lock",
			},
		},
		{"successfull relative symlink - checksum from parameter - no checksum type",
			fields{Extension: "txt", Url: []string{"./test-fixtures/root/another.txt?"}, Checksum: cs["/root/another.txt"]},
			multistep.ActionContinue,
			[]string{
				toSha1(cs["/root/another.txt"]) + ".txt.lock",
			},
		},
		{"successfull relative symlink - checksum from parameter -  checksum type",
			fields{Extension: "txt", Url: []string{"./test-fixtures/root/another.txt?"}, Checksum: "sha1:" + cs["/root/another.txt"]},
			multistep.ActionContinue,
			[]string{
				toSha1("sha1:"+cs["/root/another.txt"]) + ".txt.lock",
			},
		},
		{"successfull absolute symlink - checksum from url",
			fields{Extension: "txt", Url: []string{abs(t, "./test-fixtures/root/another.txt") + "?checksum=" + cs["/root/another.txt"]}},
			multistep.ActionContinue,
			[]string{
				toSha1(cs["/root/another.txt"]) + ".txt.lock",
			},
		},
		{"successfull absolute symlink - checksum from parameter - no checksum type",
			fields{Extension: "txt", Url: []string{abs(t, "./test-fixtures/root/another.txt") + "?"}, Checksum: cs["/root/another.txt"]},
			multistep.ActionContinue,
			[]string{
				toSha1(cs["/root/another.txt"]) + ".txt.lock",
			},
		},
		{"successfull absolute symlink - checksum from parameter - checksum type",
			fields{Extension: "txt", Url: []string{abs(t, "./test-fixtures/root/another.txt") + "?"}, Checksum: "sha1:" + cs["/root/another.txt"]},
			multistep.ActionContinue,
			[]string{
				toSha1("sha1:"+cs["/root/another.txt"]) + ".txt.lock",
			},
		},
		{"wrong first 2 urls - absolute urls - checksum from parameter - no checksum type",
			fields{
				Url: []string{
					abs(t, "./test-fixtures/root/another.txt"),
					abs(t, "./test-fixtures/root/not_found"),
					abs(t, "./test-fixtures/root/basic.txt"),
				},
				Checksum: cs["/root/basic.txt"],
			},
			multistep.ActionContinue,
			[]string{
				toSha1(cs["/root/basic.txt"]) + ".lock",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := createTempDir(t)
			defer os.RemoveAll(dir)
			s := &StepDownload{
				TargetPath:  tt.fields.TargetPath,
				Checksum:    tt.fields.Checksum,
				ResultKey:   tt.fields.ResultKey,
				Url:         tt.fields.Url,
				Extension:   tt.fields.Extension,
				Description: tt.name,
			}
			defer os.Setenv("PACKER_CACHE_DIR", os.Getenv("PACKER_CACHE_DIR"))
			os.Setenv("PACKER_CACHE_DIR", dir)

			if got := s.Run(context.Background(), testState(t)); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("StepDownload.Run() = %v, want %v", got, tt.want)
			}
			files := listFiles(t, dir)
			if diff := cmp.Diff(tt.wantFiles, files); diff != "" {
				t.Fatalf("file list differs in %s: %s", dir, diff)
			}
		})
	}
}

func TestStepDownload_download(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Log("skipping download test on windows right now.")
		return
	}
	step := &StepDownload{
		Checksum:    "sha1:f572d396fae9206628714fb2ce00f72e94f2258f",
		Description: "ISO",
		ResultKey:   "iso_path",
		Url:         nil,
	}
	ui := &packersdk.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
		PB:     &packersdk.NoopProgressTracker{},
	}

	dir := createTempDir(t)
	defer os.RemoveAll(dir)

	defer os.Setenv("PACKER_CACHE_DIR", os.Getenv("PACKER_CACHE_DIR"))
	os.Setenv("PACKER_CACHE_DIR", dir)

	// Abs path with extension provided
	step.TargetPath = "./packer"
	step.Extension = "ova"
	_, err := step.download(context.TODO(), ui, "./test-fixtures/root/basic.txt")
	if err != nil {
		t.Fatalf("Bad: non expected error %s", err.Error())
	}
	os.RemoveAll(step.TargetPath)

	// Abs path with no extension provided
	step.TargetPath = "./packer"
	step.Extension = ""
	_, err = step.download(context.TODO(), ui, "./test-fixtures/root/basic.txt")
	if err != nil {
		t.Fatalf("Bad: non expected error %s", err.Error())
	}
	os.RemoveAll(step.TargetPath)

	// Path with file
	step.TargetPath = "./packer/file.iso"
	_, err = step.download(context.TODO(), ui, "./test-fixtures/root/basic.txt")
	if err != nil {
		t.Fatalf("Bad: non expected error %s", err.Error())
	}
	os.RemoveAll(step.TargetPath)
}

func TestStepDownload_short_timeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Log("skipping download test on windows right now.")
		return
	}

	os.Setenv("PACKER_GETTER_READ_TIMEOUT", "1ns")
	prepareGetterClient()

	// we need to re prepare the getterClient with default behavior once the test
	// is completed
	defer prepareGetterClient()
	defer os.Unsetenv("PACKER_GETTER_READ_TIMEOUT")

	srvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-time.After(time.Microsecond * 100)
		w.Write([]byte("should not receive this"))
	}))

	defer srvr.Close()
	step := &StepDownload{
		Checksum:    "sha1:f572d396fae9206628714fb2ce00f72e94f2258f",
		Description: "ISO",
		ResultKey:   "iso_path",
		Url:         nil,
	}
	ui := &packersdk.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
		PB:     &packersdk.NoopProgressTracker{},
	}

	_, err := step.download(context.TODO(), ui, srvr.URL+"/root/another.txt?")

	contextDeadlineErrorMsg := "context deadline exceeded"
	if err == nil {
		t.Fatalf("Bad: expected a '%s' error", contextDeadlineErrorMsg)
	}
	if !strings.Contains(err.Error(), contextDeadlineErrorMsg) {
		t.Fatalf("Bad: expected a '%s' error, got %s", contextDeadlineErrorMsg, err.Error())
	}

}

func TestStepDownload_WindowsParseSourceURL(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skip windows specific tests")
	}

	source := "\\\\hostname\\dir\\filename.txt"
	url, err := parseSourceURL(source)
	if err != nil {
		t.Fatalf("bad: parsing source url failed: %s", err.Error())
	}
	if url.Scheme != "smb" {
		t.Fatalf("bad: url should contain smb scheme but contains: %s", url.Scheme)
	}
}

func TestStepDownload_ParseSourceSmbURL(t *testing.T) {
	source := "smb://hostname/dir/filename.txt"
	url, err := parseSourceURL(source)
	if err != nil {
		t.Fatalf("bad: parsing source url failed: %s", err.Error())
	}
	if url.Scheme != "smb" {
		t.Fatalf("bad: url should contain smb scheme but contains: %s", url.Scheme)
	}
	if url.String() != "smb://hostname/dir/filename.txt" {
		t.Fatalf("bad: url should contain smb scheme but contains: %s", url.String())
	}
}

func createTempDir(t *testing.T) string {
	dir, err := tmp.Dir("pkr")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return dir
}

func listFiles(t *testing.T, dir string) []string {
	fs, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, file := range fs {
		if file.Name() == "." {
			continue
		}
		files = append(files, file.Name())
	}

	return files
}
