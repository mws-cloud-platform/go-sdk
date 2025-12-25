package macaddress

import "go.mws.cloud/util-toolset/pkg/utils/consterr"

const (
	ErrInvalidMACString = consterr.Error("invalid MAC string")
	ErrInvalidMACFormat = consterr.Error("invalid MAC format")
	ErrEmptyMACAddress  = consterr.Error("empty MAC address")

	eui48octetsCount = 6
	eui64octetsCount = 8
)
