package teacup

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"time"

	// "github.com/carlmjohnson/requests"
	// "github.com/hashicorp/go-retryablehttp"
	"golang.org/x/net/publicsuffix"
)

// Teacup contains all of the metadata needed to construct a Session.
//
// TODO: initialization from viper instance instead of config
// TODO: add json, toml, yaml, mapstructure, validate tags
type Teacup struct {
	// UserAgent provides a convenience way to specify the User-Agent
	// HTTP header without needing to specify other headers.
	UserAgent string `mapstructure:"user_agent" json:"user_agent,omitempty"`

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
	Transport *TransportConfig `mapstructure:"transport,squash" json:"transport,omitempty"`

	// TLSConfig represents the `TLSClientConfig` field of the `http.Transport`.
	TLS *TLSConfig `mapstructure:"tls,squash" json:"tls,omitempty"`

	transport  *http.Transport
	tlsconfig  *tls.Config
	httpclient *http.Client
	// interceptors []int // TODO: add interceptors
	// options []Option
}

func (teacup *Teacup) Clone() *Teacup {
	var clone = &Teacup{
		UserAgent: teacup.UserAgent,
		NoCookieJar: teacup.NoCookieJar,
		Timeout: teacup.Timeout,
		Transport: teacup.Transport,
		TLS: teacup.TLS,
		httpclient: teacup.httpclient,
		transport: func() *http.Transport {
			if teacup.transport != nil {
				return teacup.transport.Clone()
			}
			return nil
		}(),
		tlsconfig: func() {
			if teacup.tlsconfig != nil {
				return teacup.tlsconfig.Clone()
			}
			return nil
		}(),
	}
	return clone
}

func (teacup *Teacup) WithTransport(transport *http.Transport) *Teacup {
	teacup.transport = transport
	return teacup
}

func (teacup *Teacup) WithTLS(tlsconfig *tls.Config) *Teacup {
	teacup.tls = tlsconfig
	return teacup
}

func (teacup *Teacup) WithCookieJar(jar *cookiejar.Jar) *Teacup {
	teacup.cookies = jar
	return teacup
}

func (teacup *Teacup) WithRetries(retry string) *Teacup {
	// TODO: implement retries for the client construction
	return teacup
}

func (teacup *Teacup) Session() (*teacupSession, error) {
	var err error

	if teacup.Headers == nil {
		teacup.Headers = make(http.Header)
	}
	if _, ok := teacup.Headers[http.CanonicalHeaderKey("User-Agent")]; ok {
		if teacup.UserAgent == "" {
			teacup.UserAgent = teacup.Headers.Get("User-Agent")
		} else if teacup.UserAgent != teacup.Headers.Get("User-Agent") {
			teacup.Headers.Set("User-Agent", teacup.UserAgent)
		}
	} else {
		if teacup.UserAgent == "" {
			teacup.UserAgent = UserAgent
			teacup.Headers.Set("User-Agent", UserAgent)
		} else {
			teacup.Headers.Set("User-Agent", teacup.UserAgent)
		}
	}

	var session = &teacupSession{
		teacup: teacup,
		// client: teacup.client(),
		// requestor: teacup.requestor(),
	}
	teacup.cookies = nil

	return session, err
}

func (teacup *Teacup) client() *http.Client {
	// TODO: cache the client? use it as a singleton for any new sessions once created?
	// TODO: teacup must support retries!!! and be configurable for statuses that should be retried and how many times

	if teacup.httpclient == nil {
		// A new cookie jar is always created for a new client unless disabled.
		var jar *cookiejar.Jar
		if !teacup.NoCookieJar {
			// https://golangbyexample.com/set-cookie-http-golang/
			// https://husni.dev/manage-http-cookie-in-go-with-cookie-jar/
			if jar, err = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}); err != nil {
				// NOTE: As of Go 1.16, cookiejar.New err is hardcoded nil:
				// https://cs.opensource.google/go/go/+/refs/tags/go1.20.3:src/net/http/cookiejar/jar.go;l=85
				panic(fmt.Sprintf("As of Go 1.16, cookiejar.New err is supposed to be hardcoded nil: %v", err))
			}
		}

		if teacup.tlsconfig == nil {
			teacup.tlsconfig = &tls.Config{
				// NOTE: `gosec` linter complains if this is not explicitly
				// set even though TLS 1.2 is already the default it uses.
				MinVersion: tls.VersionTLS12,

				// NOTE: `gosec` linter complains if this is not set to false
				InsecureSkipVerify: false,
			}
			teacup.TLS.Apply(teacup.tlsconfig)

			// TODO
			// // NOTE: at this time this wrapper doesn't support less than TLS 1.2
			// if teacup.TLS.TLSMinVersion > tls.VersionTLS12 {
			// 	teacup.TLS.MinVersion = teacup.TLS.TLSMinVersion
			// }
		}

		if teacup.transport == nil {
			teacup.transport = http.DefaultTransport.(*http.Transport).Clone()

			// set these to more sane defaults
			teacup.transport.MaxIdleConns = 100
			teacup.transport.MaxConnsPerHost = 100
			teacup.transport.MaxIdleConnsPerHost = 100

			// not attempted by default when TLSClientConfig is set
			teacup.transport.ForceAttemptHTTP2 = true

			if teacup.tlsconfig != nil {
				teacup.transport.TLSClientConfig = teacup.tlsconfig
			}

			// TODO: is there a default tls config that it's already using?
			teacup.Transport.Apply(teacup.transport)
		}

		teacup.httpclient = &http.Client{Transport: teacup.transport, Timeout: teacup.Timeout, Jar: jar}
	}

	return teacup.httpclient
}
