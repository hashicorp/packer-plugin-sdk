// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package ssh implements the SSH communicator. Plugin maintainers should not
// import this package directly, instead using the tooling in the
// "packer-plugin-sdk/communicator" module.

package ssh

import (
	"golang.org/x/crypto/ssh/agent"
	"log"
)

func getSSHAgent() (agent.ExtendedAgent, error) {
	log.Printf("[INFO] No Agent creation ß")
	return nil, nil
}
