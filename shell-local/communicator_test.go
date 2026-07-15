// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package shell_local

import (
	"bytes"
	"context"
	"os"
	"runtime"
	"strings"
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestCommunicator_impl(t *testing.T) {
	var _ packersdk.Communicator = new(Communicator)
}

func TestCommunicator(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows not supported for this test")
		return
	}

	c := &Communicator{
		ExecuteCommand: []string{"/bin/sh", "-c", "echo foo"},
	}

	var buf bytes.Buffer
	cmd := &packersdk.RemoteCmd{
		Stdout: &buf,
	}

	ctx := context.Background()
	if err := c.Start(ctx, cmd); err != nil {
		t.Fatalf("err: %s", err)
	}

	cmd.Wait()

	if cmd.ExitStatus() != 0 {
		t.Fatalf("err bad exit status: %d", cmd.ExitStatus())
	}

	if strings.TrimSpace(buf.String()) != "foo" {
		t.Fatalf("bad: %s", buf.String())
	}
}

func TestDefaultInlineExecuteCommandRunsScriptOpenForWriting(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows does not reject executing a file that is open for writing")
	}

	script, err := os.CreateTemp(t.TempDir(), "packer-shell")
	if err != nil {
		t.Fatalf("creating script: %s", err)
	}
	t.Cleanup(func() { script.Close() })

	if _, err := script.WriteString("echo foo\n"); err != nil {
		t.Fatalf("writing script: %s", err)
	}
	if err := script.Chmod(0700); err != nil {
		t.Fatalf("making script executable: %s", err)
	}

	config := &Config{}
	config.Inline = []string{"echo foo"}
	if err := Validate(config); err != nil {
		t.Fatalf("validating config: %s", err)
	}
	config.generatedData = make(map[string]interface{})

	executeCommand, err := createInterpolatedCommands(config, script.Name(), "")
	if err != nil {
		t.Fatalf("interpolating command: %s", err)
	}

	communicator := &Communicator{ExecuteCommand: executeCommand}
	var stdout bytes.Buffer
	cmd := &packersdk.RemoteCmd{Stdout: &stdout}
	if err := communicator.Start(context.Background(), cmd); err != nil {
		t.Fatalf("starting command: %s", err)
	}
	cmd.Wait()

	if cmd.ExitStatus() != 0 {
		t.Fatalf("bad exit status: %d", cmd.ExitStatus())
	}
	if strings.TrimSpace(stdout.String()) != "foo" {
		t.Fatalf("bad output: %s", stdout.String())
	}
}
