// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package packer

import (
	"bytes"
	"io"
	"sort"
	"strings"
	"sync"
)

type secretFilter struct {
	s map[string]struct{}
	m sync.Mutex
	w io.Writer
}

func (l *secretFilter) Set(secrets ...string) {
	l.m.Lock()
	defer l.m.Unlock()
	for _, s := range secrets {
		l.s[s] = struct{}{}
	}
}

func (l *secretFilter) SetOutput(output io.Writer) {
	l.m.Lock()
	defer l.m.Unlock()
	l.w = output
}

func (l *secretFilter) Write(p []byte) (n int, err error) {
	for _, s := range l.secrets() {
		p = bytes.ReplaceAll(p, []byte(s), []byte("<sensitive>"))
	}
	return l.w.Write(p)
}

// FilterString will overwrite any senstitive variables in a string, returning
// the filtered string.
func (l *secretFilter) FilterString(message string) string {
	for _, s := range l.secrets() {
		message = strings.ReplaceAll(message, s, "<sensitive>")
	}
	return message
}

// secrets returns the non-empty registered secrets ordered longest first.
// Ranging over the map directly is randomly ordered, so when one secret is a
// substring of another (e.g. "ubuntu" and "ubuntu-22.04") the shorter one could
// be replaced first, leaving the remainder of the longer secret in the output.
// Redacting the longest match first keeps that from happening and makes the
// result deterministic.
func (l *secretFilter) secrets() []string {
	secrets := make([]string, 0, len(l.s))
	for s := range l.s {
		if s != "" {
			secrets = append(secrets, s)
		}
	}
	sort.Slice(secrets, func(i, j int) bool {
		if len(secrets[i]) != len(secrets[j]) {
			return len(secrets[i]) > len(secrets[j])
		}
		return secrets[i] < secrets[j]
	})
	return secrets
}

var LogSecretFilter secretFilter

func init() {
	LogSecretFilter.s = make(map[string]struct{})
}
