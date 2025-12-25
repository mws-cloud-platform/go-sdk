package email

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/testing/golden"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"
)

const (
	fooEmail = "foo@bar"
)

func TestEmail_Clone(t *testing.T) {
	email, err := ParseString(fooEmail)
	require.NoError(t, err)

	clone := email.Clone()
	*email.rawValue = "rawValue"

	require.NotEqual(t, email.rawValue, clone.rawValue)
}

func TestEmail_Equal(t *testing.T) {
	email, err := ParseString(fooEmail)
	require.NoError(t, err)

	for _, testCase := range []struct {
		name   string
		email1 Email
		email2 Email
		equal  bool
	}{
		{
			name:  "empty",
			equal: true,
		},
		{
			name:   "different rawValue 1",
			email1: Email{},
			email2: Email{rawValue: ptr.Get("")},
			equal:  false,
		},
		{
			name:   "different rawValue 2",
			email1: Email{rawValue: ptr.Get("")},
			email2: Email{},
			equal:  false,
		},
		{
			name:   "different rawValue 3",
			email1: Email{rawValue: ptr.Get("hello")},
			email2: Email{rawValue: ptr.Get("world")},
			equal:  false,
		},
		{
			name:   "different username",
			email1: Email{username: "foo"},
			email2: Email{},
			equal:  false,
		},
		{
			name:   "different domain",
			email1: Email{},
			email2: Email{domain: "bar"},
			equal:  false,
		},
		{
			name:   "equal",
			email1: email,
			email2: email,
			equal:  true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.equal, testCase.email1.Equal(testCase.email2))
		})
	}
}

func TestEmail_String(t *testing.T) {
	email, err := ParseString(fooEmail)
	require.NoError(t, err)
	require.Equal(t, email.String(), fooEmail)

	email = Email{}
	require.Empty(t, email.String())
}

func TestEmail_MarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/marshal_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name  string
		email Email
	}{
		{
			name: "WithRawValue",
			email: Email{
				username: "foo",
				domain:   "bar",
				rawValue: ptr.Get("raw@value"),
			},
		},
		{
			name: "WithoutRawValue",
			email: Email{
				username: "foo",
				domain:   "bar",
			},
		},
		{
			name: "WithEmptyUsername",
			email: Email{
				username: "",
				domain:   "bar",
			},
		},
		{
			name: "WithEmptyDomain",
			email: Email{
				username: "foo",
				domain:   "",
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := json.Marshal(testCase.email)
			require.NoError(t, err)

			require.NoError(t, fs.WriteFile(testCase.name+".txt", result, 0644))
		})
	}
}

func TestEmail_UnmarshalJSON(t *testing.T) {
	dir := golden.NewDir(t, golden.WithPath("testdata/unmarshal_email_json.golden"), golden.WithRecreateOnUpdate())
	fs := golden.NewCodegenFS(t, dir)

	for _, testCase := range []struct {
		name        string
		rawValue    string
		errExpected bool
	}{
		{
			name:     "ok",
			rawValue: "foo@bar",
		},
		{
			name:     "borderlineCase",
			rawValue: "f@@",
		},
		{
			name:        "invalid 1",
			rawValue:    "f@",
			errExpected: true,
		},
		{
			name:        "invalid 2",
			rawValue:    "foo",
			errExpected: true,
		},
		{
			name:        "empty",
			rawValue:    "",
			errExpected: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			rawJSON := strconv.Quote(testCase.rawValue)

			var email Email
			err := json.Unmarshal([]byte(strconv.Quote(testCase.rawValue)), &email)
			if testCase.errExpected {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NoError(t, fs.WriteFile(testCase.name+".txt",
				[]byte("RawJSON: "+rawJSON+"\nParsedUsername: "+email.Username()+"\nParsedDomain: "+email.Domain()), 0644))
		})
	}
}
