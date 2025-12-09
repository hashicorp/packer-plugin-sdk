// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
)

// TestUi creates a simple UI for use in testing.
// It's not meant for "real" use.
func TestUi(t *testing.T) Ui {
	var buf bytes.Buffer
	return &BasicUi{
		Reader:      &buf,
		Writer:      io.Discard,
		ErrorWriter: io.Discard,
		PB:          &NoopProgressTracker{},
	}
}

type SayMessage struct {
	Message string
	SayTime time.Time
}

type MockUi struct {
	AskCalled      bool
	AskQuery       string
	ErrorCalled    bool
	ErrorMessage   string
	MachineCalled  bool
	MachineType    string
	MachineArgs    []string
	MessageCalled  bool
	MessageMessage string
	SayCalled      bool
	SayMessages    []SayMessage

	TrackProgressCalled    bool
	ProgressBarAddCalled   bool
	ProgressBarCloseCalled bool
}

func (u *MockUi) Askf(query string, args ...any) (string, error) {
	return u.Ask(fmt.Sprintf(query, args...))
}
func (u *MockUi) Ask(query string) (string, error) {
	u.AskCalled = true
	u.AskQuery = query
	return "foo", nil
}

func (u *MockUi) Errorf(message string, args ...any) {
	u.Error(fmt.Sprintf(message, args...))
}
func (u *MockUi) Error(message string) {
	u.ErrorCalled = true
	u.ErrorMessage = message
}

func (u *MockUi) Machine(t string, args ...string) {
	u.MachineCalled = true
	u.MachineType = t
	u.MachineArgs = args
}

// Deprecated: Use `Say` instead.
func (u *MockUi) Message(message string) {
	u.Say(message)
}

func (u *MockUi) Sayf(message string, args ...any) {
	u.Say(fmt.Sprintf(message, args...))
}
func (u *MockUi) Say(message string) {
	u.SayCalled = true
	sayMessage := SayMessage{
		Message: message,
		SayTime: time.Now(),
	}
	u.SayMessages = append(u.SayMessages, sayMessage)
}

func (u *MockUi) TrackProgress(_ string, _, _ int64, stream io.ReadCloser) (body io.ReadCloser) {
	u.TrackProgressCalled = true

	return &readCloser{
		read: func(p []byte) (int, error) {
			u.ProgressBarAddCalled = true
			return stream.Read(p)
		},
		close: func() error {
			u.ProgressBarCloseCalled = true
			return stream.Close()
		},
	}
}

type readCloser struct {
	read  func([]byte) (int, error)
	close func() error
}

func (c *readCloser) Close() error               { return c.close() }
func (c *readCloser) Read(p []byte) (int, error) { return c.read(p) }
