package many

import (
	"context"
)

// Sends a [SimpleRequest] and returns a [TinyResponse].
type SenderST func(context.Context, SimpleRequest) (TinyResponse, error)

func (s SenderST) SendMany(
	ctx context.Context,
	requests <-chan SimpleRequest,
	handleResponse func(TinyResponse) error,
) (uint64, error) {
	var cnt uint64 = 0
	for req := range requests {
		res, e := s(ctx, req)
		if nil != e {
			return 0, e
		}

		e = handleResponse(res)
		if nil != e {
			return 0, e
		}
		cnt += 1
	}
	return cnt, nil
}

func (s SenderST) SendManyEx(
	ctx context.Context,
	requests <-chan SimpleRequest,
	handleResponse func(TinyResponse, SimpleRequest) error,
) (uint64, error) {
	var cnt uint64 = 0
	for req := range requests {
		res, e := s(ctx, req)
		if nil != e {
			return 0, e
		}

		e = handleResponse(res, req)
		if nil != e {
			return 0, e
		}
		cnt += 1
	}
	return cnt, nil
}
