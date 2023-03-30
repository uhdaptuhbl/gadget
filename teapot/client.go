package teapot

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"

	"gadget/logging"
)

type HTTPClientFactory func() (*Client, error)
type HTTPClientAccessor func() *Client

func DumpResult(log logging.Logger, result *Result) {
	log.Debugf("REQUEST: %+v", result.Request)
	log.Debugf("RESPONSE: %+v", result.Response)
	log.Debugf("CONTENT: %v", string(result.Body))
}

func CopyHeaders(src http.Header, dst http.Header, overwrite bool) {
	for h, vals := range src {
		if overwrite {
			dst.Del(h)
		}
		for _, val := range vals {
			dst.Add(h, val)
		}
	}
}

type RequestArgs struct {
	Ctx    context.Context
	Method string
	// TODO: should this be called Location instead to avoid naming confusion?
	URL     url.URL
	Body    io.Reader
	Headers http.Header
}

/*
Result provides info to the caller about the request, response, and content.

This allows the caller the option to inspect the actual http.Request,
http.Response, and raw body content as needed, without usage assumptions.
*/
type Result struct {
	Request  *http.Request
	Response *http.Response
	Body     []byte
	Error    error
}

func (res *Result) StatusCode() int {
	if res == nil || res.Response == nil {
		return 0
	}
	return res.Response.StatusCode
}
func (res *Result) StatusMessage() string {
	var msg strings.Builder

	if res == nil {
		msg.WriteString("-  -  -")
	} else {
		if res.Response != nil {
			msg.WriteString(res.Response.Status)
		} else {
			msg.WriteString("-")
		}
		msg.WriteString("  ")
		if res.Request != nil {
			msg.WriteString(res.Request.Method)
		} else {
			msg.WriteString("-")
		}
		msg.WriteString("  ")
		if res.Request != nil {
			msg.WriteString(res.Request.URL.String())
		} else {
			msg.WriteString("-")
		}
	}

	return msg.String()
}
func (res *Result) Content() string {
	if res == nil {
		return ""
	}
	return string(res.Body)
}
func (res *Result) JSON(dest interface{}) error {
	if res == nil {
		return nil
	}
	return json.Unmarshal(res.Body, dest)
}
func (res *Result) Dump(logfunc func(string, ...interface{}), msg string) {
	logfunc(
		msg,
		"request",
		fmt.Sprintf("%+v", res.Request),
		"response",
		fmt.Sprintf("%+v", res.Response),
		"body", string(res.Body),
	)
}

func (res *Result) History() []*http.Request {
	var resp *http.Response
	var history []*http.Request

	resp = res.Response
	for resp != nil {
		req := resp.Request
		history = append(history, req)
		resp = req.Response
	}
	return history
}

func (res *Result) Locations() []url.URL {
	var locations []url.URL

	for _, req := range res.History() {
		locations = append(locations, *req.URL)
	}

	return locations
}

type Client struct {
	client    *http.Client
	transport *http.Transport
	tlsconfig *tls.Config
	headers   http.Header
	userAgent string
}

func New(conf ClientConfig) *Client {
	var cl = &Client{
		userAgent: conf.UserAgent,
	}

	if conf.Headers == nil {
		cl.headers = make(http.Header)
	} else {
		cl.headers = conf.Headers
	}
	if _, ok := cl.headers[http.CanonicalHeaderKey("User-Agent")]; true {
		if ok {
			cl.userAgent = cl.headers.Get("User-Agent")
		} else if !ok && conf.UserAgent != "" {
			cl.headers.Set("User-Agent", conf.UserAgent)
		}
	}

	cl.tlsconfig = &tls.Config{
		// false by default otherwise gosec complains
		InsecureSkipVerify: false,

		// linter gosec complains if this is not explicitly set
		// even though TLS 1.2 is already the default it uses
		MinVersion: tls.VersionTLS12,
	}
	cl.tlsconfig.InsecureSkipVerify = conf.InsecureSkipVerify
	// at this time this wrapper doesn't support less than TLS 1.2
	if conf.TLSMinVersion > tls.VersionTLS12 {
		cl.tlsconfig.MinVersion = conf.TLSMinVersion
	}

	cl.transport = &http.Transport{
		TLSClientConfig:       cl.tlsconfig,
		TLSHandshakeTimeout:   conf.TLSHandshakeTimeout,
		ResponseHeaderTimeout: conf.ResponseHeaderTimeout,
		ExpectContinueTimeout: conf.ExpectContinueTimeout,
		IdleConnTimeout:       conf.IdleConnTimeout,
		MaxIdleConns:          conf.MaxIdleConns,
		MaxIdleConnsPerHost:   conf.MaxIdleConnsPerHost,
		MaxConnsPerHost:       conf.MaxConnsPerHost,
	}

	cl.client = &http.Client{
		Timeout:   conf.Timeout,
		Transport: cl.transport,
	}

	return cl
}

func (cl *Client) StoreCookies(remember bool) error {
	if remember && cl.client.Jar == nil {
		var err error
		// https://golangbyexample.com/set-cookie-http-golang/
		// https://husni.dev/manage-http-cookie-in-go-with-cookie-jar/
		if cl.client.Jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}); err != nil {
			return err
		}
	} else if !remember && cl.client.Jar != nil {
		cl.client.Jar = nil
	}
	return nil
}

// TODO: functional options for body and headers?
func (cl *Client) GET(ctx context.Context, location url.URL, body io.Reader, headers http.Header) *Result {
	return cl.MakeRequest(RequestArgs{
		Ctx:     ctx,
		Method:  "GET",
		URL:     location,
		Body:    body,
		Headers: headers,
	})
}

// TODO: functional options for body and headers?
func (cl *Client) POST(ctx context.Context, location url.URL, body io.Reader, headers http.Header) *Result {
	return cl.MakeRequest(RequestArgs{
		Ctx:     ctx,
		Method:  "POST",
		URL:     location,
		Body:    body,
		Headers: headers,
	})
}

// TODO: functional options for body and headers?
func (cl *Client) PUT(ctx context.Context, location url.URL, body io.Reader, headers http.Header) *Result {
	return cl.MakeRequest(RequestArgs{
		Ctx:     ctx,
		Method:  "PUT",
		URL:     location,
		Body:    body,
		Headers: headers,
	})
}

// TODO: functional options for body and headers?
func (cl *Client) DELETE(ctx context.Context, location url.URL, body io.Reader, headers http.Header) *Result {
	return cl.MakeRequest(RequestArgs{
		Ctx:     ctx,
		Method:  "DELETE",
		URL:     location,
		Body:    body,
		Headers: headers,
	})
}

// TODO: add additional functional options for unmarshalling json,
// checking headers like content-length, etc.
func (cl *Client) MakeRequest(args RequestArgs) *Result {
	var err error
	var req *http.Request
	var resp *http.Response
	var body []byte
	var result Result

	var location = args.URL.String()
	if req, err = http.NewRequestWithContext(args.Ctx, strings.ToUpper(args.Method), location, args.Body); err != nil {
		result.Error = err
		return &result
	}
	CopyHeaders(cl.headers, req.Header, false)
	CopyHeaders(args.Headers, req.Header, true)
	result.Request = req

	if resp, err = cl.doRequest(req); err != nil {
		result.Error = err
		return &result
	}
	result.Response = resp

	// immediately reading the response ensures the body will be closed
	if body, err = cl.readResponseBody(resp); err != nil {
		result.Error = err
		return &result
	}
	result.Body = body

	return &result
}

func (cl *Client) doRequest(req *http.Request) (*http.Response, error) {
	var err error
	var resp *http.Response

	if resp, err = cl.client.Do(req); err != nil {
		return resp, err
	}

	return resp, err
}

// readResponseBody function ensures the body is closed each time it's read
func (cl *Client) readResponseBody(resp *http.Response) ([]byte, error) {
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

	var err error
	var body []byte
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	return body, err
}
