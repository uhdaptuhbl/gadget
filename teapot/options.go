package teapot

import (
	"crypto/tls"
	"net/http"

	"gadget/logging"
)

type Option func(tpt *teapot)

func UseLogger(log logging.Logger) Option {
	return func(tpt *teapot) {
		tpt.log = log
	}
}

func UseConfig(cfg *Config) Option {
	return func(tpt *teapot) {
		tpt.config = cfg
	}
}

func UseTransport(transport *http.Transport) Option {
	return func(tpt *teapot) {
		tpt.transport = transport
	}
}

func UseTLS(tlsconfig *tls.Config) Option {
	return func(tpt *teapot) {
		tpt.tlsconfig = tlsconfig
	}
}

func UseHeaders(headers http.Header, raw bool) Option {
	return func(tpt *teapot) {
		if raw {
			tpt.headers = headers
		} else {
			if tpt.headers == nil {
				tpt.headers = make(http.Header)
			}
			CopyHeaders(tpt.headers, headers, false)
		}
	}
}

func UseCookieJar(jar http.CookieJar) Option {
	return func(tpt *teapot) {
		tpt.cookiejar = jar
	}
}

func UseRequestHandlers(handlers []RequestInterceptor) Option {
	return func(tpt *teapot) {
		tpt.onRequest = append(tpt.onRequest, handlers...)
	}
}

func UseResponseHandlers(handlers []ResponseInterceptor) Option {
	return func(tpt *teapot) {
		tpt.onResponse = append(tpt.onResponse, handlers...)
	}
}
