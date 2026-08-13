package optional_test

import (
	"encoding/json"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/go-faster/jx"
	"github.com/stretchr/testify/require"

	"go.mws.cloud/go-sdk/internal/conv"
	"go.mws.cloud/go-sdk/pkg/optional"
)

//nolint:cyclop // this is intended for deep compare
func TestOptionalUnmarshalJSON(t *testing.T) {
	for _, v := range []struct {
		Name    string
		Raw     string
		Empty   func() any
		Compare func(any) bool
	}{
		{
			Name: "int",
			Raw:  "5",
			Empty: func() any {
				return new(optional.NewOptional(0))
			},
			Compare: func(a any) bool {
				return a.(*optional.Optional[int]).Value == 5
			},
		},
		{
			Name: "string",
			Raw:  `"hello"`,
			Empty: func() any {
				return new(optional.NewOptional(""))
			},
			Compare: func(a any) bool {
				return a.(*optional.Optional[string]).Value == "hello"
			},
		},
		{
			Name: "bool",
			Raw:  `true`,
			Empty: func() any {
				return new(optional.NewOptional(false))
			},
			Compare: func(a any) bool {
				return a.(*optional.Optional[bool]).Value
			},
		},
		{
			Name: "object",
			Raw:  `{ "int": 42 }`,
			Empty: func() any {
				return new(optional.NewOptional(plainObject{}))
			},
			Compare: func(a any) bool {
				return a.(*optional.Optional[plainObject]).Value == plainObject{Int: 42}
			},
		},
		{
			Name: "objectWithOptionals",
			Raw: `{ "int": 42, "slice": [ 1, 2, 3 ], "optional": 24, "map": { "key": 25 }, "int_nil": null, "string_nil": "xyz", 
					"slice_nil": null, "map_nil": {}, "optional_int_nil": 123, "optional_string_nil": null }`,
			Empty: func() any {
				return &objectWithOptionals{}
			},
			Compare: func(a any) bool {
				o := *a.(*objectWithOptionals)
				return o.Int.Set && o.Int.Value == 42 &&
					o.Slice.Set && slices.Equal(o.Slice.Value, []int{1, 2, 3}) &&
					!o.Object.Set &&
					o.Optional.Set && o.Optional.Value.Set && o.Optional.Value.Value == 24 &&
					o.Map.Set && maps.Equal(o.Map.Value, map[string]int{"key": 25}) &&
					o.IntNil.Set && o.IntNil.Null &&
					o.StringNil.Set && !o.StringNil.Null && o.StringNil.Value == "xyz" &&
					o.SliceNil.Set && o.SliceNil.Null &&
					o.MapNil.Set && !o.MapNil.Null && maps.Equal(o.MapNil.Value, map[string]int{}) &&
					o.OptionalIntNil.Set && !o.OptionalIntNil.Null && o.OptionalIntNil.Value.Set && o.OptionalIntNil.Value.Value == 123 &&
					o.OptionalStringNil.Set && o.OptionalStringNil.Null
			},
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			e := v.Empty()
			require.NoError(t, json.Unmarshal([]byte(v.Raw), e))
			require.True(t, v.Compare(e))

			e = v.Empty()
			u, ok := e.(unmarshaler)
			if ok {
				require.NoError(t, u.Decode(jx.DecodeStr(v.Raw)))
				require.True(t, v.Compare(e))
			}
		})
	}
}

//nolint:cyclop // this is intended for deep compare
func TestOptionalUnmarshalYAML(t *testing.T) {
	for _, v := range []struct {
		Name    string
		Raw     string
		Empty   func() any
		Compare func(any) bool
	}{
		{
			Name: "int",
			Raw:  "5",
			Empty: func() any {
				return new(optional.NewOptional(0))
			},
			Compare: func(a any) bool {
				return a.(*optional.Optional[int]).Value == 5
			},
		},
		{
			Name: "string",
			Raw:  `"hello"`,
			Empty: func() any {
				return new(optional.NewOptional(""))
			},
			Compare: func(a any) bool {
				return a.(*optional.Optional[string]).Value == "hello"
			},
		},
		{
			Name: "bool",
			Raw:  `true`,
			Empty: func() any {
				return new(optional.NewOptional(false))
			},
			Compare: func(a any) bool {
				return a.(*optional.Optional[bool]).Value
			},
		},
		{
			Name: "object",
			Raw:  `int: 42`,
			Empty: func() any {
				return new(optional.NewOptional(plainObject{}))
			},
			Compare: func(a any) bool {
				return a.(*optional.Optional[plainObject]).Value == plainObject{Int: 42}
			},
		},
		{
			Name: "objectWithOptionals",
			Raw: `int: 42
                  slice: [1, 2, 3]
                  optional: 24
                  map:
                    key: 25
                  int_nil: null
                  string_nil: xyz
                  slice_nil: null
                  map_nil: {}
                  optional_int_nil: 123
                  optional_string_nil: null`,
			Empty: func() any {
				return &objectWithOptionals{}
			},
			Compare: func(a any) bool {
				o := *a.(*objectWithOptionals)
				return o.Int.Set && o.Int.Value == 42 &&
					o.Slice.Set && slices.Equal(o.Slice.Value, []int{1, 2, 3}) &&
					!o.Object.Set &&
					o.Optional.Set && o.Optional.Value.Set && o.Optional.Value.Value == 24 &&
					o.Map.Set && maps.Equal(o.Map.Value, map[string]int{"key": 25}) &&
					o.IntNil.Set && o.IntNil.Null &&
					o.StringNil.Set && !o.StringNil.Null && o.StringNil.Value == "xyz" &&
					o.SliceNil.Set && o.SliceNil.Null &&
					o.MapNil.Set && !o.MapNil.Null && maps.Equal(o.MapNil.Value, map[string]int{}) &&
					o.OptionalIntNil.Set && !o.OptionalIntNil.Null && o.OptionalIntNil.Value.Set && o.OptionalIntNil.Value.Value == 123 &&
					o.OptionalStringNil.Set && o.OptionalStringNil.Null
			},
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			split := strings.Split(v.Raw, "\n")
			spaces := 0
			for i, s := range split {
				if spaces == 0 && s[0] == ' ' {
					spaces = len(s) - len(strings.TrimSpace(s))
				}
				split[i] = s[spaces:]
			}
			raw := strings.Join(split, "\n")
			data, err := conv.YAMLtoJSON([]byte(raw))
			require.NoError(t, err)

			e := v.Empty()
			require.NoError(t, json.Unmarshal(data, e))
			require.True(t, v.Compare(e))
		})
	}
}

type unmarshaler interface {
	Decode(d *jx.Decoder) error
}

type objectWithOptionals struct {
	Int               optional.Optional[int]                          `json:"int" yaml:"int"`
	Slice             optional.Optional[[]int]                        `json:"slice" yaml:"slice"`
	Object            optional.Optional[plainObject]                  `json:"object" yaml:"object"`
	Optional          optional.Optional[optional.Optional[int]]       `json:"optional" yaml:"optional"`
	Map               optional.Optional[map[string]int]               `json:"map" yaml:"map"`
	IntNil            optional.OptionalNil[int]                       `json:"int_nil" yaml:"int_nil"`
	StringNil         optional.OptionalNil[string]                    `json:"string_nil" yaml:"string_nil"`
	SliceNil          optional.OptionalNil[[]int]                     `json:"slice_nil" yaml:"slice_nil"`
	MapNil            optional.OptionalNil[map[string]int]            `json:"map_nil" yaml:"map_nil"`
	OptionalIntNil    optional.OptionalNil[optional.Optional[int]]    `json:"optional_int_nil" yaml:"optional_int_nil"`
	OptionalStringNil optional.OptionalNil[optional.Optional[string]] `json:"optional_string_nil" yaml:"optional_string_nil"`
}

type plainObject struct {
	Int    int    `json:"int" yaml:"int"`
	String string `json:"string" yaml:"string"`
}
