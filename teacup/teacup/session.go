package teacup

import (
	"net/http"
)

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

type Session interface {
	Headers() http.Header
	Client() *http.Client
	Request(loc *url.URL) Requestor

	*requestor
}

type teacupSession struct {
	teacup *Teacup

	*requestor
}

// func (session *teacupSession) builder() *requestor {
// 	if session.requestor == nil {
// 		session.requestor = &requestor{client: session.client}
// 	}
// 	return session.requestor
// }

func (session *teacupSession) Headers() http.Header {
	return session.teacup.Headers
}

func (session *teacupSession) Client() *http.Client {
	return session.teacup.httpclient
}

func (session *teacupSession) Request(loc *url.URL) Requestor {
	// TODO
	return nil
}

func (teacup *Teacup) Requestor() Requestor {
	// build a new one each time
	return &requestor{client: teacup.client()}
}
