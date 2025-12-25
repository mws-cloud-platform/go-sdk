package client

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/go-faster/jx"

	"go.mws.cloud/go-sdk/mws/errors"
)

// Marshaler implements json writing.
type Marshaler interface {
	Encode(e *jx.Encoder)
}

func SetJSONBody(req *http.Request, body Marshaler) error {
	var (
		e   jx.Encoder
		buf = &bytes.Buffer{}
	)

	if body == nil {
		return nil
	}

	body.Encode(&e)

	if _, err := buf.Write(e.Bytes()); err != nil {
		return errors.NewEncodeBodyError(err)
	}

	req.Body = io.NopCloser(buf)
	req.ContentLength = int64(buf.Len())

	setContentTypeRequest(req, "application/json")

	return nil
}

func SetJSONBytesBody(req *http.Request, body []byte) error {
	return SetJSONBody(req, ByteArray(body))
}

func SetBytesBody(req *http.Request, body []byte, contentType string) error {
	buf := &bytes.Buffer{}
	if _, err := buf.Write(body); err != nil {
		return errors.NewEncodeBodyError(err)
	}

	req.Body = io.NopCloser(buf)
	req.ContentLength = int64(buf.Len())
	setContentTypeRequest(req, contentType)

	return nil
}

func setContentTypeRequest(req *http.Request, contentType string) {
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
}

type ByteArray []byte

func (b ByteArray) Encode(e *jx.Encoder) {
	e.Base64(b)
}

// URLFormMarshaler implements url form writing.
type URLFormMarshaler interface {
	MarshalURLForm() (url.Values, error)
}

func SetFormURLEncodedBody(req *http.Request, body URLFormMarshaler) error {
	if body == nil {
		return nil
	}

	urlForm, err := body.MarshalURLForm()
	if err != nil {
		return errors.NewEncodeBodyError(err)
	}

	buf := &bytes.Buffer{}

	if _, err = buf.WriteString(urlForm.Encode()); err != nil {
		return errors.NewEncodeBodyError(err)
	}

	req.Body = io.NopCloser(buf)
	req.ContentLength = int64(buf.Len())

	setContentTypeRequest(req, "application/x-www-form-urlencoded")

	return nil
}
