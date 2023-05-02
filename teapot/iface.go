package teapot

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/carlmjohnson/requests"
)

type BuilderMutator interface {
	Mutate(f func(builder *requests.Builder))
}

type Requestor interface {
	URLstring(loc string) Requestor
	URL(loc *url.URL) Requestor
	Headers(headers http.Header) Requestor
	Body(body io.Reader) Requestor

	Request(ctx context.Context, method string) *Result
	Head(ctx context.Context) *Result
	Get(ctx context.Context) *Result
	Post(ctx context.Context) *Result
	Put(ctx context.Context) *Result
	Patch(ctx context.Context) *Result
	Delete(ctx context.Context) *Result
	Options(ctx context.Context) *Result
}

type Session interface {
	Requestor

	Clone() Session
	Client() *http.Client
}
