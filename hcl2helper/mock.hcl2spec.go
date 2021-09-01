// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package hcl2helper

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatMapOfNestedMockConfig is an auto-generated flat version of MapOfNestedMockConfig.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatMapOfNestedMockConfig struct {
	Nested map[string]FlatSimpleNestedMockConfig `mapstructure:"nested" cty:"nested" hcl:"nested"`
}

// FlatMapstructure returns a new FlatMapOfNestedMockConfig.
// FlatMapOfNestedMockConfig is an auto-generated flat version of MapOfNestedMockConfig.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*MapOfNestedMockConfig) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatMapOfNestedMockConfig)
}

// HCL2Spec returns the hcl spec of a MapOfNestedMockConfig.
// This spec is used by HCL to read the fields of MapOfNestedMockConfig.
// The decoded values from this spec will then be applied to a FlatMapOfNestedMockConfig.
func (*FlatMapOfNestedMockConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		// This one requires the map keys pre-defined. Without it skips whatever the key and return null map. See next definition.
		"nested": &hcldec.BlockMapSpec{TypeName: "nested", Nested: hcldec.ObjectSpec((*FlatSimpleNestedMockConfig)(nil).HCL2Spec())},

		// This one requires the map keys pre-defined. Not possible as this is dynamic like 'us-east-1'
		//"nested": &hcldec.BlockMapSpec{TypeName: "nested", LabelNames: []string{"mock"},Nested: hcldec.ObjectSpec((*FlatNestedMockConfig)(nil).HCL2Spec())},

		// This one requires the type of the map which is no possible. Empty object will always be empty in this case.
		//"nested": &hcldec.AttrSpec{Name: "nested", Type: cty.Map(cty.EmptyObject), Required: false},
	}
	return s
}

// FlatMockConfig is an auto-generated flat version of MockConfig.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatMockConfig struct {
	NotSquashed          *string                `mapstructure:"not_squashed" cty:"not_squashed" hcl:"not_squashed"`
	String               *string                `mapstructure:"string" cty:"string" hcl:"string"`
	Int                  *int                   `mapstructure:"int" cty:"int" hcl:"int"`
	Int64                *int64                 `mapstructure:"int64" cty:"int64" hcl:"int64"`
	Bool                 *bool                  `mapstructure:"bool" cty:"bool" hcl:"bool"`
	Trilean              *bool                  `mapstructure:"trilean" cty:"trilean" hcl:"trilean"`
	Duration             *string                `mapstructure:"duration" cty:"duration" hcl:"duration"`
	MapStringString      map[string]string      `mapstructure:"map_string_string" cty:"map_string_string" hcl:"map_string_string"`
	SliceString          []string               `mapstructure:"slice_string" cty:"slice_string" hcl:"slice_string"`
	SliceSliceString     [][]string             `mapstructure:"slice_slice_string" cty:"slice_slice_string" hcl:"slice_slice_string"`
	NamedMapStringString NamedMapStringString   `mapstructure:"named_map_string_string" cty:"named_map_string_string" hcl:"named_map_string_string"`
	NamedString          *NamedString           `mapstructure:"named_string" cty:"named_string" hcl:"named_string"`
	Tags                 []FlatMockTag          `mapstructure:"tag" cty:"tag" hcl:"tag"`
	Datasource           *string                `mapstructure:"data_source" cty:"data_source" hcl:"data_source"`
	Nested               *FlatNestedMockConfig  `mapstructure:"nested" cty:"nested" hcl:"nested"`
	NestedSlice          []FlatNestedMockConfig `mapstructure:"nested_slice" cty:"nested_slice" hcl:"nested_slice"`
}

// FlatMapstructure returns a new FlatMockConfig.
// FlatMockConfig is an auto-generated flat version of MockConfig.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*MockConfig) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatMockConfig)
}

// HCL2Spec returns the hcl spec of a MockConfig.
// This spec is used by HCL to read the fields of MockConfig.
// The decoded values from this spec will then be applied to a FlatMockConfig.
func (*FlatMockConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"not_squashed":            &hcldec.AttrSpec{Name: "not_squashed", Type: cty.String, Required: false},
		"string":                  &hcldec.AttrSpec{Name: "string", Type: cty.String, Required: false},
		"int":                     &hcldec.AttrSpec{Name: "int", Type: cty.Number, Required: false},
		"int64":                   &hcldec.AttrSpec{Name: "int64", Type: cty.Number, Required: false},
		"bool":                    &hcldec.AttrSpec{Name: "bool", Type: cty.Bool, Required: false},
		"trilean":                 &hcldec.AttrSpec{Name: "trilean", Type: cty.Bool, Required: false},
		"duration":                &hcldec.AttrSpec{Name: "duration", Type: cty.String, Required: false},
		"map_string_string":       &hcldec.AttrSpec{Name: "map_string_string", Type: cty.Map(cty.String), Required: false},
		"slice_string":            &hcldec.AttrSpec{Name: "slice_string", Type: cty.List(cty.String), Required: false},
		"slice_slice_string":      &hcldec.AttrSpec{Name: "slice_slice_string", Type: cty.List(cty.List(cty.String)), Required: false},
		"named_map_string_string": &hcldec.AttrSpec{Name: "named_map_string_string", Type: cty.Map(cty.String), Required: false},
		"named_string":            &hcldec.AttrSpec{Name: "named_string", Type: cty.String, Required: false},
		"tag":                     &hcldec.BlockListSpec{TypeName: "tag", Nested: hcldec.ObjectSpec((*FlatMockTag)(nil).HCL2Spec())},
		"data_source":             &hcldec.AttrSpec{Name: "data_source", Type: cty.String, Required: false},
		"nested":                  &hcldec.BlockSpec{TypeName: "nested", Nested: hcldec.ObjectSpec((*FlatNestedMockConfig)(nil).HCL2Spec())},
		"nested_slice":            &hcldec.BlockListSpec{TypeName: "nested_slice", Nested: hcldec.ObjectSpec((*FlatNestedMockConfig)(nil).HCL2Spec())},
	}
	return s
}

// FlatMockTag is an auto-generated flat version of MockTag.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatMockTag struct {
	Key   *string `mapstructure:"key" cty:"key" hcl:"key"`
	Value *string `mapstructure:"value" cty:"value" hcl:"value"`
}

// FlatMapstructure returns a new FlatMockTag.
// FlatMockTag is an auto-generated flat version of MockTag.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*MockTag) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatMockTag)
}

// HCL2Spec returns the hcl spec of a MockTag.
// This spec is used by HCL to read the fields of MockTag.
// The decoded values from this spec will then be applied to a FlatMockTag.
func (*FlatMockTag) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"key":   &hcldec.AttrSpec{Name: "key", Type: cty.String, Required: false},
		"value": &hcldec.AttrSpec{Name: "value", Type: cty.String, Required: false},
	}
	return s
}

// FlatNestedMockConfig is an auto-generated flat version of NestedMockConfig.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatNestedMockConfig struct {
	String               *string              `mapstructure:"string" cty:"string" hcl:"string"`
	Int                  *int                 `mapstructure:"int" cty:"int" hcl:"int"`
	Int64                *int64               `mapstructure:"int64" cty:"int64" hcl:"int64"`
	Bool                 *bool                `mapstructure:"bool" cty:"bool" hcl:"bool"`
	Trilean              *bool                `mapstructure:"trilean" cty:"trilean" hcl:"trilean"`
	Duration             *string              `mapstructure:"duration" cty:"duration" hcl:"duration"`
	MapStringString      map[string]string    `mapstructure:"map_string_string" cty:"map_string_string" hcl:"map_string_string"`
	SliceString          []string             `mapstructure:"slice_string" cty:"slice_string" hcl:"slice_string"`
	SliceSliceString     [][]string           `mapstructure:"slice_slice_string" cty:"slice_slice_string" hcl:"slice_slice_string"`
	NamedMapStringString NamedMapStringString `mapstructure:"named_map_string_string" cty:"named_map_string_string" hcl:"named_map_string_string"`
	NamedString          *NamedString         `mapstructure:"named_string" cty:"named_string" hcl:"named_string"`
	Tags                 []FlatMockTag        `mapstructure:"tag" cty:"tag" hcl:"tag"`
	Datasource           *string              `mapstructure:"data_source" cty:"data_source" hcl:"data_source"`
}

// FlatMapstructure returns a new FlatNestedMockConfig.
// FlatNestedMockConfig is an auto-generated flat version of NestedMockConfig.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*NestedMockConfig) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatNestedMockConfig)
}

// HCL2Spec returns the hcl spec of a NestedMockConfig.
// This spec is used by HCL to read the fields of NestedMockConfig.
// The decoded values from this spec will then be applied to a FlatNestedMockConfig.
func (*FlatNestedMockConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"string":                  &hcldec.AttrSpec{Name: "string", Type: cty.String, Required: false},
		"int":                     &hcldec.AttrSpec{Name: "int", Type: cty.Number, Required: false},
		"int64":                   &hcldec.AttrSpec{Name: "int64", Type: cty.Number, Required: false},
		"bool":                    &hcldec.AttrSpec{Name: "bool", Type: cty.Bool, Required: false},
		"trilean":                 &hcldec.AttrSpec{Name: "trilean", Type: cty.Bool, Required: false},
		"duration":                &hcldec.AttrSpec{Name: "duration", Type: cty.String, Required: false},
		"map_string_string":       &hcldec.AttrSpec{Name: "map_string_string", Type: cty.Map(cty.String), Required: false},
		"slice_string":            &hcldec.AttrSpec{Name: "slice_string", Type: cty.List(cty.String), Required: false},
		"slice_slice_string":      &hcldec.AttrSpec{Name: "slice_slice_string", Type: cty.List(cty.List(cty.String)), Required: false},
		"named_map_string_string": &hcldec.AttrSpec{Name: "named_map_string_string", Type: cty.Map(cty.String), Required: false},
		"named_string":            &hcldec.AttrSpec{Name: "named_string", Type: cty.String, Required: false},
		"tag":                     &hcldec.BlockListSpec{TypeName: "tag", Nested: hcldec.ObjectSpec((*FlatMockTag)(nil).HCL2Spec())},
		"data_source":             &hcldec.AttrSpec{Name: "data_source", Type: cty.String, Required: false},
	}
	return s
}

// FlatSimpleNestedMockConfig is an auto-generated flat version of SimpleNestedMockConfig.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatSimpleNestedMockConfig struct {
	String *string `mapstructure:"string" cty:"string" hcl:"string"`
}

// FlatMapstructure returns a new FlatSimpleNestedMockConfig.
// FlatSimpleNestedMockConfig is an auto-generated flat version of SimpleNestedMockConfig.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*SimpleNestedMockConfig) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatSimpleNestedMockConfig)
}

// HCL2Spec returns the hcl spec of a SimpleNestedMockConfig.
// This spec is used by HCL to read the fields of SimpleNestedMockConfig.
// The decoded values from this spec will then be applied to a FlatSimpleNestedMockConfig.
func (*FlatSimpleNestedMockConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"string": &hcldec.AttrSpec{Name: "string", Type: cty.String, Required: false},
	}
	return s
}
