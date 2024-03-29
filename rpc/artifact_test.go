// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"reflect"
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestArtifactRPC(t *testing.T) {
	// Create the interface to test
	a := new(packersdk.MockArtifact)

	// Start the server
	client, server := testClientServer(t)
	defer client.Close()
	defer server.Close()
	server.RegisterArtifact(a)

	aClient := client.Artifact()

	// Test
	if aClient.BuilderId() != "bid" {
		t.Fatalf("bad: %s", aClient.BuilderId())
	}

	if !reflect.DeepEqual(aClient.Files(), []string{"a", "b"}) {
		t.Fatalf("bad: %#v", aClient.Files())
	}

	if aClient.Id() != "id" {
		t.Fatalf("bad: %s", aClient.Id())
	}

	if aClient.String() != "string" {
		t.Fatalf("bad: %s", aClient.String())
	}
}

func TestArtifact_Implements(t *testing.T) {
	var _ packersdk.Artifact = new(artifact)
}
