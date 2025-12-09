// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !windows

package ssh

import (
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
)

func TestGetSSHAgentConnection_NoEnvVar(t *testing.T) {
	// Unset SSH_AUTH_SOCK
	os.Unsetenv("SSH_AUTH_SOCK")

	conn, err := GetSSHAgentConnection()
	if err == nil {
		t.Fatal("Expected error when SSH_AUTH_SOCK is not set, got nil")
	}

	if conn != nil {
		t.Fatal("Expected nil connection, got non-nil")
	}

	assert.Equal(t, "SSH_AUTH_SOCK is not set", err.Error())
}

func TestGetSSHAgentConnection_InvalidSocket(t *testing.T) {
	os.Setenv("SSH_AUTH_SOCK", "/invalid/path/to/socket")

	conn, err := GetSSHAgentConnection()
	if err == nil {
		t.Fatal("Expected error when SSH_AUTH_SOCK points to an invalid path, got nil")
	}

	if conn != nil {
		t.Fatal("Expected nil connection, got non-nil")
	}

	assert.Equal(t, "Cannot connect to SSH Agent socket \"/invalid/path/to/socket\": dial unix /invalid/path/to/socket: connect: no such file or directory", err.Error())
}

func TestGetSSHAgentConnection_ValidSocket(t *testing.T) {
	// Create a temporary Unix socket for testing
	socketPath := "/tmp/test-ssh-agent.sock"
	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to create mock SSH agent socket: %v", err)
	}
	defer ln.Close()
	defer os.Remove(socketPath) // Cleanup after test

	// Set the environment variable to use the mock socket
	os.Setenv("SSH_AUTH_SOCK", socketPath)

	conn, err := GetSSHAgentConnection()
	if err != nil {
		t.Fatalf("Expected successful connection, got error: %v", err)
	}
	if conn == nil {
		t.Fatal("Expected non-nil connection, got nil")
	}
	conn.Close() // Close the connection after testing

	assert.Nil(t, err)
}
