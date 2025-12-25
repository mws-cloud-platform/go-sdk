package conv

import (
	"time"

	"github.com/go-faster/jx"
)

const (
	dateLayout = time.DateOnly
	timeLayout = time.TimeOnly
)

func Date(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func Time(t time.Time) time.Time {
	return time.Date(0, 0, 0, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

func DateTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

// EncodeDateTime encodes date-time to json.
func EncodeDateTime(s *jx.Encoder, v time.Time) {
	encodeDateTime(s, v)
}

// EncodeDateTime encodes date-time converted to UTC time zone to json.
func EncodeDateTimeUTC(s *jx.Encoder, v time.Time) {
	encodeDateTime(s, v.UTC())
}

func encodeDateTime(s *jx.Encoder, v time.Time) {
	const (
		roundTo  = 8
		length   = len(time.RFC3339)
		allocate = ((length + roundTo - 1) / roundTo) * roundTo
	)
	b := make([]byte, allocate)
	b = v.AppendFormat(b[:0], time.RFC3339)
	s.ByteStr(b)
}
