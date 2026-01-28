package zap

import (
	"go.uber.org/zap"

	"go.mws.cloud/go-sdk/internal/client/interceptors/attribute"
)

func Convert(attrs []attribute.KeyValue) []zap.Field {
	fields := make([]zap.Field, 0, len(attrs))
	for _, attr := range attrs {
		var field zap.Field
		switch attr.Type {
		case attribute.KeyValueTypeInvalid:
			continue
		case attribute.KeyValueTypeBool:
			field = zap.Bool(attr.Key, attr.Value.(bool))
		case attribute.KeyValueTypeInt64:
			field = zap.Int64(attr.Key, attr.Value.(int64))
		case attribute.KeyValueTypeFloat64:
			field = zap.Float64(attr.Key, attr.Value.(float64))
		case attribute.KeyValueTypeString:
			field = zap.String(attr.Key, attr.Value.(string))
		case attribute.KeyValueTypeBoolSlice:
			field = zap.Bools(attr.Key, attr.Value.([]bool))
		case attribute.KeyValueTypeInt64Slice:
			field = zap.Int64s(attr.Key, attr.Value.([]int64))
		case attribute.KeyValueTypeFloat64Slice:
			field = zap.Float64s(attr.Key, attr.Value.([]float64))
		case attribute.KeyValueTypeStringSlice:
			field = zap.Strings(attr.Key, attr.Value.([]string))
		}

		fields = append(fields, field)
	}

	return fields
}
