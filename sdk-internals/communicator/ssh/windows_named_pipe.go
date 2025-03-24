// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package ssh implements the SSH communicator. Plugin maintainers should not
// import this package directly, instead using the tooling in the
// "packer-plugin-sdk/communicator" module.

//go:build windows

package ssh

import (
	"fmt"
	"github.com/Microsoft/go-winio"
	"net"
)

func getSSHAgentConnection() (net.Conn, error) {
	pipePath := "\\\\.\\pipe\\openssh-ssh-agent"

	sshAgent, err := winio.DialPipe(pipePath, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to SSH Agent named pipe %q: %s", pipePath, err)
	}

	return sshAgent, nil
}
