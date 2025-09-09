// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package secretsmanager provide methods to get data from
// AWS Secret Manager
package secretsmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManagerAPI defines the interface for AWS Secrets Manager operations
type SecretsManagerAPI interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// Client represents an AWS Secrets Manager client
type Client struct {
	config *AWSConfig
	api    SecretsManagerAPI
}

// New creates an AWS Secrets Manager Client
func New(config *AWSConfig) *Client {
	c := &Client{
		config: config,
	}

	cfg := c.loadConfig(config)
	c.api = secretsmanager.NewFromConfig(cfg)
	return c
}

func (c *Client) loadConfig(awsConfig *AWSConfig) aws.Config {
	ctx := context.Background()

	var opts []func(*config.LoadOptions) error

	if awsConfig.Region != "" {
		opts = append(opts, config.WithRegion(awsConfig.Region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		// In v1, session.Must would panic on error, maintaining same behavior
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}

	return cfg
}

// GetSecret return an AWS Secret Manager secret
// in plain text from a given secret name
func (c *Client) GetSecret(spec *SecretSpec, raw bool) (string, error) {
	ctx := context.Background()

	params := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(spec.Name),
		VersionStage: aws.String("AWSCURRENT"),
	}

	resp, err := c.api.GetSecretValue(ctx, params)
	if err != nil {
		return "", err
	}

	if resp == nil || resp.SecretString == nil {
		return "", errors.New("Secret is not string")
	}

	secret := SecretString{
		Name:         *resp.Name,
		SecretString: *resp.SecretString,
	}
	value, err := getSecretValue(&secret, spec, raw)
	if err != nil {
		return "", err
	}

	return value, nil
}

func getSecretValue(s *SecretString, spec *SecretSpec, raw bool) (string, error) {
	var secretValue map[string]interface{}
	blob := []byte(s.SecretString)

	//For those plaintext secrets just return the value or if raw is requested
	if !json.Valid(blob) || raw {
		return s.SecretString, nil
	}

	err := json.Unmarshal(blob, &secretValue)
	if err != nil {
		return "", err
	}

	// If key is not set and secret has multiple keys, return error
	if spec.Key == "" && len(secretValue) > 1 {
		return "", errors.New("Secret has multiple values and no key was set")
	}

	if spec.Key == "" {
		for _, v := range secretValue {
			return getStringSecretValue(v)
		}
	}

	if v, ok := secretValue[spec.Key]; ok {
		return getStringSecretValue(v)
	}

	return "", fmt.Errorf("No secret found for key %q", spec.Key)
}

func getStringSecretValue(v interface{}) (string, error) {
	switch valueType := v.(type) {
	case string:
		return valueType, nil
	case float64:
		return strconv.FormatFloat(valueType, 'f', 0, 64), nil
	default:
		return "", fmt.Errorf("Unsupported secret value type: %T", valueType)
	}
}
