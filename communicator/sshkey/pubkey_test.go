// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package sshkey

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	privDSA = `-----BEGIN DSA PRIVATE KEY-----
MIIBuwIBAAKBgQD2/wltLNU9YDnye+wEtF4Y6NPvIUtxbXjFbbcyq4voi6nm6Jql
XrkMEQ9uu94ibsZSPfx+fWJmx0jcuk5vBIGYD5hcCmyuZhXerogwRhyYNis+D60d
oWKOGHm51UFdZxB3vS++T/iVBNuf5ZVAcSBXQxm3SYLv2GTSfVq/HS98nwIVAPqN
aymInPvLNGYjei7WDq6YMRkdAoGAQZ16m9JgxHowIdMazOpeBQJdSDL4pMpPc/L9
hox9yF6ZTQUXPdZAFKlOsrTwQuQY3lKn+C9XFJs6htrprGKegCCYlQg7F5nSHpav
/du5KkAYCOlzNArncfX06PRLbqdy6GRvSIWOqvo5HLjauAQxAJJBhcVmvfqkYxQB
nrn22ocCgYEAkuvhJlazyMuaGqK5uTfa+SwjsqPi7IwfXTWyvpdYqdNJhutgmMCp
NN8bmYJ4a/4Cqf9JQNX//QtiQMyD6Hw6/4K93yc0ao4aqdFwrRxX74H+1aaKmisX
q9ayvQ91BWzOPMuXOQ4h2yc8rujxSl/eewkB4vtL81GVDO5bAZk4Rk8CFEPAePES
KAXYNk9l3L3a6L6CvBF8
-----END DSA PRIVATE KEY-----
`
	priv_RSA = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCviEbCEYL3HbR/CfERnnN/5PVyH2FVPlQO1fH0QNmaBXQor2lK
/wKR1jHhpjCyyJqzf81UdiXPQP+dC4TNQtQXyXKAbPwCbOw1AveON429nH72dz3B
HYF9vyxBzYCGPUwtferViPe15/3m68W0sNKlldt+OXASVKNK0objAm0xgQIDAQAB
AoGBAKvyLI4SpV34xUTksjb2JmIUILMoNipQofrebOM9W2tbCEyKd/Q1FYlSbw6B
w5Z+l7xZ5wNjsOny2/I0xGRloGgECuNtYDiEkZxL+DFV3vGjBlTbsgjL8MwVBmYb
BPtH7njbTM2gbLFqHUAoPofOVPxeGIfez7/XFrsOKA1/fm/5AkEA5eRnA/1gp1N8
ntRhJaMcXxPCzQ/NjEzhOvAsR1GtSbKmPVOv1lVJdySUwgHzWycHkNigD9qgHDGJ
10euK2ojEwJBAMN3e7ykKxBfagLmeHwC3GCIxE4f7ckCr43qMme2ga6Q97ik9PKc
e40D3kSI/B4KEPaxIPt63J+7A71GMbUX15sCQBH2q/oK1X+drXI3xDONsEzZnXIq
lvQsmbjiRYN5JWJJ3GWUYXwNBAWUTS9vuZVY0mWGF9PFUJeDY3L3/esUixMCQE03
zHgPWrvTFawjpAFJmAnCvdonHubD2tpzZIo4PS7bMiGNeP7G9sAUgSAOBZtQWrc5
7k+qj7HCTY3eRxDFZHUCQQDEM1tof0JVvGzkV4ieBClsWcmqgZAg4C6qGFGQLubX
QqmcLecpWCOdSy82JBTMUMFjcWF1NX6yTpA2mC6S8LFX
-----END RSA PRIVATE KEY-----
`
	priv_ECDSA = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIHpwAwKnukGd3q39ImG5cdpZV1rVRVKn85jhXgX9ASc5oAoGCCqGSM49
AwEHoUQDQgAEI+QrxZsRkHFbkHCIo1hh8WrikT1bAaHZpTEfVi1w7bQtXErA/2hv
lvVZHAGSJinAHEnO5ZDM0nP5wjUf3OrAqw==
-----END EC PRIVATE KEY-----
`
	priv_ED25519 = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACCttg8eP6zngqQlQB5YvPepOLADodc6mqIwXTxAdr4ZMAAAAJj+rLIa/qyy
GgAAAAtzc2gtZWQyNTUxOQAAACCttg8eP6zngqQlQB5YvPepOLADodc6mqIwXTxAdr4ZMA
AAAECOJsldUUkAGQ1DXzX0eM9/mTEzYZMFmxsiVquCXCCBRq22Dx4/rOeCpCVAHli896k4
sAOh1zqaojBdPEB2vhkwAAAAD3hnZW5AaTEwNzUxNjk2NgECAwQFBg==
-----END OPENSSH PRIVATE KEY-----
`
	pub_DSA     = "ssh-dss AAAAB3NzaC1kc3MAAACBAPb/CW0s1T1gOfJ77AS0Xhjo0+8hS3FteMVttzKri+iLqebomqVeuQwRD2673iJuxlI9/H59YmbHSNy6Tm8EgZgPmFwKbK5mFd6uiDBGHJg2Kz4PrR2hYo4YebnVQV1nEHe9L75P+JUE25/llUBxIFdDGbdJgu/YZNJ9Wr8dL3yfAAAAFQD6jWspiJz7yzRmI3ou1g6umDEZHQAAAIBBnXqb0mDEejAh0xrM6l4FAl1IMvikyk9z8v2GjH3IXplNBRc91kAUqU6ytPBC5BjeUqf4L1cUmzqG2umsYp6AIJiVCDsXmdIelq/927kqQBgI6XM0Cudx9fTo9Etup3LoZG9IhY6q+jkcuNq4BDEAkkGFxWa9+qRjFAGeufbahwAAAIEAkuvhJlazyMuaGqK5uTfa+SwjsqPi7IwfXTWyvpdYqdNJhutgmMCpNN8bmYJ4a/4Cqf9JQNX//QtiQMyD6Hw6/4K93yc0ao4aqdFwrRxX74H+1aaKmisXq9ayvQ91BWzOPMuXOQ4h2yc8rujxSl/eewkB4vtL81GVDO5bAZk4Rk8=\n"
	pub_RSA     = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCviEbCEYL3HbR/CfERnnN/5PVyH2FVPlQO1fH0QNmaBXQor2lK/wKR1jHhpjCyyJqzf81UdiXPQP+dC4TNQtQXyXKAbPwCbOw1AveON429nH72dz3BHYF9vyxBzYCGPUwtferViPe15/3m68W0sNKlldt+OXASVKNK0objAm0xgQ==\n"
	pub_ECDSA   = "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBCPkK8WbEZBxW5BwiKNYYfFq4pE9WwGh2aUxH1YtcO20LVxKwP9ob5b1WRwBkiYpwBxJzuWQzNJz+cI1H9zqwKs=\n"
	pub_ED25519 = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIK22Dx4/rOeCpCVAHli896k4sAOh1zqaojBdPEB2vhkw\n"
)

func TestPublicKeyFromPrivate(t *testing.T) {
	tests := []struct {
		t       Algorithm
		privKey []byte
		pubKey  []byte
	}{
		{DSA,
			[]byte(privDSA),
			[]byte(pub_DSA),
		},
		{RSA,
			[]byte(priv_RSA),
			[]byte(pub_RSA),
		},
		{ECDSA,
			[]byte(priv_ECDSA),
			[]byte(pub_ECDSA),
		},
		{ED25519,
			[]byte(priv_ED25519),
			[]byte(pub_ED25519),
		},
	}
	for _, tt := range tests {
		t.Run(tt.t.String(), func(t *testing.T) {
			got, err := PublicKeyFromPrivate(tt.privKey)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(got, tt.pubKey); diff != "" {
				t.Errorf("wrong PublicKeyFromPrivate(): %s", diff)
			}
		})
	}
}
