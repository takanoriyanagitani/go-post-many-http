package hstd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	pm "github.com/takanoriyanagitani/go-post-many-http"
)

type SenderStdST struct {
	*http.Client

	// Converts the [*http.Response].
	ResponseConv func(*http.Response) (pm.TinyResponse, error)
}

// Posts the body to the url.
func (s SenderStdST) Post(
	url string,
	contentType string,
	body []byte,
) (*http.Response, error) {
	return s.Client.Post(url, contentType, bytes.NewReader(body))
}

// Handles the [pm.SimpleRequest] and returns the [pm.TinyResponse].
func (s SenderStdST) Handle(req pm.SimpleRequest) (pm.TinyResponse, error) {
	var url string = req.Url
	var typ string = req.ContentType
	var bdy []byte = req.Body
	res, e := s.Post(url, typ, bdy)
	if nil != e {
		return pm.TinyResponse{}, fmt.Errorf(
			"unable to post(url=%s): %w",
			url,
			e,
		)
	}

	return s.ResponseConv(res)
}

// Converts to the [pm.SenderST].
func (s SenderStdST) ToSender() pm.SenderST {
	return func(
		_ context.Context,
		req pm.SimpleRequest,
	) (pm.TinyResponse, error) {
		return s.Handle(req)
	}
}

// Consumes the response body and convert to the [pm.TinyResponse].
func ResponseToTinyDiscard(
	original *http.Response,
) (pm.TinyResponse, error) {
	var bdy io.ReadCloser = original.Body
	defer bdy.Close()
	_, e := io.Copy(io.Discard, bdy)
	if nil != e {
		return pm.TinyResponse{}, fmt.Errorf("unable to discard body: %w", e)
	}
	return pm.TinyResponse{StatusCode: original.StatusCode}, nil
}

// Creates a handler which rejects any status other than the specified one.
func StatusHandlerNew(acceptStatus int) func(pm.TinyResponse) error {
	return func(res pm.TinyResponse) error {
		var sts int = res.StatusCode
		switch sts {
		case acceptStatus:
			return nil
		default:
			return fmt.Errorf("unexpected http status: %v", sts)
		}
	}
}

var StatusHandlerDefault func(pm.TinyResponse) error = StatusHandlerNew(
	http.StatusOK,
)

func SenderStdNewST(
	client *http.Client,
	responseConverter func(*http.Response) (pm.TinyResponse, error),
) SenderStdST {
	return SenderStdST{
		Client:       client,
		ResponseConv: responseConverter,
	}
}

func SenderStdNewSTdefault(client *http.Client) SenderStdST {
	return SenderStdNewST(
		client,
		ResponseToTinyDiscard,
	)
}
