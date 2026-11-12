// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"testing"
	"time"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

func TestPluginServerCloseBeforeAccept(t *testing.T) {
	client, server := testClientServer(t)
	defer client.Close()
	defer server.Close()
	if err := server.RegisterArtifact(&packersdk.MockArtifact{IdValue: "parent"}); err != nil {
		t.Fatal(err)
	}

	sub := newServerWithMux(server.mux, server.mux.NextId())
	done := make(chan struct{})
	go func() {
		defer close(done)
		sub.Serve()
	}()

	if err := sub.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Serve did not return after Close")
	}
	if got := client.Artifact().Id(); got != "parent" {
		t.Fatalf("shared mux closed: id %q", got)
	}
}

func TestPluginServerCloseStopsOneStream(t *testing.T) {
	client, server := testClientServer(t)
	defer client.Close()
	defer server.Close()
	if err := server.RegisterArtifact(&packersdk.MockArtifact{IdValue: "parent"}); err != nil {
		t.Fatal(err)
	}

	id := server.mux.NextId()
	sub := newServerWithMux(server.mux, id)
	if err := sub.RegisterArtifact(&packersdk.MockArtifact{IdValue: "child"}); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		sub.Serve()
	}()

	child, err := newClientWithMux(client.mux, id)
	if err != nil {
		t.Fatal(err)
	}
	defer child.Close()
	if got := child.Artifact().Id(); got != "child" {
		t.Fatalf("child id %q", got)
	}

	if err := sub.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Serve did not return after Close")
	}
	if got := client.Artifact().Id(); got != "parent" {
		t.Fatalf("shared mux closed: id %q", got)
	}
}

func TestArtifactDestroyClosesItsStream(t *testing.T) {
	client, server := testClientServer(t)
	defer client.Close()
	defer server.Close()
	if err := server.RegisterArtifact(&packersdk.MockArtifact{IdValue: "parent"}); err != nil {
		t.Fatal(err)
	}

	local := &packersdk.MockArtifact{IdValue: "child"}
	id := server.mux.NextId()
	sub := newServerWithMux(server.mux, id)
	if err := sub.RegisterArtifact(local); err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		sub.Serve()
	}()

	remoteClient, err := newClientWithMux(client.mux, id)
	if err != nil {
		t.Fatal(err)
	}
	remote := remoteClient.Artifact()
	if got := remote.Id(); got != "child" {
		t.Fatalf("child id %q", got)
	}
	if err := remote.Destroy(); err != nil {
		t.Fatal(err)
	}
	if !local.DestroyCalled {
		t.Fatal("artifact Destroy was not called")
	}
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("artifact server still running after Destroy")
	}
	if got := client.Artifact().Id(); got != "parent" {
		t.Fatalf("shared mux closed: id %q", got)
	}
}
