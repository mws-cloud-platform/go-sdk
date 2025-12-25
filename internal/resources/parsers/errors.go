package parsers

import "go.mws.cloud/util-toolset/pkg/utils/consterr"

const (
	ErrReferenceParsing = consterr.Error("reference parsing")
	ErrPatternMatches   = consterr.Error("pattern doesn't match")
)
