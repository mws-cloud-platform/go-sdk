package otel

import (
	otelattribute "go.opentelemetry.io/otel/attribute"

	"go.mws.cloud/go-sdk/internal/client/interceptors/attribute"
)

func Convert(attrs []attribute.KeyValue) []otelattribute.KeyValue {
	out := make([]otelattribute.KeyValue, 0, len(attrs))
	for _, attr := range attrs {
		var kv otelattribute.KeyValue
		switch attr.Type {
		case attribute.KeyValueTypeInvalid:
			continue
		case attribute.KeyValueTypeBool:
			kv = otelattribute.Bool(attr.Key, attr.Value.(bool))
		case attribute.KeyValueTypeInt64:
			kv = otelattribute.Int64(attr.Key, attr.Value.(int64))
		case attribute.KeyValueTypeFloat64:
			kv = otelattribute.Float64(attr.Key, attr.Value.(float64))
		case attribute.KeyValueTypeString:
			kv = otelattribute.String(attr.Key, attr.Value.(string))
		case attribute.KeyValueTypeBoolSlice:
			kv = otelattribute.BoolSlice(attr.Key, attr.Value.([]bool))
		case attribute.KeyValueTypeInt64Slice:
			kv = otelattribute.Int64Slice(attr.Key, attr.Value.([]int64))
		case attribute.KeyValueTypeFloat64Slice:
			kv = otelattribute.Float64Slice(attr.Key, attr.Value.([]float64))
		case attribute.KeyValueTypeStringSlice:
			kv = otelattribute.StringSlice(attr.Key, attr.Value.([]string))
		}

		out = append(out, kv)
	}

	return out
}
