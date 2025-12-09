// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

// Package ssh implements the SSH communicator. Plugin maintainers should not
// import this package directly, instead using the tooling in the
// "packer-plugin-sdk/communicator" module.

//go:build !windows

package ssh

import (
	"fmt"
	"net"
	"os"
)

func GetSSHAgentConnection() (net.Conn, error) {
	authSock := os.Getenv("SSH_AUTH_SOCK")
	if authSock == "" {
		return nil, fmt.Errorf("SSH_AUTH_SOCK is not set")
	}

	sshAgent, err := net.Dial("unix", authSock)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to SSH Agent socket %q: %s", authSock, err)
	}

	return sshAgent, nil
}
