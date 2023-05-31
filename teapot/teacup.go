package teapot

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/carlmjohnson/requests"

	"gadget/logging"
)

// type BuilderMutator interface {
// 	Mutate(f func(builder *requests.Builder))
// }

func newBuilder() *requests.Builder {
	return requests.New() //.Client(client)
}

func newTeacup() *teacup {
	return &teacup{Builder: newBuilder()}
}

// func newTeacup(c *http.Client, j http.CookieJar, h http.Header) *teacup {
// 	return &teacup{client: c, jar: j, header: h, Builder: requests.New().Client(c)}
// }

// func NewSession(c *http.Client, opts ...func(tcup *teacup)) *teacup {
// 	var tcup = &teacup{client: c, Builder: requests.New().Client(c)}
// 	for _, opt := range opts {
// 		opt(tcup)
// 	}
// 	return tcup
// }

type teacup struct {
	*requests.Builder

	log        logging.Logger
	client     *http.Client
	jar        http.CookieJar
	headers    http.Header
	onRequest  []RequestInterceptor
	onResponse []ResponseInterceptor

	location *url.URL
	method   string
	body     io.Reader
	header   http.Header
	err      error

	// retryStatus []int
	// blacklistStatus []int
	// whitelistStatus []int

}

func (session *teacup) clone() *teacup {
	var loc url.URL
	var locptr *url.URL

	if session.location != nil {
		loc = *session.location
		locptr = &loc
	}

	var clone = &teacup{
		Builder: session.Builder,

		log:    session.log,
		client: session.client,
		jar:    session.jar,
		headers: func() http.Header {
			if session.headers == nil {
				return make(http.Header)
			}
			return session.headers.Clone()
		}(),
		onRequest:  session.onRequest[:],
		onResponse: session.onResponse[:],

		location: locptr,
		method:   session.method,
		body:     session.body,
		header: func() http.Header {
			if session.header == nil {
				return make(http.Header)
			}
			return session.header.Clone()
		}(),
		err: session.err,
	}
	// CopyHeaders(clone.headers, session.headers, true)
	return clone
}

func (session *teacup) Clone() Session {
	return session.clone()
}

func (session *teacup) Mutate() SessionMutator {
	return newSessionMutator(session)
}

func (session *teacup) Client() *http.Client {
	return session.client
}

func (session *teacup) Jar() http.CookieJar {
	return session.jar
}

// func (session *teacup) SetHeader(key string, value string) {
// 	session.headers.Set(key, value)
// }

func (session *teacup) URLstring(loc string) Requestor {
	var clone = session.clone()
	clone.location, clone.err = url.Parse(loc)
	return clone
}

func (session *teacup) URL(loc *url.URL) Requestor {
	var clone = session.clone()
	clone.location = loc
	return clone
}

func (session *teacup) Headers(headers http.Header) Requestor {
	var clone = session.clone()
	clone.headers = headers
	return clone
}

func (session *teacup) Body(body io.Reader) Requestor {
	var clone = session.clone()
	clone.body = body
	return clone
}

func (session *teacup) Request(ctx context.Context, method string) *Result {
	session.method = method
	return session.fetch(ctx)
}

func (session *teacup) Head(ctx context.Context) *Result {
	session.method = "HEAD"
	return session.fetch(ctx)
}

func (session *teacup) Get(ctx context.Context) *Result {
	session.method = "GET"
	return session.fetch(ctx)
}

func (session *teacup) Post(ctx context.Context) *Result {
	session.method = "POST"
	return session.fetch(ctx)
}

func (session *teacup) Put(ctx context.Context) *Result {
	session.method = "PUT"
	return session.fetch(ctx)
}

func (session *teacup) Patch(ctx context.Context) *Result {
	session.method = "PATCH"
	return session.fetch(ctx)
}

func (session *teacup) Delete(ctx context.Context) *Result {
	session.method = "DELETE"
	return session.fetch(ctx)
}

func (session *teacup) Options(ctx context.Context) *Result {
	session.method = "OPTIONS"
	return session.fetch(ctx)
}

func (session *teacup) fetch(ctx context.Context) *Result {
	// TODO: add additional functional options for unmarshalling json,
	// checking headers like content-length, etc. so add hooks for intercepting
	// the request and responses and propagating errors, maybe outside of the
	// result even to avoid the need for result.Error != nil checks.

	var err error
	var req *http.Request
	var resp *http.Response
	var body []byte
	var result Result

	if session.err != nil {
		result.Error = session.err
		return &result
	}

	var loc = session.location.String()
	if req, err = http.NewRequestWithContext(ctx, strings.ToUpper(session.method), loc, session.body); err != nil {
		result.Error = err
		return &result
	}
	result.Request = req

	for _, handler := range session.onRequest {
		if err = handler(req); err != nil {
			result.Error = err
			return &result
		}
	}

	// TODO: should this be done here or before applying onRequest?
	CopyHeaders(req.Header, session.headers, false)

	if resp, err = session.Client().Do(req); err != nil {
		result.Error = err
		return &result
	}
	defer resp.Body.Close()
	result.Response = resp

	// immediately reading the response ensures the body
	// will be closed and the connection released
	if body, err = io.ReadAll(resp.Body); err != nil {
		result.Error = err
		return &result
	}
	result.Body = body

	for _, handler := range session.onResponse {
		if err = handler(resp); err != nil {
			result.Error = err
			return &result
		}
	}

	return &result
}

// TODO: should gzip be handled manually?
// TODO: https://stackoverflow.com/questions/71011274/golang-default-http-client-doesnt-handle-compression
// var body *strings.Reader
// var reader *bufio.Reader
// 	contentEncoding := resp.Header.Get("Content-Encoding")
// 	if strings.Contains(contentEncoding, "gzip") {
// 		var gzipReader *gzip.Reader
// 		if gzipReader, err = gzip.NewReader(resp.Body); err != nil {
// 			return stats, errors.Wrap(err, "runTask() gzip.NewReader() error")
// 		}
// 		defer gzipReader.Close()
// 		reader = bufio.NewReader(gzipReader)
// 	} else {
// 		reader = bufio.NewReader(resp.Body)
// 	}
// 	if resp.StatusCode != 200 {
// 		log.Println(task.LogFormat("[ERROR] runTask failed with StatusCode:%d", resp.StatusCode))
// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			return stats, errors.Wrap(err, fmt.Sprintf("runTask() failed, StatusCode:%d", resp.StatusCode))
// 		}
// 		return stats, errors.New(fmt.Sprintf("StatusCode:%d body:%s", resp.StatusCode, string(body)))
// 	}
