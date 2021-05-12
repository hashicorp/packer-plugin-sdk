package commonsteps

import (
	"bytes"
	"context"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestStepCreateCD_Impl(t *testing.T) {
	var raw interface{}
	raw = new(StepCreateCD)
	if _, ok := raw.(multistep.Step); !ok {
		t.Fatalf("StepCreateCD should be a step")
	}
}

func testStepCreateCDState(t *testing.T) multistep.StateBag {
	state := new(multistep.BasicStateBag)
	state.Put("ui", &packersdk.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	})
	return state
}

func createFiles(t *testing.T, rootFolder string, expected map[string]string) {
	for fname, content := range expected {
		path := filepath.Join(rootFolder, fname)
		err := os.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			t.Fatalf("mkdir -p: %s", err)
		}
		err = ioutil.WriteFile(path, []byte(content), 0666)
		if err != nil {
			t.Fatalf("writing file: %s", err)
		}
	}
}

func checkFiles(t *testing.T, rootFolder string, expected map[string]string) {
	err := filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatalf("walking folder: %s", err)
		}

		if !d.IsDir() {
			name, _ := filepath.Rel(rootFolder, path)

			expectedContent, ok := expected[name]
			if !ok {
				t.Fatalf("unexpected file: %s", name)
			}

			content, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatalf("reading file: %s", err)
			}
			if string(content) != expectedContent {
				t.Fatalf("unexpected content: %s", name)
			}

			delete(expected, name)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir: %v", err)
	}
	if len(expected) != 0 {
		t.Fatalf("missing files: %v", expected)
	}
}

func TestStepCreateCD(t *testing.T) {
	if os.Getenv("PACKER_ACC") == "" {
		t.Skip("This test is only run with PACKER_ACC=1 due to the requirement of access to the disk management binaries.")
	}
	state := testStepCreateCDState(t)
	step := new(StepCreateCD)

	dir, err := ioutil.TempDir("", "packer")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	expected := map[string]string{
		"test folder/b/test1": "1",
		"test folder/b/test2": "2",
		"test folder 2/x":     "3",
		"test_cd_roms.tmp":    "4",
		"test cd files.tmp":   "5",
		"Test-Test-Test5.tmp": "6",
	}

	createFiles(t, dir, expected)
	files := []string{"test folder", "test folder 2/", "test_cd_roms.tmp", "test cd files.tmp", "Test-Test-Test5.tmp"}

	step.Files = make([]string, len(files))
	for i, fname := range files {
		step.Files[i] = filepath.Join(dir, fname)
	}
	action := step.Run(context.Background(), state)

	if err, ok := state.GetOk("error"); ok {
		t.Fatalf("state should be ok for %v: %s", step.Files, err)
	}

	if action != multistep.ActionContinue {
		t.Fatalf("bad action: %#v for %v", action, step.Files)
	}

	CD_path := state.Get("cd_path").(string)

	if _, err := os.Stat(CD_path); err != nil {
		t.Fatalf("file not found: %s for %v", CD_path, step.Files)
	}

	checkFiles(t, step.rootFolder, expected)

	step.Cleanup(state)

	if _, err := os.Stat(CD_path); err == nil {
		t.Fatalf("file found: %s for %v", CD_path, step.Files)
	}
	if _, err := os.Stat(step.rootFolder); err == nil {
		t.Fatalf("folder found: %s", step.rootFolder)
	}
}

func TestStepCreateCD_missing(t *testing.T) {
	if os.Getenv("PACKER_ACC") == "" {
		t.Skip("This test is only run with PACKER_ACC=1 due to the requirement of access to the disk management binaries.")
	}
	state := testStepCreateCDState(t)
	step := new(StepCreateCD)

	dir, err := ioutil.TempDir("", "packer")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(dir)

	step.Files = []string{"missing file.tmp"}
	if action := step.Run(context.Background(), state); action != multistep.ActionHalt {
		t.Fatalf("bad action: %#v for %v", action, step.Files)
	}

	if _, ok := state.GetOk("error"); !ok {
		t.Fatalf("state should not be ok for %v", step.Files)
	}

	CD_path := state.Get("cd_path")

	if CD_path != nil {
		t.Fatalf("CD_path is not nil for %v", step.Files)
	}

	checkFiles(t, step.rootFolder, nil)

	step.Cleanup(state)

	if _, err := os.Stat(step.rootFolder); err == nil {
		t.Fatalf("folder found: %s", step.rootFolder)
	}

	step.Cleanup(state)

	if _, err := os.Stat(step.rootFolder); err == nil {
		t.Fatalf("folder found: %s", step.rootFolder)
	}
}
