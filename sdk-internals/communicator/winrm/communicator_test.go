// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package winrm

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/dylanmei/winrmtest"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

const PAYLOAD = "stuff"
const BASE64_ENCODED_PAYLOAD = "c3R1ZmY="

func newMockWinRMServer(t *testing.T) *winrmtest.Remote {
	wrm := winrmtest.NewRemote()

	wrm.CommandFunc(
		winrmtest.MatchText("echo foo"),
		func(out, err io.Writer) int {
			out.Write([]byte("foo"))
			return 0
		})

	wrm.CommandFunc(
		winrmtest.MatchPattern(`^echo c29tZXRoaW5n >> ".*"$`),
		func(out, err io.Writer) int {
			return 0
		})

	wrm.CommandFunc(
		winrmtest.MatchPattern(`^echo `+BASE64_ENCODED_PAYLOAD+` >> ".*"$`),
		func(out, err io.Writer) int {
			return 0
		})

	wrm.CommandFunc(
		winrmtest.MatchPattern(`^powershell.exe -EncodedCommand .*$`),
		func(out, err io.Writer) int {
			out.Write([]byte(BASE64_ENCODED_PAYLOAD))
			return 0
		})

	wrm.CommandFunc(
		winrmtest.MatchText("powershell"),
		func(out, err io.Writer) int {
			return 0
		})
	wrm.CommandFunc(
		winrmtest.MatchText(`powershell -Command "(Get-Item C:/Temp/packer.cmd) -is [System.IO.DirectoryInfo]"`),
		func(out, err io.Writer) int {
			out.Write([]byte("False"))
			return 0
		})

	return wrm
}

func TestStart(t *testing.T) {
	wrm := newMockWinRMServer(t)
	defer wrm.Close()

	c, err := New(&Config{
		Host:     wrm.Host,
		Port:     wrm.Port,
		Username: "user",
		Password: "pass",
		Timeout:  30 * time.Second,
	})
	if err != nil {
		t.Fatalf("error creating communicator: %s", err)
	}

	var cmd packersdk.RemoteCmd
	stdout := new(bytes.Buffer)
	cmd.Command = "echo foo"
	cmd.Stdout = stdout
	ctx := context.Background()
	err = c.Start(ctx, &cmd)
	if err != nil {
		t.Fatalf("error executing remote command: %s", err)
	}
	cmd.Wait()

	if stdout.String() != "foo" {
		t.Fatalf("bad command response: expected %q, got %q", "foo", stdout.String())
	}
}

func TestUpload(t *testing.T) {
	wrm := newMockWinRMServer(t)
	defer wrm.Close()

	c, err := New(&Config{
		Host:     wrm.Host,
		Port:     wrm.Port,
		Username: "user",
		Password: "pass",
		Timeout:  30 * time.Second,
	})
	if err != nil {
		t.Fatalf("error creating communicator: %s", err)
	}
	file := "C:/Temp/packer.cmd"
	err = c.Upload(file, strings.NewReader(PAYLOAD), nil)
	if err != nil {
		t.Fatalf("error uploading file: %s", err)
	}

	dest := new(bytes.Buffer)
	err = c.Download(file, dest)
	if err != nil {
		t.Fatalf("error downloading file: %s", err)
	}
	downloadedPayload := dest.String()

	if downloadedPayload != PAYLOAD {
		t.Fatalf("files are not equal: expected [%s] length: %v, got [%s] length %v", PAYLOAD, len(PAYLOAD), downloadedPayload, len(downloadedPayload))
	}
}

func TestUpload_nilFileInfo(t *testing.T) {
	wrm := newMockWinRMServer(t)
	defer wrm.Close()

	c, err := New(&Config{
		Host:     wrm.Host,
		Port:     wrm.Port,
		Username: "user",
		Password: "pass",
		Timeout:  30 * time.Second,
	})
	if err != nil {
		t.Fatalf("error creating communicator: %s", err)
	}
	file := "C:\\Temp\\"
	err = c.Upload(file, strings.NewReader(PAYLOAD), nil)
	if err == nil {
		t.Fatalf("Should have errored because of nil fileinfo")
	}
}
