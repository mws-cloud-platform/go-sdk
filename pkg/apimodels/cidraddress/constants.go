package cidraddress

import "go.mws.cloud/util-toolset/pkg/utils/consterr"

const (
	ErrInvalidCIDRString     = consterr.Error("invalid cidr string")
	ErrInvalidCIDRVersion    = consterr.Error("invalid cidr version")
	ErrEmptyIPAddressOrIPNet = consterr.Error("empty ip address or ip network")
)
