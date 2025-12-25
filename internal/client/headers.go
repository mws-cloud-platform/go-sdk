package client

import (
	"net/http"
	"strings"
)

func Header(headers http.Header, name string) string {
	for headerName := range headers {
		if strings.EqualFold(headerName, name) {
			return headers.Get(headerName)
		}
	}
	return ""
}
