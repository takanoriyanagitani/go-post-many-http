package many

import (
	"context"
)

type SimpleRequest struct {
	Url         string
	ContentType string
	Body        []byte
}

type SimpleRequestSourceMany func(context.Context) (SimpleRequest, error)

func (s SimpleRequestSourceMany) ToChannel(
	ctx context.Context,
	isEnd func(error) bool,
	errorHandler func(error),
) <-chan SimpleRequest {
	var ch chan SimpleRequest = make(chan SimpleRequest)

	go func() {
		defer close(ch)
		for {
			req, e := s(ctx)
			switch e {
			case nil:
				ch <- req
			default:
				if isEnd(e) {
					return
				} else {
					errorHandler(e)
					return
				}
			}
		}
	}()
	return ch
}

type RawRequest []byte

type RawRequestSourceMany func(context.Context) (RawRequest, error)

func (r RawRequestSourceMany) ToSimpleReqSource(
	url string,
	contentType string,
) SimpleRequestSourceMany {
	return func(ctx context.Context) (SimpleRequest, error) {
		raw, e := r(ctx)
		return SimpleRequest{
			Url:         url,
			ContentType: contentType,
			Body:        raw,
		}, e
	}
}
