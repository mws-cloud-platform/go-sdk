package http

import (
	"bytes"
	"io"
	"net/http"
)

type mockRoundTripper struct {
	calls int
	fn    func(req *http.Request, call int) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.calls++
	if m.fn != nil {
		return m.fn(req, m.calls)
	}
	return fakeResponse(req, http.StatusOK, []byte("OK")), nil
}

func fakeResponse(req *http.Request, code int, body []byte) *http.Response {
	var bodyReadCloser io.ReadCloser
	var contentLength int64 = -1

	if len(body) != 0 {
		bodyReadCloser = io.NopCloser(bytes.NewReader(body))
		contentLength = int64(len(body))
	}

	return &http.Response{
		Status:        http.StatusText(code),
		StatusCode:    code,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Uncompressed:  true,
		ContentLength: contentLength,
		Body:          bodyReadCloser,
		Request:       req,
	}
}
