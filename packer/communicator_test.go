// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/iochan"
	"golang.org/x/sync/errgroup"
)

func TestRemoteCmd_StartWithUi(t *testing.T) {
	data := []string{
		"hello",
		"world",
		"foo",
		"there",
	}

	originalOutputReader, originalOutputWriter := io.Pipe()
	uilOutputReader, uilOutputWriter := io.Pipe()

	testComm := new(MockCommunicator)
	testComm.StartStdout = strings.Join(data, "\n") + "\n"
	testUi := &BasicUi{
		Reader: new(bytes.Buffer),
		Writer: uilOutputWriter,
	}

	rc := &RemoteCmd{
		Command: "test",
		Stdout:  originalOutputWriter,
	}
	ctx := context.TODO()

	wg := errgroup.Group{}

	testPrintFn := func(in io.Reader, expected []string) error {
		i := 0
		got := []string{}
		for output := range iochan.DelimReader(in, '\n') {
			got = append(got, strings.TrimSpace(output))
			i++
			if i == len(expected) {
				// here ideally the LineReader chan should be closed, but since
				// the stream virtually has no ending we need to leave early.
				break
			}
		}
		if diff := cmp.Diff(got, expected); diff != "" {
			t.Fatalf("bad output: %s", diff)
		}
		return nil
	}

	wg.Go(func() error { return testPrintFn(uilOutputReader, data) })
	wg.Go(func() error { return testPrintFn(originalOutputReader, data) })

	err := rc.RunWithUi(ctx, testComm, testUi)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	wg.Wait()
}

func TestRemoteCmd_Wait(t *testing.T) {
	var cmd RemoteCmd

	result := make(chan bool)
	go func() {
		cmd.Wait()
		result <- true
	}()

	cmd.SetExited(42)

	select {
	case <-result:
		// Success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("never got exit notification")
	}
}

func TestRemoteCmd_RunWithUi_StartErrorDoesNotLeakGoroutines(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("boom")
	rc := &RemoteCmd{Command: "test"}
	ui := &BasicUi{
		Reader: new(bytes.Buffer),
		Writer: io.Discard,
	}

	err := rc.RunWithUi(context.Background(), errorCommunicator{err: errExpected}, ui)
	if !errors.Is(err, errExpected) {
		t.Fatalf("expected %v, got %v", errExpected, err)
	}

	assertNoGoroutineStack(t,
		"github.com/hashicorp/packer-plugin-sdk/packer.(*RemoteCmd).RunWithUi.func",
		"github.com/hashicorp/packer-plugin-sdk/packer.(*RemoteCmd).Wait",
	)
}

func TestRemoteCmd_RunWithUi_ContextCancelDoesNotLeakGoroutines(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rc := &RemoteCmd{Command: "test"}
	ui := &BasicUi{
		Reader: new(bytes.Buffer),
		Writer: io.Discard,
	}
	comm := startedCommunicator{started: make(chan struct{})}
	errCh := make(chan error, 1)

	go func() {
		errCh <- rc.RunWithUi(ctx, comm, ui)
	}()

	select {
	case <-comm.started:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Start was never called")
	}

	cancel()

	select {
	case err := <-errCh:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected %v, got %v", context.Canceled, err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("RunWithUi did not return promptly after cancellation")
	}

	assertNoGoroutineStack(t,
		"github.com/hashicorp/packer-plugin-sdk/packer.(*RemoteCmd).RunWithUi.func",
		"github.com/hashicorp/packer-plugin-sdk/packer.(*RemoteCmd).Wait",
	)
}

type errorCommunicator struct {
	err error
}

func (c errorCommunicator) Start(context.Context, *RemoteCmd) error {
	return c.err
}

func (errorCommunicator) Upload(string, io.Reader, *os.FileInfo) error {
	panic("unexpected Upload call")
}

func (errorCommunicator) UploadDir(string, string, []string) error {
	panic("unexpected UploadDir call")
}

func (errorCommunicator) Download(string, io.Writer) error {
	panic("unexpected Download call")
}

func (errorCommunicator) DownloadDir(string, string, []string) error {
	panic("unexpected DownloadDir call")
}

type startedCommunicator struct {
	started chan struct{}
}

func (c startedCommunicator) Start(context.Context, *RemoteCmd) error {
	close(c.started)
	return nil
}

func (startedCommunicator) Upload(string, io.Reader, *os.FileInfo) error {
	panic("unexpected Upload call")
}

func (startedCommunicator) UploadDir(string, string, []string) error {
	panic("unexpected UploadDir call")
}

func (startedCommunicator) Download(string, io.Writer) error {
	panic("unexpected Download call")
}

func (startedCommunicator) DownloadDir(string, string, []string) error {
	panic("unexpected DownloadDir call")
}

func assertNoGoroutineStack(t *testing.T, needles ...string) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		var stacks bytes.Buffer
		if err := pprof.Lookup("goroutine").WriteTo(&stacks, 2); err != nil {
			t.Fatalf("failed to dump goroutines: %v", err)
		}

		allGone := true
		for _, needle := range needles {
			if strings.Contains(stacks.String(), needle) {
				allGone = false
				break
			}
		}

		if allGone {
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	var stacks bytes.Buffer
	if err := pprof.Lookup("goroutine").WriteTo(&stacks, 2); err != nil {
		t.Fatalf("failed to dump goroutines: %v", err)
	}

	t.Fatalf("expected goroutines to exit, still found matching stacks:\n%s", stacks.String())
}
