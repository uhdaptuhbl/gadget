package teapot

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"gadget/logging"
)

type RequestMutator interface {
	URLstring(loc string) Requestor
	URL(loc *url.URL) Requestor
	Headers(headers http.Header) Requestor
	Body(body io.Reader) Requestor
}

type Requestor interface {
	RequestMutator

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
	Mutate() SessionMutator
	Client() *http.Client
	Jar() http.CookieJar
}

type SessionOption func(tcup *teacup)

type SessionMutator interface {
	Logger(log logging.Logger) SessionMutator
	AddHeaders(headers http.Header) SessionMutator
	SetHeaders(headers http.Header) SessionMutator
	CookieJar(jar http.CookieJar) SessionMutator
	OnRequest(handlers ...RequestInterceptor) SessionMutator
	OnResponse(handlers ...ResponseInterceptor) SessionMutator
	Make() Session
}

type sessionMutator struct {
	tcup *teacup              `mapstructure:"-" json:"-"`
	opts []func(tcup *teacup) `mapstructure:"-" json:"-"`
}

func newSessionMutator(tcup *teacup, opts ...func(*teacup)) SessionMutator {
	return &sessionMutator{tcup: tcup, opts: append(make([]func(*teacup), 0, 10), opts...)}
}

func (mttr *sessionMutator) Logger(log logging.Logger) SessionMutator {
	mttr.opts = append(mttr.opts, func(tcup *teacup) { tcup.log = log })
	return mttr
}

func (mttr *sessionMutator) AddHeaders(headers http.Header) SessionMutator {
	mttr.opts = append(mttr.opts, func(tcup *teacup) {
		if tcup.headers == nil {
			tcup.headers = make(http.Header)
		}
		CopyHeaders(tcup.headers, headers, false)
	})
	return mttr
}

func (mttr *sessionMutator) SetHeaders(headers http.Header) SessionMutator {
	mttr.opts = append(mttr.opts, func(tcup *teacup) { tcup.headers = headers.Clone() })
	return mttr
}

func (mttr *sessionMutator) CookieJar(jar http.CookieJar) SessionMutator {
	mttr.opts = append(mttr.opts, func(tcup *teacup) { tcup.jar = jar })
	return mttr
}

func (mttr *sessionMutator) OnRequest(handlers ...RequestInterceptor) SessionMutator {
	// NOTE: OnRequest will ADD handlers to what already exists!
	mttr.opts = append(mttr.opts, func(tcup *teacup) { tcup.onRequest = append(tcup.onRequest, handlers...) })
	return mttr
}

func (mttr *sessionMutator) OnResponse(handlers ...ResponseInterceptor) SessionMutator {
	// NOTE: OnResponse will ADD handlers to what already exists!
	mttr.opts = append(mttr.opts, func(tcup *teacup) { tcup.onResponse = append(tcup.onResponse, handlers...) })
	return mttr
}

// Make will clone the teapot then modify and return the clone.
func (mttr *sessionMutator) Make() Session {
	var tcup *teacup
	if mttr.tcup != nil {
		tcup = mttr.tcup.clone()
	} else {
		tcup = newTeacup()
	}
	for _, optfunc := range mttr.opts {
		optfunc(tcup)
	}
	return tcup
}
