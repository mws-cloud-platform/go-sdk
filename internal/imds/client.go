package imds

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.mws.cloud/util-toolset/pkg/os/env"
)

const (
	metadataHost    = "http://169.254.169.254"
	metadataHostEnv = "MWS_COMPUTE_METADATA_HOST"
	metadataPath    = "/computeMetadata/v1/"

	metadataTimeout    = 500 * time.Millisecond
	metadataTimeoutEnv = "MWS_COMPUTE_METADATA_CLIENT_TIMEOUT"

	metadataFlavorHeader = "Metadata-Flavor"
	metadataFlavorValue  = "Google"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPClient
	env        env.Env
}

func NewClient(client HTTPClient, env env.Env) Client {
	return Client{
		httpClient: client,
		env:        env,
	}
}

func (p Client) GetWithContext(ctx context.Context, key string) (_ string, err error) {
	host := metadataHost
	if hostEnv, ok := p.env.LookupEnv(metadataHostEnv); ok {
		host = hostEnv
	}

	timeout := metadataTimeout
	if timeoutEnv, ok := p.env.LookupEnv(metadataTimeoutEnv); ok {
		timeout, err = time.ParseDuration(timeoutEnv)
		if err != nil {
			return "", fmt.Errorf("parse compute metadata client timeout: %w", err)
		}
	}

	key = strings.TrimLeft(key, "/")
	host = strings.TrimRight(host, "/")
	addr := host + metadataPath + key

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set(metadataFlavorHeader, metadataFlavorValue)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		cErr := resp.Body.Close()
		if cErr != nil {
			err = errors.Join(err, cErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", Error{
			Code:    resp.StatusCode,
			Message: string(body),
		}
	}

	return string(body), nil
}

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (code: %d)", e.Message, e.Code)
}
