// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package secretsmanager

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type mockedSecret struct {
	Resp secretsmanager.GetSecretValueOutput
}

// GetSecretValue return mocked secret value
func (m mockedSecret) GetSecretValue(ctx context.Context, in *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return &m.Resp, nil
}

func TestGetSecret(t *testing.T) {
	testCases := []struct {
		description string
		arg         *SecretSpec
		raw         bool
		mock        secretsmanager.GetSecretValueOutput
		want        string
		ok          bool
	}{
		{
			description: "input has valid secret name, secret has single key",
			arg:         &SecretSpec{Name: "test/secret"},
			raw:         false,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test/secret"),
				SecretString: aws.String(`{"key": "test"}`),
			},
			want: "test",
			ok:   true,
		},
		{
			description: "input has valid secret name and key, secret has single key",
			arg: &SecretSpec{
				Name: "test/secret",
				Key:  "key",
			},
			raw: false,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test/secret"),
				SecretString: aws.String(`{"key": "test"}`),
			},
			want: "test",
			ok:   true,
		},
		{
			description: "input has valid secret name and key, secret has multiple keys",
			arg: &SecretSpec{
				Name: "test/secret",
				Key:  "second_key",
			},
			raw: false,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test/secret"),
				SecretString: aws.String(`{"first_key": "first_val", "second_key": "second_val"}`),
			},
			want: "second_val",
			ok:   true,
		},
		{
			description: "input has valid secret name and no key, secret has multiple keys",
			arg: &SecretSpec{
				Name: "test/secret",
			},
			raw: false,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test/secret"),
				SecretString: aws.String(`{"first_key": "first_val", "second_key": "second_val"}`),
			},
			ok: false,
		},
		{
			description: "input has valid secret name and invalid key, secret has single key",
			arg: &SecretSpec{
				Name: "test/secret",
				Key:  "nonexistent",
			},
			raw: false,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test/secret"),
				SecretString: aws.String(`{"key": "test"}`),
			},
			ok: false,
		},
		{
			description: "input has valid secret name and invalid key, secret has multiple keys",
			arg: &SecretSpec{
				Name: "test/secret",
				Key:  "nonexistent",
			},
			raw: false,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test/secret"),
				SecretString: aws.String(`{"first_key": "first_val", "second_key": "second_val"}`),
			},
			ok: false,
		},
		{
			description: "input has secret and key, secret is empty",
			arg: &SecretSpec{
				Name: "test/secret",
				Key:  "nonexistent",
			},
			raw:  false,
			mock: secretsmanager.GetSecretValueOutput{},
			ok:   false,
		},
		{
			description: "input has secret stored as plaintext",
			arg: &SecretSpec{
				Name: "test",
			},
			raw: false,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test"),
				SecretString: aws.String("ThisIsThePassword"),
			},
			want: "ThisIsThePassword",
			ok:   true,
		},
		{
			description: "input as secret stored with 'String: int' value",
			arg:         &SecretSpec{Name: "test"},
			raw:         false,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test"),
				SecretString: aws.String(`{"port": 5432}`),
			},
			want: "5432",
			ok:   true,
		},
		{
			description: "input as secret stored as json object, returned as json",
			arg:         &SecretSpec{Name: "test"},
			raw:         true,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test"),
				SecretString: aws.String(`{"foo":{"bar":"baz"}}`),
			},
			want: `{"foo":{"bar":"baz"}}`,
			ok:   true,
		},
		{
			description: "input as secret stored as json with object, fails without raw",
			arg:         &SecretSpec{Name: "test"},
			raw:         false,
			mock: secretsmanager.GetSecretValueOutput{
				Name:         aws.String("test"),
				SecretString: aws.String(`{"foo":{"bar":"baz"}}`),
			},
			ok: false,
		},
	}

	for _, test := range testCases {
		c := &Client{
			api: mockedSecret{Resp: test.mock},
		}
		got, err := c.GetSecret(test.arg, test.raw)
		if test.ok {
			if got != test.want {
				t.Fatalf("want %v, got %v, error %v, using arg %v", test.want, got, err, test.arg)
			}
		}
		if !test.ok {
			if err == nil {
				t.Fatalf("error expected but got %q, using arg %v", err, test.arg)
			}
		}
		t.Logf("arg (%v), want %v, got %v, err %v", test.arg, test.want, got, err)
	}
}
