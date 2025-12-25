package client

import (
	"context"
	"mime"
	"net/http"

	"go.mws.cloud/go-sdk/mws/errors"
	"go.mws.cloud/go-sdk/pkg/context/metadata"
)

func GetContentType(resp *http.Response) (string, error) {
	rawContentType := resp.Header.Get("Content-Type")

	if rawContentType == "" {
		return "", nil
	}

	ct, _, err := mime.ParseMediaType(rawContentType)
	if err != nil {
		return "", errors.NewInvalidContentTypeError(rawContentType)
	}

	return ct, nil
}

func AddOutgoingMetadataToHeader(ctx context.Context, req *http.Request) {
	md := metadata.FromOutgoingContext(ctx)
	for key, values := range md {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}
