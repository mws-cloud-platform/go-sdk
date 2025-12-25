package decode

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-faster/jx"
)

func Str(d *jx.Decoder) (string, error) {
	return primitive(d, d.Str)
}

func StrBytes(d *jx.Decoder) ([]byte, error) { return primitive(d, d.StrBytes) }
func Bool(d *jx.Decoder) (bool, error)       { return primitive(d, d.Bool) }
func UInt(d *jx.Decoder) (uint, error)       { return primitive(d, d.UInt) }
func Int(d *jx.Decoder) (int, error)         { return primitive(d, d.Int) }
func UInt8(d *jx.Decoder) (uint8, error)     { return primitive(d, d.UInt8) }
func Int8(d *jx.Decoder) (int8, error)       { return primitive(d, d.Int8) }
func UInt16(d *jx.Decoder) (uint16, error)   { return primitive(d, d.UInt16) }
func Int16(d *jx.Decoder) (int16, error)     { return primitive(d, d.Int16) }
func UInt32(d *jx.Decoder) (uint32, error)   { return primitive(d, d.UInt32) }
func Int32(d *jx.Decoder) (int32, error)     { return primitive(d, d.Int32) }
func UInt64(d *jx.Decoder) (uint64, error)   { return primitive(d, d.UInt64) }
func Int64(d *jx.Decoder) (int64, error)     { return primitive(d, d.Int64) }
func Float32(d *jx.Decoder) (float32, error) { return primitive(d, d.Float32) }
func Float64(d *jx.Decoder) (float64, error) { return primitive(d, d.Float64) }

func DateTime(d *jx.Decoder) (time.Time, error) {
	v, err := d.Str()
	if err != nil {
		return time.Time{}, errorWithRawValue(err, d)
	}
	return time.Parse(time.RFC3339, v)
}

func StringInt(d *jx.Decoder) (int64, error) {
	s, err := d.StrBytes()
	if err != nil {
		return 0, errorWithRawValue(err, d)
	}
	v, err := jx.DecodeBytes(s).Int64()
	if err != nil {
		return 0, errorWithRawValue(err, d)
	}
	return v, nil
}

func errorWithRawValue(err error, d *jx.Decoder) error {
	r, rawErr := d.Raw()
	err = errors.Join(err, rawErr)
	rawString := r.String()
	if rawString == "" {
		return err
	}
	return fmt.Errorf("%w: raw value '%s'", err, rawString)
}

func primitive[T any](d *jx.Decoder, decode func() (T, error)) (v T, err error) {
	v, err = decode()
	if err != nil {
		return v, errorWithRawValue(err, d)
	}
	return v, nil
}
