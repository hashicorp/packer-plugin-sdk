// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/ssh"
)

func TestGeneratePair_parseable(t *testing.T) {
	tests := []struct {
		t Algorithm
	}{
		{RSA},
		{ECDSA},
		{ED25519},
	}
	for _, tt := range tests {
		t.Run(tt.t.String(), func(t *testing.T) {
			got, err := GeneratePair(tt.t, nil, 0)
			if err != nil {
				t.Errorf("GeneratePair() error = %v", err)
				return
			}

			privateKey, err := ssh.ParsePrivateKey(got.Private)
			if err != nil {
				t.Fatal(err)
			}
			publicKey, _, _, _, err := ssh.ParseAuthorizedKey(got.Public)
			if err != nil {
				t.Fatalf("%v: %s", err, got.Public)
			}
			if diff := cmp.Diff(privateKey.PublicKey().Marshal(), publicKey.Marshal()); diff != "" {
				t.Fatalf("wrong public key: %s", diff)
			}
		})
	}
}

func TestGeneratePair_dsaNotSupported(t *testing.T) {
	for _, bits := range []int{0, 1024} {
		pair, err := GeneratePair(DSA, nil, bits)
		if !errors.Is(err, ErrDSANotSupported) {
			t.Fatalf("GeneratePair(DSA, %d) error = %v, want %v", bits, err, ErrDSANotSupported)
		}
		if pair != nil {
			t.Fatalf("GeneratePair(DSA, %d) returned a pair, want nil", bits)
		}
	}
}
