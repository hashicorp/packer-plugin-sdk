// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shell_local

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type Communicator struct {
	ExecuteCommand []string
}

func (c *Communicator) Start(ctx context.Context, cmd *packersdk.RemoteCmd) error {
	if len(c.ExecuteCommand) == 0 {
		return fmt.Errorf("Error launching command via shell-local communicator: No ExecuteCommand provided")
	}

	// Build the local command to execute
	log.Printf("[INFO] (shell-local communicator): Executing local shell command %s", c.ExecuteCommand)
	localCmd := exec.CommandContext(ctx, c.ExecuteCommand[0], c.ExecuteCommand[1:]...)
	localCmd.Stdin = cmd.Stdin
	localCmd.Stdout = cmd.Stdout
	localCmd.Stderr = cmd.Stderr

	// Start it. If it doesn't work, then error right away.
	if err := localCmd.Start(); err != nil {
		return err
	}

	// We've started successfully. Start a goroutine to wait for
	// it to complete and track exit status.
	go func() {
		var exitStatus int
		err := localCmd.Wait()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitStatus = 1

				// There is no process-independent way to get the REAL
				// exit status so we just try to go deeper.
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					exitStatus = status.ExitStatus()
				}
			}
		}

		cmd.SetExited(exitStatus)
	}()

	return nil
}

func (c *Communicator) Upload(string, io.Reader, *os.FileInfo) error {
	return fmt.Errorf("upload not supported")
}

func (c *Communicator) UploadDir(string, string, []string) error {
	return fmt.Errorf("uploadDir not supported")
}

func (c *Communicator) Download(string, io.Writer) error {
	return fmt.Errorf("download not supported")
}

func (c *Communicator) DownloadDir(string, string, []string) error {
	return fmt.Errorf("downloadDir not supported")
}
