// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"crypto/dsa" //nolint:all
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

type Algorithm int

//go:generate enumer -type Algorithm -transform snake
const (
	RSA Algorithm = iota
	// Deprecated: DSA is no longer supported by golang.org/x/crypto/ssh.
	// GeneratePair returns ErrDSANotSupported for this algorithm.
	DSA
	ECDSA
	ED25519
)

var (
	ErrUnknownAlgorithm    = fmt.Errorf("sshkey: unknown private key algorithm")
	ErrInvalidRSAKeySize   = fmt.Errorf("sshkey: invalid private key rsa size: must be more than 1024")
	ErrInvalidECDSAKeySize = fmt.Errorf("sshkey: invalid private key ecdsa size, must be one of 256, 384 or 521")
	ErrDSANotSupported     = fmt.Errorf("sshkey: DSA keys are no longer supported, use ecdsa, ed25519 or rsa instead")

	// Deprecated: DSA is no longer supported and this error is no longer
	// returned. See ErrDSANotSupported.
	ErrInvalidDSAKeySize = fmt.Errorf("sshkey: invalid private key dsa size, must be one of 1024, 2048 or 3072")
)

// Pair represents an ssh key pair, as in
type Pair struct {
	Private []byte
	Public  []byte
}

func NewPair(public, private any) (*Pair, error) {
	kb, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return nil, err
	}

	privBlk := &pem.Block{
		Type:    "PRIVATE KEY",
		Headers: nil,
		Bytes:   kb,
	}

	publicKey, err := ssh.NewPublicKey(public)
	if err != nil {
		return nil, err
	}
	return &Pair{
		Private: pem.EncodeToMemory(privBlk),
		Public:  ssh.MarshalAuthorizedKey(publicKey),
	}, nil
}

// PairFromED25519 marshals a valid pair of openssh pem for ED25519 keypairs.
func PairFromED25519(public ed25519.PublicKey, private ed25519.PrivateKey) (*Pair, error) {
	// see https://github.com/golang/crypto/blob/7f63de1d35b0f77fa2b9faea3e7deb402a2383c8/ssh/keys.go#L1273-L1443
	// and https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.key
	// Also the right cipher block size for padding could be found here
	// https://github.com/openssh/openssh-portable/blob/eba523f0a130f1cce829e6aecdcefa841f526a1a/cipher.c#L112
	const noneCipherBlockSize = 8

	key := struct {
		Pub     []byte
		Priv    []byte
		Comment string
		Pad     []byte `ssh:"rest"`
	}{
		Pub:  public,
		Priv: private,
	}
	keyBytes := ssh.Marshal(key)

	pk1 := struct {
		Check1  uint32
		Check2  uint32
		Keytype string
		Rest    []byte `ssh:"rest"`
	}{
		Keytype: ssh.KeyAlgoED25519,
		Rest:    keyBytes,
	}
	pk1Bytes := ssh.Marshal(pk1)

	pubk1 := struct {
		Keytype string
		Key     []byte
	}{
		Keytype: ssh.KeyAlgoED25519,
		Key:     public,
	}
	pubk1Bytes := ssh.Marshal(pubk1)

	padLen := noneCipherBlockSize - (len(pk1Bytes) % noneCipherBlockSize)
	for i := 1; i <= padLen; i++ {
		pk1Bytes = append(pk1Bytes, byte(i))
	}

	k := struct {
		CipherName   string
		KdfName      string
		KdfOpts      string
		NumKeys      uint32
		PubKey       []byte
		PrivKeyBlock []byte
	}{
		CipherName:   "none",
		KdfName:      "none",
		KdfOpts:      "",
		NumKeys:      1,
		PubKey:       pubk1Bytes,
		PrivKeyBlock: pk1Bytes,
	}

	const opensshV1Magic = "openssh-key-v1\x00"

	privBlk := &pem.Block{
		Type:    "OPENSSH PRIVATE KEY",
		Headers: nil,
		Bytes:   append([]byte(opensshV1Magic), ssh.Marshal(k)...),
	}
	publicKey, err := ssh.NewPublicKey(public)
	if err != nil {
		return nil, err
	}
	return &Pair{
		Private: pem.EncodeToMemory(privBlk),
		Public:  ssh.MarshalAuthorizedKey(publicKey),
	}, nil
}

// PairFromDSA always returns ErrDSANotSupported.
//
// Deprecated: DSA is no longer supported by golang.org/x/crypto/ssh. Use
// ECDSA, ED25519 or RSA keys instead.
func PairFromDSA(_ *dsa.PrivateKey) (*Pair, error) {
	return nil, ErrDSANotSupported
}

// GeneratePair generates a Private/Public key pair using algorithm t.
//
// When rand is nil "crypto/rand".Reader will be used.
//
// bits specifies the number of bits in the key to create. For RSA keys, the
// minimum size is 1024 bits and the default is 4096 bits. Generally, 3072 bits
// is considered sufficient. For ECDSA keys, bits determines the key length by
// selecting from one of three elliptic curve sizes: 256, 384 or 521 bits.
// Attempting to use bit lengths other than these three values for ECDSA keys
// will fail. Ed25519 keys have a fixed length and the bits will be ignored.
//
// DSA is no longer supported and returns ErrDSANotSupported.
func GeneratePair(t Algorithm, rand io.Reader, bits int) (*Pair, error) {
	if rand == nil {
		rand = cryptorand.Reader
	}
	switch t {
	case DSA:
		return nil, ErrDSANotSupported
	case ECDSA:
		if bits == 0 {
			bits = 521
		}
		var ecdsakey *ecdsa.PrivateKey
		var err error
		switch bits {
		case 256:
			ecdsakey, err = ecdsa.GenerateKey(elliptic.P256(), rand)
		case 384:
			ecdsakey, err = ecdsa.GenerateKey(elliptic.P384(), rand)
		case 521:
			ecdsakey, err = ecdsa.GenerateKey(elliptic.P521(), rand)
		default:
			ecdsakey, err = nil, ErrInvalidECDSAKeySize
		}
		if err != nil {
			return nil, err
		}
		return NewPair(&ecdsakey.PublicKey, ecdsakey)
	case ED25519:
		publicKey, privateKey, err := ed25519.GenerateKey(rand)
		if err != nil {
			return nil, err
		}
		return PairFromED25519(publicKey, privateKey)
	case RSA:
		if bits == 0 {
			bits = 4096
		}
		if bits < 1024 {
			return nil, ErrInvalidRSAKeySize
		}
		rsakey, err := rsa.GenerateKey(rand, bits)
		if err != nil {
			return nil, err
		}
		return NewPair(&rsakey.PublicKey, rsakey)
	default:
		return nil, ErrUnknownAlgorithm
	}
}
