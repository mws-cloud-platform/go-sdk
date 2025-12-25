package page_test

import (
	"context"
	"iter"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/internal/page"
)

func TestPage(t *testing.T) {
	for _, v := range []struct {
		Name     string
		Pager    func() page.Pager[any]
		Expected []any
		Error    error
		Limit    int
	}{
		{
			Name: "no items",
			Pager: func() page.Pager[any] {
				return page.PagerFunc[any](func(context.Context) ([]any, error) { return nil, page.ErrNoItems })
			},
		},
		{
			Name:     "single no error",
			Pager:    ptr.Get(staticPager{total: 1}).clone,
			Expected: []any{24, 42},
		},
		{
			Name:     "multiple no error",
			Pager:    ptr.Get(staticPager{total: 3}).clone,
			Expected: []any{24, 42, 24, 42, 24, 42},
		},
		{
			Name:  "error",
			Pager: ptr.Get(staticPager{err: errPager}).clone,
			Error: errPager,
		},
		{
			Name:     "single error",
			Pager:    ptr.Get(staticPager{err: errPager, errFrom: 1, total: 3}).clone,
			Error:    errPager,
			Expected: []any{24, 42},
		},
		{
			Name:     "multiple error",
			Pager:    ptr.Get(staticPager{err: errPager, errFrom: 3, total: 5}).clone,
			Error:    errPager,
			Expected: []any{24, 42, 24, 42, 24, 42},
		},
		{
			Name:     "break on limit",
			Pager:    ptr.Get(staticPager{total: 3}).clone,
			Expected: []any{24, 42, 24, 42},
			Limit:    2,
		},
	} {
		t.Run(v.Name, func(t *testing.T) {
			t.Parallel()

			iterateActual, iterateErr := iterate(v.Pager(), v.Limit)
			allActual, allErr := all(v.Pager(), v.Limit*2)

			if v.Error != nil {
				require.ErrorIs(t, iterateErr, v.Error)
				require.ErrorIs(t, allErr, v.Error)
			} else {
				require.NoError(t, iterateErr)
				require.NoError(t, allErr)
			}

			require.Equal(t, v.Expected, iterateActual)
			require.Equal(t, v.Expected, allActual)
		})
	}
}

func iterate(pager page.Pager[any], limit int) (actual []any, err error) {
	return collect(page.Iterate(context.Background(), pager), limit, func(arr []any, items []any) []any {
		return append(arr, items...)
	})
}

func all(pager page.Pager[any], limit int) (actual []any, err error) {
	return collect(page.All(context.Background(), pager), limit, func(arr []any, item any) []any {
		return append(arr, item)
	})
}

func collect[K any](seq iter.Seq2[K, error], limit int, add func([]any, K) []any) (actual []any, err error) {
	var item K
	i := 0
	for item, err = range seq {
		if limit > 0 && i >= limit {
			break
		}
		if err != nil {
			return actual, err
		}
		actual = add(actual, item)
		i++
	}

	return actual, nil
}

const errPager = consterr.Error("pager error")

type staticPager struct {
	cur, total int
	err        error
	errFrom    int
}

func (p *staticPager) Next(context.Context) ([]any, error) {
	if p.err != nil && p.cur >= p.errFrom {
		return nil, p.err
	}
	if p.cur >= p.total {
		return nil, page.ErrNoItems
	}
	p.cur++

	return []any{24, 42}, nil
}

func (p *staticPager) clone() page.Pager[any] {
	return &staticPager{
		cur:     p.cur,
		total:   p.total,
		err:     p.err,
		errFrom: p.errFrom,
	}
}
