package merge_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	commonclient "go.mws.cloud/go-sdk/internal/client"
	"go.mws.cloud/go-sdk/internal/merge"
)

func TestSlice(t *testing.T) {
	for _, v := range []struct {
		Name     string
		Dst      []model
		Src      []updateModel
		Expected []model
	}{
		{
			Name: "empty_all",
		},
		{
			Name: "empty_src",
			Dst:  []model{{name: "foo"}},
		},
		{
			Name:     "empty_dst",
			Src:      []updateModel{{name: commonclient.NewOptional("foo")}},
			Expected: []model{{name: "foo"}},
		},
		{
			Name:     "no_common",
			Dst:      []model{{name: "qux"}, {name: "www"}},
			Src:      []updateModel{{name: commonclient.NewOptional("foo")}, {name: commonclient.NewOptional("bar")}},
			Expected: []model{{name: "foo"}, {name: "bar"}},
		},
		{
			Name: "has_common",
			Dst: []model{
				{name: "qux"}, {name: "www"},
				{name: "foo", number: 42},
			},
			Src: []updateModel{
				{name: commonclient.NewOptional("bar")},
				{name: commonclient.NewOptional("foo"), number: commonclient.NewOptional(24)},
			},
			Expected: []model{
				{name: "bar"},
				{name: "foo", number: 24},
			},
		},
		{
			Name: "all_common",
			Dst: []model{
				{name: "foo", number: 42},
				{name: "bar"},
				{name: "baz", number: -1},
			},
			Src: []updateModel{
				{name: commonclient.NewOptional("bar"), number: commonclient.NewOptional(24)},
				{name: commonclient.NewOptional("baz"), number: commonclient.NewOptional(1)},
				{name: commonclient.NewOptional("foo")},
			},
			Expected: []model{
				{name: "bar", number: 24},
				{name: "baz", number: 1},
				{name: "foo", number: 42},
			},
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			actual := merge.Slice(v.Dst, v.Src, (*model).withChanges, (*model).getName, (*updateModel).getName)
			require.Equal(t, v.Expected, actual)
		})
	}
}

func TestMap(t *testing.T) {
	for _, v := range []struct {
		Name     string
		Dst      map[string]model
		Src      map[string]updateModel
		Expected map[string]model
	}{
		{
			Name: "empty_all",
		},
		{
			Name: "empty_src",
			Dst:  map[string]model{"foo": {}},
		},
		{
			Name:     "empty_dst",
			Src:      map[string]updateModel{"foo": {}},
			Expected: map[string]model{"foo": {}},
		},
		{
			Name:     "no_common",
			Dst:      map[string]model{"qux": {}, "www": {}},
			Src:      map[string]updateModel{"foo": {}, "bar": {}},
			Expected: map[string]model{"foo": {}, "bar": {}},
		},
		{
			Name: "has_common",
			Dst: map[string]model{
				"qux": {},
				"www": {},
				"foo": {name: "hello", number: 42},
			},
			Src: map[string]updateModel{
				"foo": {name: commonclient.NewOptional("world")},
				"bar": {},
			},
			Expected: map[string]model{
				"foo": {name: "world", number: 42},
				"bar": {},
			},
		},
		{
			Name: "all_common",
			Dst: map[string]model{
				"foo": {name: "hello", number: 42},
				"bar": {name: "bar"},
			},
			Src: map[string]updateModel{
				"foo": {name: commonclient.NewOptional("world")},
				"bar": {number: commonclient.NewOptional(24)},
			},
			Expected: map[string]model{
				"foo": {name: "world", number: 42},
				"bar": {name: "bar", number: 24},
			},
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			actual := merge.Map(v.Dst, v.Src, (*model).withChanges)
			require.Equal(t, v.Expected, actual)
		})
	}
}

func TestInapplicableSlice(t *testing.T) {
	for _, v := range []struct {
		Name     string
		Src      []updateModel
		Expected []model
	}{
		{
			Name: "empty",
		},
		{
			Name:     "single_empty",
			Src:      []updateModel{{}},
			Expected: []model{{}},
		},
		{
			Name:     "multiple_empty",
			Src:      []updateModel{{}, {}},
			Expected: []model{{}, {}},
		},
		{
			Name:     "single_one_field",
			Src:      []updateModel{{number: commonclient.NewOptional(20)}},
			Expected: []model{{number: 20}},
		},
		{
			Name:     "single_multiple_fields",
			Src:      []updateModel{{name: commonclient.NewOptional("foo"), number: commonclient.NewOptional(20)}},
			Expected: []model{{name: "foo", number: 20}},
		},
		{
			Name: "multiple",
			Src: []updateModel{
				{name: commonclient.NewOptional("bar")},
				{name: commonclient.NewOptional("foo"), number: commonclient.NewOptional(20)},
				{number: commonclient.NewOptional(42)},
			},
			Expected: []model{
				{name: "bar"},
				{name: "foo", number: 20},
				{number: 42},
			},
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			actual := merge.InapplicableSlice(v.Src, (*model).withChanges)
			require.Equal(t, v.Expected, actual)
		})
	}
}

type model struct {
	name   string
	number int
}

func (m *model) getName() string {
	return m.name
}

func (m *model) withChanges(u updateModel) model {
	var out model
	if m != nil {
		out = *m
	}

	if u.name.IsSet() {
		out.name = u.name.Value
	}
	if u.number.IsSet() {
		out.number = u.number.Value
	}
	return out
}

type updateModel struct {
	name   commonclient.Optional[string]
	number commonclient.Optional[int]
}

func (u *updateModel) getName() string {
	if u.name.IsSet() {
		return u.name.Value
	}
	return ""
}
