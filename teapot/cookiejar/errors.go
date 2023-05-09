package cookiejar

import (
	"fmt"
)

type ErrorChannel chan error
type ErrorHandler func(err error) bool

type NoCookiesFoundError struct {
	Hosts []string
}

func (e *NoCookiesFoundError) Error() string {
	if len(e.Hosts) == 0 {
		return "No Firefox cookies found"
	}
	return fmt.Sprintf("No Firefox cookies found for: %v", e.Hosts)
}

type NilCookieJarError struct{}

func (e *NilCookieJarError) Error() string {
	return "Nil cookie jar value"
}

type CookieError struct {
	Message string
}

func (e *CookieError) Error() string {
	if e.Message == "" {
		return "Unexpected cookie error"
	}
	return e.Message
}
