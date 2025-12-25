package attribute

import (
	"context"
	"fmt"
)

type keyAttribute struct{}

func FromContext(ctx context.Context) []KeyValue {
	if ctx == nil {
		return nil
	}

	val := ctx.Value(keyAttribute{})
	if val == nil {
		return nil
	}

	attributes, ok := val.([]KeyValue)
	if !ok {
		panic(fmt.Sprintf("Unexpected type '%T' for attributes holder", val))
	}

	return attributes
}

func WithContext(ctx context.Context, attributes ...KeyValue) context.Context {
	return context.WithValue(ctx, keyAttribute{}, attributes)
}
