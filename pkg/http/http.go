// Package http provides batteries included HTTP client for the SDK service
// clients.
package http

import (
	"net"
	"net/http"
	"time"
)

const (
	defaultTimeout               = 5 * time.Second
	defaultIdleConnTimeout       = 90 * time.Second
	defaultTLSHandshakeTimeout   = 10 * time.Second
	defaultExpectContinueTimeout = time.Second
	defaultTCPTimeout            = 15 * time.Second
	defaultKeepAlive             = 15 * time.Second
)

// RoundTripperFunc is an adapter to allow the use of ordinary functions as HTTP
// transport.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewClient creates a new HTTP client.
func NewClient() *http.Client {
	return &http.Client{
		Timeout:   defaultTimeout,
		Transport: DefaultTransport(),
	}
}

// DefaultTransport creates a new default HTTP client transport.
func DefaultTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	t.IdleConnTimeout = defaultIdleConnTimeout
	t.TLSHandshakeTimeout = defaultTLSHandshakeTimeout
	t.ExpectContinueTimeout = defaultExpectContinueTimeout
	t.DialContext = DefaultDialer().DialContext
	return t
}

// DefaultDialer creates a new default net.Dialer.
func DefaultDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   defaultTCPTimeout,
		KeepAlive: defaultKeepAlive,
	}
}
