// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communicator

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/masterzen/winrm"
	"golang.org/x/crypto/ssh"
)

func testConfig() *Config {
	return &Config{
		SSH: SSH{
			SSHUsername: "root",
		},
	}
}

func TestConfigType(t *testing.T) {
	c := testConfig()
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.Type != "ssh" {
		t.Fatalf("bad: %#v", c)
	}
}

func TestConfig_none(t *testing.T) {
	c := &Config{Type: "none"}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}
}

func TestConfig_badtype(t *testing.T) {
	c := &Config{Type: "foo"}
	if err := c.Prepare(testContext(t)); len(err) != 1 {
		t.Fatalf("bad: %#v", err)
	}
}

func TestConfig_winrm_noport(t *testing.T) {
	c := &Config{
		Type: "winrm",
		WinRM: WinRM{
			WinRMUser: "admin",
		},
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMPort != 5985 {
		t.Fatalf("WinRMPort doesn't match default port 5985 when SSL is not enabled and no port is specified.")
	}

}

func TestConfig_winrm_noport_ssl(t *testing.T) {
	c := &Config{
		Type: "winrm",
		WinRM: WinRM{
			WinRMUser:   "admin",
			WinRMUseSSL: true,
		},
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMPort != 5986 {
		t.Fatalf("WinRMPort doesn't match default port 5986 when SSL is enabled and no port is specified.")
	}

}

func TestConfig_winrm_port(t *testing.T) {
	c := &Config{
		Type: "winrm",
		WinRM: WinRM{
			WinRMUser: "admin",
			WinRMPort: 5509,
		},
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMPort != 5509 {
		t.Fatalf("WinRMPort doesn't match custom port 5509 when SSL is not enabled.")
	}

}

func TestConfig_winrm_port_ssl(t *testing.T) {
	c := &Config{
		Type: "winrm",
		WinRM: WinRM{
			WinRMUser:   "admin",
			WinRMPort:   5510,
			WinRMUseSSL: true,
		},
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMPort != 5510 {
		t.Fatalf("WinRMPort doesn't match custom port 5510 when SSL is enabled.")
	}

}

func TestConfig_winrm_use_ntlm(t *testing.T) {
	c := &Config{
		Type: "winrm",
		WinRM: WinRM{
			WinRMUser:    "admin",
			WinRMUseNTLM: true,
		},
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}

	if c.WinRMTransportDecorator == nil {
		t.Fatalf("WinRMTransportDecorator not set.")
	}

	expected := &winrm.ClientNTLM{}
	actual := c.WinRMTransportDecorator()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("WinRMTransportDecorator isn't ClientNTLM.")
	}

}

// generateSSHPrivateKey generates a new RSA SSH private key for use in tests
//
// It returns the path in which the key was created.
// Removing the key after testing is the caller's responsibility.
func generateSSHPrivateKey() (path string, signer ssh.Signer, err error) {
	pk, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		err = fmt.Errorf("failed to generate key: %s", err)
		return
	}

	sshKeyFile, err := os.CreateTemp("", "")
	if err != nil {
		err = fmt.Errorf("failed to open a temp file: %s", err)
		return
	}

	defer sshKeyFile.Close()

	path = sshKeyFile.Name()

	rawPkey := x509.MarshalPKCS1PrivateKey(pk)

	err = pem.Encode(sshKeyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: rawPkey,
	})
	if err != nil {
		err = fmt.Errorf("failed to encode to PEM: %s", err)
		return
	}

	signer, err = ssh.NewSignerFromKey(pk)
	if err != nil {
		err = fmt.Errorf("failed to create SSH signer: %s", err)
		return
	}

	return
}

// generateSSHKeys generates a new SSH key, the CA key and a cert linked to the SSH key for use in tests
//
// It returns the paths in which the keys and cert were created.
// Removing the keys and certs after testing is the caller's responsibility.
func generateSSHKeys() (
	privKeyPath string,
	certKeyPath string,
	certPath string,
	err error,
) {
	var sshPrivKey, certSSHKey ssh.Signer

	privKeyPath, sshPrivKey, err = generateSSHPrivateKey()
	if err != nil {
		err = fmt.Errorf("failed to generate private key: %s", err)
		return
	}

	certKeyPath, certSSHKey, err = generateSSHPrivateKey()
	if err != nil {
		err = fmt.Errorf("failed to generate CA private key: %s", err)
		return
	}

	cert := &ssh.Certificate{
		CertType:        ssh.HostCert,
		Key:             sshPrivKey.PublicKey(),
		ValidAfter:      0,
		ValidBefore:     ssh.CertTimeInfinity,
		KeyId:           "TestSSHCert",
		ValidPrincipals: []string{"authority.example.com"},
	}

	certFile, err := os.CreateTemp("", "")
	if err != nil {
		err = fmt.Errorf("failed to create cert file: %s", err)
		return
	}
	defer certFile.Close()

	certPath = certFile.Name()

	err = cert.SignCert(rand.Reader, certSSHKey)
	if err != nil {
		err = fmt.Errorf("failed to sign cert: %s", err)
		return
	}

	rawCert := ssh.MarshalAuthorizedKey(cert)

	_, err = certFile.Write(rawCert)
	if err != nil {
		err = fmt.Errorf("failed to write marshalled certificate: %s", err)
	}

	return
}

func TestSSHBastion(t *testing.T) {
	privKeyPath, certKeyPath, certPath, err := generateSSHKeys()
	if err != nil {
		t.Fatalf("failed to generate SSH keys and certificates: %s", err)
	}

	defer func() {
		os.Remove(privKeyPath)
		os.Remove(certKeyPath)
		os.Remove(certPath)
	}()

	t.Logf("generated private key (%q), CA key (%q), certificate (%q)", privKeyPath, certKeyPath, certPath)

	bastionPrivKeyPath, bastionCertKeyPath, bastionCertPath, err := generateSSHKeys()
	if err != nil {
		t.Fatalf("failed to generate bastion SSH keys and certificates: %s", err)
	}

	defer func() {
		os.Remove(bastionPrivKeyPath)
		os.Remove(bastionCertKeyPath)
		os.Remove(bastionCertPath)
	}()

	t.Logf("generated bastion private key (%q), CA key (%q), certificate (%q)", bastionPrivKeyPath, bastionCertKeyPath, bastionCertPath)

	testcases := []struct {
		name           string
		config         *Config
		expectedConfig *Config
		expectError    bool
	}{
		{
			"OK - with host and password",
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:        "root",
					SSHBastionHost:     "mybastionhost.company.com",
					SSHBastionPassword: "test",
				},
			},
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:           "root",
					SSHBastionHost:        "mybastionhost.company.com",
					SSHBastionPassword:    "test",
					SSHPort:               22,
					SSHTimeout:            time.Minute * 5,
					SSHFileTransferMethod: "scp",
					SSHKeepAliveInterval:  time.Second * 5,
					SSHHandshakeAttempts:  10,
					SSHBastionPort:        22,
				},
			},
			false,
		},
		{
			"OK - bastion config with bastion SSH private key",
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:              "root",
					SSHBastionHost:           "my.bastion",
					SSHBastionPrivateKeyFile: bastionPrivKeyPath,
				},
			},
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:              "root",
					SSHBastionHost:           "my.bastion",
					SSHBastionPrivateKeyFile: bastionPrivKeyPath,
					SSHPort:                  22,
					SSHTimeout:               time.Minute * 5,
					SSHFileTransferMethod:    "scp",
					SSHKeepAliveInterval:     time.Second * 5,
					SSHHandshakeAttempts:     10,
					SSHBastionPort:           22,
				},
			},
			false,
		},
		{
			"OK - bastion config with SSH private key, bastion key should be the same as SSH key",
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:       "root",
					SSHBastionHost:    "my.bastion",
					SSHPrivateKeyFile: privKeyPath,
				},
			},
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:              "root",
					SSHBastionHost:           "my.bastion",
					SSHBastionPrivateKeyFile: privKeyPath,
					SSHPort:                  22,
					SSHTimeout:               time.Minute * 5,
					SSHFileTransferMethod:    "scp",
					SSHKeepAliveInterval:     time.Second * 5,
					SSHHandshakeAttempts:     10,
					SSHBastionPort:           22,
					SSHPrivateKeyFile:        privKeyPath,
				},
			},
			false,
		},
		{
			"OK - bastion config with SSH private key and cert, bastion should have both set",
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:        "root",
					SSHBastionHost:     "my.bastion",
					SSHPrivateKeyFile:  privKeyPath,
					SSHCertificateFile: certPath,
				},
			},
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:               "root",
					SSHBastionHost:            "my.bastion",
					SSHBastionPrivateKeyFile:  privKeyPath,
					SSHBastionCertificateFile: certPath,
					SSHPort:                   22,
					SSHTimeout:                time.Minute * 5,
					SSHFileTransferMethod:     "scp",
					SSHKeepAliveInterval:      time.Second * 5,
					SSHHandshakeAttempts:      10,
					SSHBastionPort:            22,
					SSHPrivateKeyFile:         privKeyPath,
					SSHCertificateFile:        certPath,
				},
			},
			false,
		},
		{
			"OK - bastion config with SSH private key and cert, and a bastion private key, bastion cert should not be set",
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:              "root",
					SSHBastionHost:           "my.bastion",
					SSHBastionPrivateKeyFile: bastionPrivKeyPath,
					SSHPrivateKeyFile:        privKeyPath,
					SSHCertificateFile:       certPath,
				},
			},
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:              "root",
					SSHBastionHost:           "my.bastion",
					SSHBastionPrivateKeyFile: bastionPrivKeyPath,
					SSHPort:                  22,
					SSHTimeout:               time.Minute * 5,
					SSHFileTransferMethod:    "scp",
					SSHKeepAliveInterval:     time.Second * 5,
					SSHHandshakeAttempts:     10,
					SSHBastionPort:           22,
					SSHPrivateKeyFile:        privKeyPath,
					SSHCertificateFile:       certPath,
				},
			},
			false,
		},
		{
			"OK - bastion config with SSH private key and cert, and a bastion private key and cert",
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:               "root",
					SSHBastionHost:            "my.bastion",
					SSHBastionPrivateKeyFile:  bastionPrivKeyPath,
					SSHBastionCertificateFile: bastionCertPath,
					SSHPrivateKeyFile:         privKeyPath,
					SSHCertificateFile:        certPath,
				},
			},
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:               "root",
					SSHBastionHost:            "my.bastion",
					SSHBastionPrivateKeyFile:  bastionPrivKeyPath,
					SSHBastionCertificateFile: bastionCertPath,
					SSHPort:                   22,
					SSHTimeout:                time.Minute * 5,
					SSHFileTransferMethod:     "scp",
					SSHKeepAliveInterval:      time.Second * 5,
					SSHHandshakeAttempts:      10,
					SSHBastionPort:            22,
					SSHPrivateKeyFile:         privKeyPath,
					SSHCertificateFile:        certPath,
				},
			},
			false,
		},
		{
			"Fail - ssh certificate file specified without an ssh private key file",
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:        "root",
					SSHCertificateFile: certPath,
				},
			},
			nil,
			true,
		},
		{
			"Fail - ssh bastion certificate file specified without an ssh bastion private key file",
			&Config{
				Type: "ssh",
				SSH: SSH{
					SSHUsername:               "root",
					SSHBastionCertificateFile: certPath,
				},
			},
			nil,
			true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.config.Prepare(testContext(t))

			for _, err := range errs {
				t.Logf("%s", err)
			}
			if (len(errs) != 0) != tt.expectError {
				t.Fatalf("Expected %t error, got %d", tt.expectError, len(errs))
			}
			if tt.expectError {
				return
			}

			diff := cmp.Diff(tt.config, tt.expectedConfig)
			if diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestSSHConfigFunc_ciphers(t *testing.T) {
	state := new(multistep.BasicStateBag)

	// No ciphers set
	c := &Config{
		Type: "ssh",
	}

	f := c.SSHConfigFunc()
	sshConfig, _ := f(state)
	if sshConfig.Config.Ciphers != nil {
		t.Fatalf("Shouldn't set SSHCiphers if communicator config option " +
			"ssh_ciphers is unset.")
	}

	// Ciphers are set
	c = &Config{
		Type: "ssh",
		SSH: SSH{
			SSHCiphers: []string{"partycipher"},
		},
	}
	f = c.SSHConfigFunc()
	sshConfig, _ = f(state)
	if sshConfig.Config.Ciphers == nil {
		t.Fatalf("Shouldn't set SSHCiphers if communicator config option " +
			"ssh_ciphers is unset.")
	}
	if sshConfig.Config.Ciphers[0] != "partycipher" {
		t.Fatalf("ssh_ciphers should be a direct passthrough.")
	}
	if c.SSHCertificateFile != "" {
		t.Fatalf("Identity certificate somehow set")
	}
}

func TestSSHConfigFunc_kexAlgos(t *testing.T) {
	state := new(multistep.BasicStateBag)

	// No ciphers set
	c := &Config{
		Type: "ssh",
	}

	f := c.SSHConfigFunc()
	sshConfig, _ := f(state)
	if sshConfig.Config.KeyExchanges != nil {
		t.Fatalf("Shouldn't set KeyExchanges if communicator config option " +
			"ssh_key_exchange_algorithms is unset.")
	}

	// Ciphers are set
	c = &Config{
		Type: "ssh",
		SSH: SSH{
			SSHKEXAlgos: []string{"partyalgo"},
		},
	}
	f = c.SSHConfigFunc()
	sshConfig, _ = f(state)
	if sshConfig.Config.KeyExchanges == nil {
		t.Fatalf("Should set SSHKEXAlgos if communicator config option " +
			"ssh_key_exchange_algorithms is set.")
	}
	if sshConfig.Config.KeyExchanges[0] != "partyalgo" {
		t.Fatalf("ssh_key_exchange_algorithms should be a direct passthrough.")
	}
	if c.SSHCertificateFile != "" {
		t.Fatalf("Identity certificate somehow set")
	}
}

func TestConfig_winrm(t *testing.T) {
	c := &Config{
		Type: "winrm",
		WinRM: WinRM{
			WinRMUser: "admin",
		},
	}
	if err := c.Prepare(testContext(t)); len(err) > 0 {
		t.Fatalf("bad: %#v", err)
	}
}

func testContext(t *testing.T) *interpolate.Context {
	return nil
}
