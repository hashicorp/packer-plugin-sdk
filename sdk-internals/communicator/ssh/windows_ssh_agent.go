// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package ssh implements the SSH communicator. Plugin maintainers should not
// import this package directly, instead using the tooling in the
// "packer-plugin-sdk/communicator" module.

//go:build windows

package ssh

import (
	"errors"
	"github.com/Microsoft/go-winio"
	"golang.org/x/crypto/ssh/agent"
	"log"
)

func getSSHAgent() (agent.ExtendedAgent, error) {

	pipePath := "\\\\.\\pipe\\openssh-ssh-agent"
	agentConn, err := winio.DialPipe(pipePath, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to SSH agent pipe: %v", err)
	}

	forwardingAgent := agent.NewClient(agentConn)
	if forwardingAgent == nil {
		log.Printf("[ERROR] Could not create agent client")
		agentConn.Close()
		return nil, errors.New("[ERROR] Could not create agent client")
	}

	return forwardingAgent, nil
}
