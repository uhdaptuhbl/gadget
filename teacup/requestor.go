package teacup

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Requestor interface {
	Send(ctx context.Context) *Result
	Headers(headers http.Header) Requestor
	BodyBytes(body []byte) Requestor
	BodyReader(body io.Reader) Requestor
	Request(method string, loc *url.URL) Requestor
	Head(loc *url.URL) Requestor
	Get(loc *url.URL) Requestor
	Post(loc *url.URL) Requestor
	Put(loc *url.URL) Requestor
	Patch(loc *url.URL) Requestor
	Delete(loc *url.URL) Requestor
	Options(loc *url.URL) Requestor
}

type requestor struct {
	client   *http.Client
	location *url.URL
	method   string
	body     io.Reader
	headers  http.Header
	retryStatus []int
	blacklistStatus []int
	whitelistStatus []int
}

// TODO: add additional functional options for unmarshalling json,
// checking headers like content-length, etc. so add hooks for intercepting
// the request and responses and propagating errors, maybe outside of the
// result even to avoid the need for result.Error != nil checks.
func (builder *requestor) fetch(ctx context.Context) *Result {
	var err error
	var req *http.Request
	var resp *http.Response
	var body []byte
	var result Result

	var loc = builder.location.String()
	if req, err = http.NewRequestWithContext(ctx, strings.ToUpper(builder.method), loc, builder.body); err != nil {
		result.Error = err
		return &result
	}
	CopyHeaders(builder.headers, req.Header, false)
	result.Request = req

	if resp, err = builder.client.Do(req); err != nil {
		result.Error = err
		return &result
	}
	result.Response = resp

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

	// immediately reading the response ensures the body
	// will be closed and the connection released
	defer resp.Body.Close()
	if body, err = io.ReadAll(resp.Body); err != nil {
		result.Error = err
		return &result
	}
	result.Body = body

	return &result
}

func (builder *requestor) Send(ctx context.Context) *Result {
	var result = new(Result)

	return result
}

func (builder *requestor) Headers(headers http.Header) Requestor {
	builder.headers = headers
	return builder
}

func (builder *requestor) Body(body io.Reader) Requestor {
	builder.body = body
	return builder
}

func (builder *requestor) Request(method string, loc *url.URL) Requestor {
	builder.method = method
	builder.location = loc
	return builder
}

func (builder *requestor) Head(loc *url.URL) Requestor {
	builder.method = "HEAD"
	builder.location = loc
	return builder
}

func (builder *requestor) Get(loc *url.URL) Requestor {
	builder.method = "GET"
	builder.location = loc
	return builder
}

func (builder *requestor) Post(loc *url.URL) Requestor {
	builder.method = "POST"
	builder.location = loc
	return builder
}

func (builder *requestor) Put(loc *url.URL) Requestor {
	builder.method = "PUT"
	builder.location = loc
	return builder
}

func (builder *requestor) Patch(loc *url.URL) Requestor {
	builder.method = "PATCH"
	builder.location = loc
	return builder
}

func (builder *requestor) Delete(loc *url.URL) Requestor {
	builder.method = "DELETE"
	builder.location = loc
	return builder
}

func (builder *requestor) Options(loc *url.URL) Requestor {
	builder.method = "OPTIONS"
	builder.location = loc
	return builder
}
