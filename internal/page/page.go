package page

import (
	"context"
	"errors"
	"iter"

	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const ErrNoItems = consterr.Error("no items left")

type Pager[T any] interface {
	Next(context.Context) ([]T, error)
}

type PagerWithHasNext[T any] interface {
	Next(context.Context) ([]T, error)
	HasNext() bool
}

type PagerFunc[T any] func(context.Context) ([]T, error)

func (f PagerFunc[T]) Next(ctx context.Context) ([]T, error) {
	return f(ctx)
}

func Iterate[T any](ctx context.Context, pager Pager[T]) iter.Seq2[[]T, error] {
	return func(yield func([]T, error) bool) {
		for {
			items, err := pager.Next(ctx)
			if errors.Is(err, ErrNoItems) || !yield(items, err) {
				return
			}
		}
	}
}

func All[T any](ctx context.Context, pager Pager[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for items, err := range Iterate(ctx, pager) {
			if err != nil {
				var empty T
				yield(empty, err)
				return
			}
			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}
		}
	}
}
