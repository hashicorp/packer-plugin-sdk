// Code generated by "packer-sdc mapstructure-to-hcl2 "; DO NOT EDIT.

package config

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatKeyValue is an auto-generated flat version of KeyValue.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatKeyValue struct {
	Key   *string `cty:"key" hcl:"key"`
	Value *string `cty:"value" hcl:"value"`
}

// FlatMapstructure returns a new FlatKeyValue.
// FlatKeyValue is an auto-generated flat version of KeyValue.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*KeyValue) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatKeyValue)
}

// HCL2Spec returns the hcl spec of a KeyValue.
// This spec is used by HCL to read the fields of KeyValue.
// The decoded values from this spec will then be applied to a FlatKeyValue.
func (*FlatKeyValue) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"key":   &hcldec.AttrSpec{Name: "key", Type: cty.String, Required: false},
		"value": &hcldec.AttrSpec{Name: "value", Type: cty.String, Required: false},
	}
	return s
}

// FlatKeyValueFilter is an auto-generated flat version of KeyValueFilter.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatKeyValueFilter struct {
	Filters map[string]string `cty:"filters" hcl:"filters"`
	Filter  []FlatKeyValue    `cty:"filter" hcl:"filter"`
}

// FlatMapstructure returns a new FlatKeyValueFilter.
// FlatKeyValueFilter is an auto-generated flat version of KeyValueFilter.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*KeyValueFilter) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatKeyValueFilter)
}

// HCL2Spec returns the hcl spec of a KeyValueFilter.
// This spec is used by HCL to read the fields of KeyValueFilter.
// The decoded values from this spec will then be applied to a FlatKeyValueFilter.
func (*FlatKeyValueFilter) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"filters": &hcldec.AttrSpec{Name: "filters", Type: cty.Map(cty.String), Required: false},
		"filter":  &hcldec.BlockListSpec{TypeName: "filter", Nested: hcldec.ObjectSpec((*FlatKeyValue)(nil).HCL2Spec())},
	}
	return s
}

// FlatNameValue is an auto-generated flat version of NameValue.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatNameValue struct {
	Name  *string `cty:"name" hcl:"name"`
	Value *string `cty:"value" hcl:"value"`
}

// FlatMapstructure returns a new FlatNameValue.
// FlatNameValue is an auto-generated flat version of NameValue.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*NameValue) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatNameValue)
}

// HCL2Spec returns the hcl spec of a NameValue.
// This spec is used by HCL to read the fields of NameValue.
// The decoded values from this spec will then be applied to a FlatNameValue.
func (*FlatNameValue) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"name":  &hcldec.AttrSpec{Name: "name", Type: cty.String, Required: false},
		"value": &hcldec.AttrSpec{Name: "value", Type: cty.String, Required: false},
	}
	return s
}

// FlatNameValueFilter is an auto-generated flat version of NameValueFilter.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatNameValueFilter struct {
	Filters map[string]string `cty:"filters" hcl:"filters"`
	Filter  []FlatNameValue   `cty:"filter" hcl:"filter"`
}

// FlatMapstructure returns a new FlatNameValueFilter.
// FlatNameValueFilter is an auto-generated flat version of NameValueFilter.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*NameValueFilter) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatNameValueFilter)
}

// HCL2Spec returns the hcl spec of a NameValueFilter.
// This spec is used by HCL to read the fields of NameValueFilter.
// The decoded values from this spec will then be applied to a FlatNameValueFilter.
func (*FlatNameValueFilter) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"filters": &hcldec.AttrSpec{Name: "filters", Type: cty.Map(cty.String), Required: false},
		"filter":  &hcldec.BlockListSpec{TypeName: "filter", Nested: hcldec.ObjectSpec((*FlatNameValue)(nil).HCL2Spec())},
	}
	return s
}
