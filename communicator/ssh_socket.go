// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package communicator

import (
	"fmt"
	"net"
	"os"
)

func getSSHAgentConnection() (net.Conn, error) {
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
