package rangeunit

import "go.mws.cloud/util-toolset/pkg/utils/consterr"

const (
	ErrDifferentUnitValues         = consterr.Error("unit values are not equal")
	ErrMinValueGreaterThanMaxValue = consterr.Error("min value is greater than max value")
	ErrIncorrectMinValue           = consterr.Error("incorrect min value")
	ErrIncorrectMaxValue           = consterr.Error("incorrect max value")
)
