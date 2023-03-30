package teapot

import (
	"fmt"
	"net/url"
)

type UserAgentError struct {
	UserAgent string
}

func (e *UserAgentError) Error() string {
	return fmt.Sprintf("invalid user agent: '%s'", e.UserAgent)
}

type StatusCodeError struct {
	StatusCode int
	Expected   int
}

func (e *StatusCodeError) Error() string {
	if e.Expected == 0 {
		return fmt.Sprintf("unexpected HTTP status code: '%d'", e.StatusCode)
	}
	return fmt.Sprintf("unexpected HTTP status code: '%d' (expected '%d')", e.StatusCode, e.Expected)
}

type ContentTypeError struct {
	ContentType string
	Expected    string
}

func (e *ContentTypeError) Error() string {
	if e.Expected == "" {
		return fmt.Sprintf("unsupported HTTP content type: '%s'", e.ContentType)
	}
	return fmt.Sprintf("unsupported HTTP content type: '%s' (expected '%s')", e.ContentType, e.Expected)
}

type RequestFailedError struct {
	Status   string
	Location *url.URL
	Reason   string
}

func (e *RequestFailedError) Error() string {
	if e.Location == nil {
		return fmt.Sprintf("%s %s", e.Status, e.Reason)
	}
	return fmt.Sprintf("%s %s %s", e.Status, e.Location.String(), e.Reason)
}
