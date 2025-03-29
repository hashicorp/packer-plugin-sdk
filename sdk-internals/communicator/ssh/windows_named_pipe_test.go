// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package ssh

import "testing"

func TestGetSSHAgentConnection_NoPipe(t *testing.T) {
	conn, err := GetSSHAgentConnection()
	if err == nil {
		t.Fatal("Expected error when named pipe does not exist, got nil")
	}
	if conn != nil {
		t.Fatal("Expected nil connection, got non-nil")
	}
}
