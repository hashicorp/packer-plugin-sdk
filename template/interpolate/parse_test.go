// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package interpolate

import (
	"reflect"
	"testing"
	"text/template"
)

func TestFunctionsCalled(t *testing.T) {
	cases := []struct {
		Input  string
		Result map[string]struct{}
	}{
		{
			"foo",
			map[string]struct{}{},
		},

		{
			"foo {{user `bar`}}",
			map[string]struct{}{
				"user": {},
			},
		},
	}

	funcs := Funcs(&Context{})
	for _, tc := range cases {
		tpl, err := template.New("root").Funcs(funcs).Parse(tc.Input)
		if err != nil {
			t.Fatalf("err parsing: %v\n\n%s", tc.Input, err)
		}

		actual := functionsCalled(tpl)
		if !reflect.DeepEqual(actual, tc.Result) {
			t.Fatalf("bad: %v\n\ngot: %#v", tc.Input, actual)
		}
	}
}
