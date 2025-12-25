// Package page provides utilities for working with paginated API responses.
package page

import (
	"context"
	"iter"

	"go.mws.cloud/util-toolset/pkg/utils/ptr"

	"go.mws.cloud/go-sdk/internal/page"
)

// Request represents a request for methods that return paginated response.
type Request[T any] interface {
	WithPageToken(*string) T
}

// Response represents a paginated response.
type Response[Token ~string, T any] interface {
	GetNextPageToken() *Token
	GetItems() []T
}

// Do is a function that performs an API request with a paginated response.
type Do[Token ~string, T any, Req Request[Req], Resp Response[Token, T]] func(context.Context, Req) (Resp, error)

// NewPager creates a new pager for the given request and do function.
func NewPager[Token ~string, T any, Req Request[Req], Resp Response[Token, T]](
	req Req,
	do Do[Token, T, Req, Resp],
) *Pager[Token, T, Req, Resp] {
	return &Pager[Token, T, Req, Resp]{
		req: req,
		do:  do,
	}
}

// Pager is an iterator over paginated API method.
type Pager[Token ~string, T any, Req Request[Req], Resp Response[Token, T]] struct {
	req  Req
	do   Do[Token, T, Req, Resp]
	stop bool
}

// Next returns items from the next page.
func (p *Pager[Token, T, Req, Resp]) Next(ctx context.Context) (data []T, err error) {
	if p.stop {
		return nil, page.ErrNoItems
	}

	resp, err := p.do(ctx, p.req)
	if err != nil {
		return nil, err
	}

	var token *string
	if nextToken := resp.GetNextPageToken(); nextToken == nil || *nextToken == "" {
		p.stop = true
	} else {
		token = ptr.Get(string(*nextToken))
	}

	p.req = p.req.WithPageToken(token)
	return resp.GetItems(), nil
}

// HasNext returns true if there are more pages to iterate.
func (p *Pager[Token, T, Req, Resp]) HasNext() bool {
	return !p.stop
}

// Pages returns an iterator over the pages of the response.
func (p *Pager[Token, T, Req, Resp]) Pages(ctx context.Context) iter.Seq2[[]T, error] {
	return page.Iterate(ctx, p)
}

// All returns an iterator over all items in the response.
func (p *Pager[Token, T, Req, Resp]) All(ctx context.Context) iter.Seq2[T, error] {
	return page.All(ctx, p)
}
