package conv_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"

	"go.mws.cloud/go-sdk/internal/conv"
	"go.mws.cloud/go-sdk/pkg/apimodels/cidraddress"
	"go.mws.cloud/go-sdk/pkg/resources/models"
)

func TestYAML(t *testing.T) {
	expected := object{
		Ref: models.NewAnyResourceRef("hello/world"),
		Test: test{
			Obj: obj{
				KK: "testdata_kk",
				HH: "testdata_hh",
				AA: "testdata_aa",
			},
			B: 121,
			G: "testdata_g",
		},
		ID: models.NewAnyResourceID("foo/bar"),
		IP: cidraddress.MustParseCIDR4AddressString("192.0.2.0/24"),
	}

	data, err := conv.JSONtoYAML(expected)
	require.NoError(t, err)

	dir := golden.NewDir(t, golden.WithPath("testdata/"+"json_to_yaml"+"/golden"),
		golden.WithRecreateOnUpdate())
	dir.Bytes(t, "expected.yaml", data)

	data, err = conv.YAMLtoJSON(data)
	require.NoError(t, err)

	dir = golden.NewDir(t, golden.WithPath("testdata/"+"yaml_to_json"+"/golden"),
		golden.WithRecreateOnUpdate())
	dir.Bytes(t, "expected.json", data)

	actual := object{}
	require.NoError(t, json.Unmarshal(data, &actual))

	require.Equal(t, expected, actual)
}

func TestStringToYAML(t *testing.T) {
	data, err := conv.JSONtoYAML("string")
	require.NoError(t, err)
	require.Equal(t, "string\n", string(data))
}

type object struct {
	Ref  models.AnyResourceRef    `json:"ref" yaml:"ref"`
	Test test                     `json:"test" yaml:"test"`
	ID   models.AnyResourceID     `json:"id" yaml:"id"`
	IP   cidraddress.CIDR4Address `json:"ip" yaml:"ip"`
}

type test struct {
	Obj obj    `json:"object" yaml:"object"`
	B   int    `json:"b" yaml:"b"`
	G   string `json:"g" yaml:"g"`
}

type obj struct {
	KK string `json:"kk" yaml:"kk"`
	AA string `json:"aa" yaml:"аа"`
	HH string `json:"hh" yaml:"hh"`
}
