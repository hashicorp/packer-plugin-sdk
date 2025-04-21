// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rpc

import (
	"fmt"
	"log"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// An implementation of packersdk.Ui where the Ui is actually executed
// over an RPC connection.
type Ui struct {
	commonClient
	endpoint string
}

var _ packersdk.Ui = new(Ui)

// UiServer wraps a packersdk.Ui implementation and makes it exportable
// as part of a Golang RPC server.
type UiServer struct {
	ui       packersdk.Ui
	register func(name string, rcvr interface{}) error
}

// The arguments sent to Ui.Machine
type UiMachineArgs struct {
	Category string
	Args     []string
}

func (u *Ui) Askf(query string, args ...any) (string, error) {
	return u.Ask(fmt.Sprintf(query, args...))
}
func (u *Ui) Ask(query string) (result string, err error) {
	err = u.client.Call("Ui.Ask", query, &result)
	return
}

func (u *Ui) Errorf(message string, args ...any) {
	u.Error(fmt.Sprintf(message, args...))
}
func (u *Ui) Error(message string) {
	if err := u.client.Call("Ui.Error", message, new(interface{})); err != nil {
		log.Printf("Error in Ui.Error RPC call: %s", err)
	}
}

func (u *Ui) Machine(t string, args ...string) {
	rpcArgs := &UiMachineArgs{
		Category: t,
		Args:     args,
	}

	if err := u.client.Call("Ui.Machine", rpcArgs, new(interface{})); err != nil {
		log.Printf("Error in Ui.Machine RPC call: %s", err)
	}
}

func (u *Ui) Sayf(message string, args ...any) {
	u.Say(fmt.Sprintf(message, args...))
}
func (u *Ui) Say(message string) {
	if err := u.client.Call("Ui.Say", message, new(interface{})); err != nil {
		log.Printf("Error in Ui.Say RPC call: %s", err)
	}
}

func (u *UiServer) Ask(query string, reply *string) (err error) {
	*reply, err = u.ui.Ask(query)
	return
}

func (u *UiServer) Error(message *string, reply *interface{}) error {
	u.ui.Error(*message)

	*reply = nil
	return nil
}

func (u *UiServer) Machine(args *UiMachineArgs, reply *interface{}) error {
	u.ui.Machine(args.Category, args.Args...)

	*reply = nil
	return nil
}

func (u *UiServer) Say(message *string, reply *interface{}) error {
	u.ui.Say(*message)

	*reply = nil
	return nil
}
