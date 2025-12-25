package ipaddress

import "go.mws.cloud/util-toolset/pkg/utils/consterr"

const (
	ErrInvalidIPAddressString = consterr.Error("invalid ip address string")
	ErrInvalidIPVersion       = consterr.Error("invalid ip version")
	ErrEmptyIPAddress         = consterr.Error("empty ip address")
)
