// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package ssh

import "testing"

func TestGetSSHAgentConnection(t *testing.T) {
	conn, err := GetSSHAgentConnection()
	if err != nil {
		t.Fatal("Expected nil error when named pipe exist, got non-nil")
	}
	if conn == nil {
		t.Fatal("Expected connection, got nil")
	}
}
