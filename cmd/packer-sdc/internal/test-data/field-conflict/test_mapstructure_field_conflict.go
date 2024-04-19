// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package test

type NestedOne struct {
	Arg int `mapstructure:"test"`
}

type NestedTwo struct {
	Arg int `mapstructure:"test"`
}

type Config struct {
	NestedOne `mapstructure:",squash"`
	NestedTwo `mapstructure:",squash"`
}
