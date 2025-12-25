package attribute

type KeyValueType int8

const (
	KeyValueTypeInvalid KeyValueType = iota
	KeyValueTypeBool
	KeyValueTypeInt64
	KeyValueTypeFloat64
	KeyValueTypeString
	KeyValueTypeBoolSlice
	KeyValueTypeInt64Slice
	KeyValueTypeFloat64Slice
	KeyValueTypeStringSlice
)

type KeyValue struct {
	Key   string
	Value any
	Type  KeyValueType
}

func Bool(k string, v bool) KeyValue {
	return keyValue(k, v, KeyValueTypeBool)
}

func Int64(k string, v int64) KeyValue {
	return keyValue(k, v, KeyValueTypeInt64)
}

func Float64(k string, v float64) KeyValue {
	return keyValue(k, v, KeyValueTypeFloat64)
}

func String(k, v string) KeyValue {
	return keyValue(k, v, KeyValueTypeString)
}

func BoolSlice(k string, v []bool) KeyValue {
	return keyValue(k, v, KeyValueTypeBoolSlice)
}

func Int64Slice(k string, v []int64) KeyValue {
	return keyValue(k, v, KeyValueTypeInt64Slice)
}

func Float64Slice(k string, v []float64) KeyValue {
	return keyValue(k, v, KeyValueTypeFloat64Slice)
}

func StringSlice(k string, v []string) KeyValue {
	return keyValue(k, v, KeyValueTypeStringSlice)
}

func keyValue(key string, value any, typ KeyValueType) KeyValue {
	return KeyValue{
		Key:   key,
		Value: value,
		Type:  typ,
	}
}
