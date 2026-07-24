// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

func TestDecode(t *testing.T) {
	type Target struct {
		Name    string
		Address string
		Time    time.Duration
		Trilean Trilean
	}

	cases := map[string]struct {
		Input  []any
		Output *Target
		Opts   *DecodeOpts
	}{
		"basic": {
			[]any{
				map[string]any{
					"name":    "bar",
					"time":    "5s",
					"trilean": "true",
				},
			},
			&Target{
				Name:    "bar",
				Time:    5 * time.Second,
				Trilean: TriTrue,
			},
			nil,
		},

		"empty-string-trilean": {
			[]any{
				map[string]any{
					"trilean": "",
				},
			},
			&Target{
				Trilean: TriUnset,
			},
			nil,
		},

		"variables": {
			[]any{
				map[string]any{
					"name": "{{user `name`}}",
				},
				map[string]any{
					"packer_user_variables": map[string]string{
						"name": "bar",
					},
				},
			},
			&Target{
				Name: "bar",
			},
			nil,
		},

		"filter": {
			[]any{
				map[string]any{
					"name":    "{{user `name`}}",
					"address": "{{user `name`}}",
				},
				map[string]any{
					"packer_user_variables": map[string]string{
						"name": "bar",
					},
				},
			},
			&Target{
				Name:    "bar",
				Address: "{{user `name`}}",
			},
			&DecodeOpts{
				Interpolate: true,
				InterpolateFilter: &interpolate.RenderFilter{
					Include: []string{"name"},
				},
			},
		},

		"build name": {
			[]any{
				map[string]any{
					"name": "{{build_name}}",
				},
				map[string]any{
					"packer_build_name": "foo",
				},
			},
			&Target{
				Name: "foo",
			},
			nil,
		},

		"build type": {
			[]any{
				map[string]any{
					"name": "{{build_type}}",
				},
				map[string]any{
					"packer_builder_type": "foo",
				},
			},
			&Target{
				Name: "foo",
			},
			nil,
		},
	}

	for k, tc := range cases {
		var result Target
		err := Decode(&result, tc.Opts, tc.Input...)
		if err != nil {
			t.Fatalf("err: %s\n\n%s", k, err)
		}

		if !reflect.DeepEqual(&result, tc.Output) {
			t.Fatalf("bad:\n\n%#v\n\n%#v", &result, tc.Output)
		}
	}
}

func TestDecode_fixerRecommendations(t *testing.T) {
	type TestConfig struct {
		Name string
	}

	cases := []struct {
		Reason   string
		Input    []any
		Opts     *DecodeOpts
		Expected string
	}{
		{
			Reason: "If no plugin type is provided, don't try to match fixer options",
			Input: []any{
				map[string]any{
					"name":    "bar",
					"iso_md5": "13123412341234",
				},
			},
			Opts:     &DecodeOpts{},
			Expected: `unknown configuration key: '"iso_md5"'`,
		},
		{
			Reason: "iso_md5 should always recommend packer fix regardless of plugin type",
			Input: []any{
				map[string]any{
					"name":    "bar",
					"iso_md5": "13123412341234",
				},
			},
			Opts:     &DecodeOpts{PluginType: "someplugin"},
			Expected: `Deprecated configuration key: 'iso_md5'`,
		},
		{
			Reason: "filename option should generate a fixer recommendation for the manifest postprocessor",
			Input: []any{
				map[string]any{
					"name":     "bar",
					"filename": "fakefilename",
				},
			},
			Opts:     &DecodeOpts{PluginType: "packer.post-processor.manifest"},
			Expected: `Deprecated configuration key: 'filename'`,
		},
		{
			Reason: "filename option should generate an unknown key error for other plugins",
			Input: []any{
				map[string]any{
					"name":     "bar",
					"filename": "fakefilename",
				},
			},
			Opts:     &DecodeOpts{PluginType: "randomplugin"},
			Expected: `unknown configuration key: '"filename"'`,
		},
	}

	for _, tc := range cases {
		var result TestConfig
		err := Decode(&result, tc.Opts, tc.Input...)
		if err == nil {
			t.Fatalf("Should have had an error: %s", tc.Reason)
		}

		if !strings.Contains(err.Error(), tc.Expected) {
			t.Fatalf("Expected: %s\nActual: %s\n; Reason: %s", tc.Expected, err.Error(), tc.Reason)
		}
	}
}
