package conv

import (
	"encoding/base64"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func StringToInt[T constraints.Integer](s string, bitSize int) (T, error) {
	i, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}

func StringToIntFn[T constraints.Integer](bitSize int) func(string) (T, error) {
	return func(s string) (T, error) {
		return StringToInt[T](s, bitSize)
	}
}

func StringToUint[T constraints.Unsigned](s string, bitSize int) (T, error) {
	i, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}

func StringToFloat[T constraints.Float](s string, bitSize int) (T, error) {
	i, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}

func StringToBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func StringToStringArray(s string) []string {
	return strings.Split(s, ",")
}

func StringToStringArrayErr(s string) ([]string, error) {
	return StringToStringArray(s), nil
}

func StringToIntArray[T constraints.Integer](s string, bitSize int) ([]T, error) {
	return stringToArray(s, func(n string) (T, error) {
		return StringToInt[T](n, bitSize)
	})
}

func StringToUintArray[T constraints.Unsigned](s string, bitSize int) ([]T, error) {
	return stringToArray(s, func(n string) (T, error) {
		return StringToUint[T](n, bitSize)
	})
}

func StringToFloatArray[T constraints.Float](s string, bitSize int) ([]T, error) {
	return stringToArray(s, func(n string) (T, error) {
		return StringToFloat[T](n, bitSize)
	})
}

func StringToBoolArray(s string) ([]bool, error) {
	return stringToArray(s, StringToBool)
}

func StringToStringSubsetArray[T ~string](s string) []T {
	split := StringToStringArray(s)
	out := make([]T, len(split))
	for i, n := range split {
		out[i] = StringToStringSubset[T](n)
	}
	return out
}

func StringToStringSubsetArrayErr[T ~string](s string) ([]T, error) {
	return StringToStringSubsetArray[T](s), nil
}

func DecodeBase64(r io.Reader) ([]byte, error) {
	return io.ReadAll(base64.NewDecoder(base64.StdEncoding, r))
}

func stringToArray[T any](s string, converter func(string) (T, error)) ([]T, error) {
	var err error
	split := StringToStringArray(s)
	result := make([]T, len(split))

	for i, n := range split {
		result[i], err = converter(n)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func StringToDateTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func StringToBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
