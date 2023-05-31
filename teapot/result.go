package teapot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func CopyHeaders(dst http.Header, src http.Header, overwrite bool) {
	for h, vals := range src {
		if overwrite {
			dst.Del(h)
		}
		for _, val := range vals {
			dst.Add(h, val)
		}
	}
}

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

func (res *Result) Dump() string {
	// fmt.Sprintf("\nREQUEST: %+v\nRESPONSE: %+v\nBODY: %s", res.Response, res.Request, res.Body)
	var dump strings.Builder
	if res.Request != nil {
		// dump.WriteString("\nREQUEST: " + res.Request.URL.String() + "\n" + res.Request.Method)

		dump.WriteString("\n" + res.Response.Proto)
		dump.WriteString("\n" + res.Request.Method + " " + res.Request.URL.String())
		for header, values := range res.Request.Header {
			dump.WriteString("\n\t" + header + ": " + fmt.Sprintf("%v", values))
		}
	}
	if res.Response != nil {
		// dump.WriteString("\nRESPONSE: " + fmt.Sprintf("%+v", res.Response) + "\n")

		dump.WriteString("\n" + res.Response.Proto)
		dump.WriteString("\n" + res.Response.Status)
		dump.WriteString(" " + res.Response.Request.URL.String())
		for header, values := range res.Response.Header {
			dump.WriteString("\n\t" + header + ": " + fmt.Sprintf("%v", values))
		}
	}
	// dump.WriteString(fmt.Sprintf("\nBODY: %d bytes", len(res.Body)))
	dump.WriteString(fmt.Sprintf("\nBODY: %s", res.Text()))
	return dump.String()
}

func (res *Result) DumpLog(logfunc func(string, ...any), msg string) {
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
