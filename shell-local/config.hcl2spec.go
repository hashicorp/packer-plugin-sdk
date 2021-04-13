// Code generated by "packer-sdc mapstructure-to-hcl2 "; DO NOT EDIT.

package shell_local

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	PackerBuildName     *string           `mapstructure:"packer_build_name" cty:"packer_build_name" hcl:"packer_build_name"`
	PackerBuilderType   *string           `mapstructure:"packer_builder_type" cty:"packer_builder_type" hcl:"packer_builder_type"`
	PackerCoreVersion   *string           `mapstructure:"packer_core_version" cty:"packer_core_version" hcl:"packer_core_version"`
	PackerDebug         *bool             `mapstructure:"packer_debug" cty:"packer_debug" hcl:"packer_debug"`
	PackerForce         *bool             `mapstructure:"packer_force" cty:"packer_force" hcl:"packer_force"`
	PackerOnError       *string           `mapstructure:"packer_on_error" cty:"packer_on_error" hcl:"packer_on_error"`
	PackerUserVars      map[string]string `mapstructure:"packer_user_variables" cty:"packer_user_variables" hcl:"packer_user_variables"`
	PackerSensitiveVars []string          `mapstructure:"packer_sensitive_variables" cty:"packer_sensitive_variables" hcl:"packer_sensitive_variables"`
	Inline              []string          `cty:"inline" hcl:"inline"`
	Script              *string           `cty:"script" hcl:"script"`
	Scripts             []string          `cty:"scripts" hcl:"scripts"`
	ValidExitCodes      []int             `mapstructure:"valid_exit_codes" cty:"valid_exit_codes" hcl:"valid_exit_codes"`
	Vars                []string          `mapstructure:"environment_vars" cty:"environment_vars" hcl:"environment_vars"`
	EnvVarFormat        *string           `mapstructure:"env_var_format" cty:"env_var_format" hcl:"env_var_format"`
	Command             *string           `cty:"command" hcl:"command"`
	ExecuteCommand      []string          `mapstructure:"execute_command" cty:"execute_command" hcl:"execute_command"`
	InlineShebang       *string           `mapstructure:"inline_shebang" cty:"inline_shebang" hcl:"inline_shebang"`
	OnlyOn              []string          `mapstructure:"only_on" cty:"only_on" hcl:"only_on"`
	TempfileExtension   *string           `mapstructure:"tempfile_extension" cty:"tempfile_extension" hcl:"tempfile_extension"`
	UseLinuxPathing     *bool             `mapstructure:"use_linux_pathing" cty:"use_linux_pathing" hcl:"use_linux_pathing"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"packer_build_name":          &hcldec.AttrSpec{Name: "packer_build_name", Type: cty.String, Required: false},
		"packer_builder_type":        &hcldec.AttrSpec{Name: "packer_builder_type", Type: cty.String, Required: false},
		"packer_core_version":        &hcldec.AttrSpec{Name: "packer_core_version", Type: cty.String, Required: false},
		"packer_debug":               &hcldec.AttrSpec{Name: "packer_debug", Type: cty.Bool, Required: false},
		"packer_force":               &hcldec.AttrSpec{Name: "packer_force", Type: cty.Bool, Required: false},
		"packer_on_error":            &hcldec.AttrSpec{Name: "packer_on_error", Type: cty.String, Required: false},
		"packer_user_variables":      &hcldec.AttrSpec{Name: "packer_user_variables", Type: cty.Map(cty.String), Required: false},
		"packer_sensitive_variables": &hcldec.AttrSpec{Name: "packer_sensitive_variables", Type: cty.List(cty.String), Required: false},
		"inline":                     &hcldec.AttrSpec{Name: "inline", Type: cty.List(cty.String), Required: false},
		"script":                     &hcldec.AttrSpec{Name: "script", Type: cty.String, Required: false},
		"scripts":                    &hcldec.AttrSpec{Name: "scripts", Type: cty.List(cty.String), Required: false},
		"valid_exit_codes":           &hcldec.AttrSpec{Name: "valid_exit_codes", Type: cty.List(cty.Number), Required: false},
		"environment_vars":           &hcldec.AttrSpec{Name: "environment_vars", Type: cty.List(cty.String), Required: false},
		"env_var_format":             &hcldec.AttrSpec{Name: "env_var_format", Type: cty.String, Required: false},
		"command":                    &hcldec.AttrSpec{Name: "command", Type: cty.String, Required: false},
		"execute_command":            &hcldec.AttrSpec{Name: "execute_command", Type: cty.List(cty.String), Required: false},
		"inline_shebang":             &hcldec.AttrSpec{Name: "inline_shebang", Type: cty.String, Required: false},
		"only_on":                    &hcldec.AttrSpec{Name: "only_on", Type: cty.List(cty.String), Required: false},
		"tempfile_extension":         &hcldec.AttrSpec{Name: "tempfile_extension", Type: cty.String, Required: false},
		"use_linux_pathing":          &hcldec.AttrSpec{Name: "use_linux_pathing", Type: cty.Bool, Required: false},
	}
	return s
}
