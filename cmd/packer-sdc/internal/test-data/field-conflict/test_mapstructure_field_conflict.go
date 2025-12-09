// Copyright IBM Corp. 2013, 2025
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
