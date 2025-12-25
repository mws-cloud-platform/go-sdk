package page_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	mwspage "go.mws.cloud/go-sdk/internal/page"
	"go.mws.cloud/go-sdk/mws/page"
)

func TestPager(t *testing.T) {
	for _, v := range []struct {
		Name      string
		Client    *client
		Times     int
		Responses [][]any
		Error     error
	}{
		{
			Name:      "no items",
			Client:    &client{t: t},
			Times:     1,
			Responses: [][]any{nil},
		},
		{
			Name:      "no items drain",
			Client:    &client{t: t},
			Times:     2,
			Responses: [][]any{nil},
			Error:     mwspage.ErrNoItems,
		},
		{
			Name:   "no items error",
			Client: &client{t: t, err: errClient},
			Times:  2,
			Error:  errClient,
		},
		{
			Name:      "has items single",
			Client:    &client{t: t, items: [][]any{{1, "2"}, {3.0, 4}, {'5', 6}}},
			Times:     1,
			Responses: [][]any{{1, "2"}},
		},
		{
			Name:      "has items multiple",
			Client:    &client{t: t, items: [][]any{{1, "2"}, {3.0, 4}, {'5', 6}}},
			Times:     2,
			Responses: [][]any{{1, "2"}, {3.0, 4}},
		},
		{
			Name:      "has items drain",
			Client:    &client{t: t, items: [][]any{{1, "2"}, {3.0, 4}, {'5', 6}}},
			Times:     5,
			Responses: [][]any{{1, "2"}, {3.0, 4}, {'5', 6}, nil},
			Error:     mwspage.ErrNoItems,
		},
		{
			Name:   "has items error first",
			Client: &client{t: t, items: [][]any{{1, "2"}, {3.0, 4}, {'5', 6}}, err: errClient},
			Times:  3,
			Error:  errClient,
		},
		{
			Name:      "has items error",
			Client:    &client{t: t, items: [][]any{{1, "2"}, {3.0, 4}, {'5', 6}}, err: errClient, errFrom: 2},
			Times:     4,
			Responses: [][]any{{1, "2"}, {3.0, 4}},
			Error:     errClient,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			var (
				actual []any
				err    error
			)
			pager := page.NewPager(request{}, v.Client.do)
			for i := range v.Times {
				actual, err = pager.Next(t.Context())
				if err != nil {
					break
				}
				require.Equal(t, v.Responses[i], actual)
			}

			if v.Error != nil {
				require.ErrorIs(t, err, v.Error)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

const errClient = consterr.Error("client error")

type client struct {
	t       *testing.T
	cur     int
	err     error
	errFrom int
	items   [][]any
}

func (c *client) do(_ context.Context, req request) (response, error) {
	cur, _ := strconv.Atoi(ptr.Value(req.token))
	require.Equal(c.t, c.cur, cur)

	if c.err != nil && c.cur >= c.errFrom {
		return response{}, c.err
	}

	var (
		token *string
		items []any
	)
	if c.cur < len(c.items) {
		token = ptr.Get(strconv.Itoa(c.cur + 1))
		items = c.items[c.cur]
	}
	c.cur++

	return response{
		token: token,
		items: items,
	}, nil
}

type request struct {
	token *string
}

func (r request) WithPageToken(token *string) request {
	r.token = token
	return r
}

type response struct {
	token *string
	items []any
}

func (r response) GetNextPageToken() *string {
	return r.token
}

func (r response) GetItems() []any {
	return r.items
}
