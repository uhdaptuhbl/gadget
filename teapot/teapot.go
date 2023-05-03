package teapot

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	// "github.com/carlmjohnson/requests"
	// "github.com/hashicorp/go-retryablehttp"
	"golang.org/x/net/publicsuffix"
)

// Teapot contains all of the metadata needed to construct a Session.
//
// TODO: initialization from viper instance instead of config
// TODO: add json, toml, yaml, mapstructure, validate tags
type Teapot struct {
	// Headers indicates which HTTP headers should be sent and with
	// what values with every request made by the Session.
	Headers http.Header `mapstructure:"headers" json:"headers,omitempty"`

	// NoCookieJar disables Session cookies when set to true.
	NoCookieJar bool `mapstructure:"no_cookie_jar" json:"no_cookie_jar,omitempty"`

	// Timeout specifies a time limit for requests made by this
	// Client. The timeout includes connection time, any
	// redirects, and reading the response body. The timer remains
	// running after Get, Head, Post, or Do return and will
	// interrupt reading of the Response.Body.
	//
	// A Timeout of zero means no timeout. The Client cancels requests
	// to the underlying Transport as if the Request's Context ended.
	Timeout time.Duration `mapstructure:"timeout" json:"timeout,omitempty"`

	// Transport represents the `http.Transport` configuration for the `http.Client`.
	Transport *TransportConfig `mapstructure:"transport" json:"transport,omitempty"`

	// TLSConfig represents the `TLSClientConfig` field of the `http.Transport`.
	TLS *TLSConfig `mapstructure:"tls" json:"tls,omitempty"`

	onRequest  []RequestInterceptor
	onResponse []ResponseInterceptor
	transport  *http.Transport
	tlsconfig  *tls.Config
	cookiejar  *cookiejar.Jar
	httpclient *http.Client

	// UserAgent provides a convenience way to specify the User-Agent
	// HTTP header without needing to specify other headers.
	// UserAgent string `mapstructure:"user_agent" json:"user_agent,omitempty"`
	// UserAgentFunc func() string `mapstructure:"-" json:"-"`
}

func (teapot *Teapot) Clone() *Teapot {
	var clone = &Teapot{
		// UserAgentFunc: teapot.UserAgentFunc,
		NoCookieJar: teapot.NoCookieJar,
		Timeout:     teapot.Timeout,
		Transport:   teapot.Transport,
		TLS:         teapot.TLS,

		onRequest:  append(make([]RequestInterceptor, 0, len(teapot.onRequest)), teapot.onRequest...),
		onResponse: append(make([]ResponseInterceptor, 0, len(teapot.onResponse)), teapot.onResponse...),
		httpclient: teapot.httpclient,
		transport: func() *http.Transport {
			if teapot.transport != nil {
				return teapot.transport.Clone()
			}
			return nil
		}(),
		tlsconfig: func() *tls.Config {
			if teapot.tlsconfig != nil {
				return teapot.tlsconfig.Clone()
			}
			return nil
		}(),
	}

	return clone
}

func (teapot *Teapot) OnRequest(handlers ...RequestInterceptor) *Teapot {
	teapot.onRequest = append(teapot.onRequest, handlers...)
	return teapot
}

func (teapot *Teapot) OnResponse(handlers ...ResponseInterceptor) *Teapot {
	teapot.onResponse = append(teapot.onResponse, handlers...)
	return teapot
}

func (teapot *Teapot) WithTransport(transport *http.Transport) *Teapot {
	teapot.transport = transport
	return teapot
}

func (teapot *Teapot) WithTLS(tlsconfig *tls.Config) *Teapot {
	teapot.tlsconfig = tlsconfig
	return teapot
}

func (teapot *Teapot) WithCookieJar(jar *cookiejar.Jar) *Teapot {
	teapot.cookiejar = jar
	return teapot
}

func (teapot *Teapot) Session() Session {
	if teapot.Headers == nil {
		teapot.Headers = make(http.Header)
	}

	// if _, ok := teapot.Headers[http.CanonicalHeaderKey("User-Agent")]; ok {
	// 	if teapot.UserAgent == "" {
	// 		teapot.UserAgent = teapot.Headers.Get("User-Agent")
	// 	} else if teapot.UserAgent != teapot.Headers.Get("User-Agent") {
	// 		teapot.Headers.Set("User-Agent", teapot.UserAgent)
	// 	}
	// } else {
	// 	if teapot.UserAgent == "" {
	// 		teapot.UserAgent = DefaultUserAgent
	// 		teapot.Headers.Set("User-Agent", DefaultUserAgent)
	// 	} else {
	// 		teapot.Headers.Set("User-Agent", teapot.UserAgent)
	// 	}
	// }

	return NewTeacup(teapot)
}

func (teapot *Teapot) Client() *http.Client {
	// TODO: teapot must support retries!!! and be configurable for statuses that should be retried and how many times

	if teapot.httpclient == nil {
		// A new cookie jar is always created for a new client unless disabled.
		var err error
		var jar *cookiejar.Jar
		if !teapot.NoCookieJar {
			// https://golangbyexample.com/set-cookie-http-golang/
			// https://husni.dev/manage-http-cookie-in-go-with-cookie-jar/
			if teapot.cookiejar != nil {
				jar = teapot.cookiejar
			} else if jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}); err != nil {
				// NOTE: As of Go 1.16, cookiejar.New err is hardcoded nil:
				// https://cs.opensource.google/go/go/+/refs/tags/go1.20.3:src/net/http/cookiejar/jar.go;l=85
				panic(fmt.Sprintf("As of Go 1.16, cookiejar.New err is supposed to be hardcoded nil: %v", err))
			}
		}

		if teapot.tlsconfig == nil {
			teapot.tlsconfig = &tls.Config{
				// NOTE: `gosec` linter complains if this is not explicitly
				// set even though TLS 1.2 is already the default it uses.
				MinVersion: tls.VersionTLS12,

				// NOTE: `gosec` linter complains if this is not set to false
				InsecureSkipVerify: false,
			}
			if teapot.TLS != nil {
				teapot.TLS.Apply(teapot.tlsconfig)
			}

			// TODO
			// // NOTE: at this time this wrapper doesn't support less than TLS 1.2
			// if teapot.TLS.TLSMinVersion > tls.VersionTLS12 {
			// 	teapot.TLS.MinVersion = teapot.TLS.TLSMinVersion
			// }
		}

		if teapot.transport == nil {
			teapot.transport = http.DefaultTransport.(*http.Transport).Clone()

			// set these to more sane defaults
			teapot.transport.MaxIdleConns = 100
			teapot.transport.MaxConnsPerHost = 100
			teapot.transport.MaxIdleConnsPerHost = 100

			// not attempted by default when TLSClientConfig is set
			teapot.transport.ForceAttemptHTTP2 = true

			if teapot.tlsconfig != nil {
				teapot.transport.TLSClientConfig = teapot.tlsconfig
			}

			if teapot.Transport != nil {
				teapot.Transport.Apply(teapot.transport)
			}
		}

		teapot.httpclient = &http.Client{Transport: teapot.transport, Timeout: teapot.Timeout, Jar: jar}
	}

	return teapot.httpclient
}
