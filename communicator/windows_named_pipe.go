// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package communicator

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
