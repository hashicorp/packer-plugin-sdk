package internal

import (
	"reflect"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestMapOfInterfaceToMapOfCTY(t *testing.T) {
	input := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
		},
		"bar": []interface{}{
			"foo",
		},
		"bFalse": false,
		"bTrue":  true,
		"bNil":   nil,
		"bStr":   "bar",
		"bInt":   1,
		"bFloat": 4.5,
	}

	expectedOutput := map[string]cty.Value{
		"foo": cty.MapVal(map[string]cty.Value{
			"bar": cty.StringVal("baz"),
		}),
		"bar": cty.TupleVal([]cty.Value{
			cty.StringVal("foo"),
		}),
		"bFalse": cty.False,
		"bTrue":  cty.True,
		"bStr":   cty.StringVal("bar"),
		"bInt":   cty.NumberIntVal(1),
		"bFloat": cty.NumberFloatVal(4.5),
	}

	output, err := MapOfInterfaceToMapOfCTY(reflect.TypeOf(input), reflect.TypeOf(expectedOutput), input)
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if _, ok := output.(map[string]cty.Value); !ok {
		t.Fatalf("expecting output to be of type map[string]cty.Value")
	}

	outCty := output.(map[string]cty.Value)
	for k, v := range expectedOutput {
		outVal, ok := outCty[k]
		if !ok {
			t.Fatalf("expected key %s not found", k)
		}
		equals := v.Equals(outVal)
		if equals.False() {
			t.Fatalf("Unexpected val for key: \nkey: %s \nval: %s\nexpecting: %s", k, outVal.GoString(), v.GoString())
		}
	}
}

func TestInterfaceToCTY(t *testing.T) {
	tests := []struct {
		Name  string
		Input interface{}
		Want  cty.Value
	}{
		{
			Name:  "bool true",
			Input: true,
			Want:  cty.True,
		},
		{
			Name:  "bool false",
			Input: false,
			Want:  cty.False,
		},
		{
			Name:  "int",
			Input: int(12),
			Want:  cty.NumberIntVal(12),
		},
		{
			Name:  "float64",
			Input: float64(12.5),
			Want:  cty.NumberFloatVal(12.5),
		},
		{
			Name:  "string",
			Input: "hello world",
			Want:  cty.StringVal("hello world"),
		},
		{
			Name: "nested map[string]interface{}",
			Input: map[string]interface{}{
				"name": "Ermintrude",
				"age":  int(19),
				"address": map[string]interface{}{
					"street": []interface{}{"421 Shoreham Loop"},
					"city":   "Fridgewater",
					"state":  "MA",
					"zip":    "91037",
				},
			},
			Want: cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Ermintrude"),
				"age":  cty.NumberIntVal(19),
				"address": cty.ObjectVal(map[string]cty.Value{
					"street": cty.TupleVal([]cty.Value{cty.StringVal("421 Shoreham Loop")}),
					"city":   cty.StringVal("Fridgewater"),
					"state":  cty.StringVal("MA"),
					"zip":    cty.StringVal("91037"),
				}),
			}),
		},
		{
			Name: "simple map[string]interface{}",
			Input: map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
			Want: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
				"bar": cty.StringVal("baz"),
			}),
		},
		{
			Name: "[]interface{} as tuple",
			Input: []interface{}{
				"foo",
				true,
			},
			Want: cty.TupleVal([]cty.Value{
				cty.StringVal("foo"),
				cty.True,
			}),
		},
		{
			Name:  "nil",
			Input: nil,
			Want:  cty.NullVal(cty.DynamicPseudoType),
		},
		{
			Name: "SliceString",
			Input: []string{
				"a",
				"b",
				"c",
			},
			Want: cty.ListVal([]cty.Value{
				cty.StringVal("a"),
				cty.StringVal("b"),
				cty.StringVal("c"),
			}),
		},
		{
			Name:  "Empty SliceString",
			Input: []string{},
			Want:  cty.ListValEmpty(cty.String),
		},
		{
			Name: "map[interface{}]interface{}",
			Input: map[interface{}]interface{}{
				"name": "Ermintrude",
				1:      2,
			},
			Want: cty.TupleVal([]cty.Value{
				cty.StringVal("name"),
				cty.StringVal("Ermintrude"),
				cty.NumberIntVal(1),
				cty.NumberIntVal(2),
			}),
		},
		{
			Name:  "empty map[interface{}]interface{}",
			Input: map[interface{}]interface{}{},
			Want:  cty.EmptyTupleVal,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			got := InterfaceToCTY(test.Input)
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("wrong result\ninput: %#v\ngot:   %#v\nwant:  %#v", test.Input, got, test.Want)
			}
		})
	}
}
