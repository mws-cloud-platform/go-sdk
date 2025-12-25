package conv

import (
	"strconv"

	"github.com/go-faster/jx"
	"golang.org/x/exp/constraints"
)

func EncodeStringInt[T constraints.Integer](e *jx.Encoder, v T) {
	var (
		buf  [32]byte
		n    int
		base = 10
	)
	// Write first quote
	buf[n] = '"'
	n++
	// Write integer
	n += len(strconv.AppendInt(buf[n:n], int64(v), base))
	// Write second quote
	buf[n] = '"'
	n++
	e.Raw(buf[:n])
}
