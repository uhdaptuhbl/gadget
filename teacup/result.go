package teacup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Result provides info to the caller about the request, response, and content.
//
// This allows the caller the option to inspect the actual http.Request,
// http.Response, and raw body content as needed, without usage assumptions.
//
// WARNING: this should not be used for large downloads as the entire response
// body is read into a byte slice at once without chunking or lazy-loading.
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

	if res == nil || res.Response == nil {
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

func (res *Result) Content() []byte {
	if res == nil {
		return make([]byte, 0)
	}
	return res.Body
}

func (res *Result) Text() string {
	if res == nil {
		return ""
	}
	return string(res.Body)
}

func (res *Result) JSON(dest any) error {
	if res == nil {
		return nil
	}
	return json.Unmarshal(res.Body, dest)
}

func (res *Result) Dump(logfunc func(string, ...any), msg string) {
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

func (res *Result) Locations() []*url.URL {
	var locations []*url.URL

	for _, req := range res.History() {
		locations = append(locations, req.URL)
	}

	return locations
}
