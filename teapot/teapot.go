package teapot

import (
	"crypto/tls"
	"net/http"
	"time"

	"gadget/logging"
	"gadget/teapot/cookiejar"
)

type Teapot interface {
	Clone() Teapot
	Mutate() Constructor

	Session() Session
	OnRequest() []RequestInterceptor
	OnResponse() []ResponseInterceptor

	Client() *http.Client
	// Jar() http.CookieJar
}

// Teapot contains all of the metadata needed to construct a Session.
//
// TODO: initialization from viper instance instead of config?
type teapot struct {
	log        logging.Logger
	config     *Config
	transport  *http.Transport
	tlsconfig  *tls.Config
	headers    http.Header
	httpclient *http.Client
	cookiejar  http.CookieJar
	jarLoaders []cookiejar.Loader
	onRequest  []RequestInterceptor
	onResponse []ResponseInterceptor
}

func (tpt *teapot) clone() *teapot {
	var clone = &teapot{
		log:    tpt.log,
		config: tpt.config,
		transport: func() *http.Transport {
			if tpt.transport != nil {
				return tpt.transport.Clone()
			}
			return nil
		}(),
		tlsconfig: func() *tls.Config {
			if tpt.tlsconfig != nil {
				return tpt.tlsconfig.Clone()
			}
			return nil
		}(),
		headers:    tpt.headers.Clone(),
		httpclient: tpt.httpclient,
		cookiejar:  tpt.cookiejar,
		jarLoaders: append(make([]cookiejar.Loader, 0, len(tpt.jarLoaders)), tpt.jarLoaders...),
		onRequest:  append(make([]RequestInterceptor, 0, len(tpt.onRequest)), tpt.onRequest...),
		onResponse: append(make([]ResponseInterceptor, 0, len(tpt.onResponse)), tpt.onResponse...),
	}

	return clone
}

func (tpt *teapot) Clone() Teapot {
	return tpt.clone()
}

func (tpt *teapot) Mutate() Constructor {
	return constructorFromTeapot(tpt)
}

func (tpt *teapot) Session() Session {
	var cup = newTeacup()
	cup.client = tpt.Client()
	cup.jar = tpt.cookiejar
	cup.headers = tpt.headers.Clone()
	if cup.jar == nil && cup.client.Jar != nil {
		cup.jar = cup.client.Jar
	} else if cup.jar == nil && cup.client.Jar == nil {
		cup.jar = cookiejar.New()
	} else if cup.client.Jar == nil && cup.jar != nil {
		cup.client.Jar = cup.jar
	}
	if tpt.onRequest != nil {
		cup.onRequest = tpt.onRequest[:]
	}
	if tpt.onResponse != nil {
		cup.onResponse = tpt.onResponse[:]
	}
	return cup
}

func (tpt *teapot) OnRequest() []RequestInterceptor {
	return tpt.onRequest
}

func (tpt *teapot) OnResponse() []ResponseInterceptor {
	return tpt.onResponse
}

func (tpt *teapot) Client() *http.Client {
	// TODO: teapot must support retries!!! and be configurable for statuses that should be retried and how many times

	if tpt.httpclient == nil {
		if tpt.config == nil {
			tpt.config = &Config{Timeout: time.Second * 10}
		}

		// A new cookie jar is always created for a new client unless disabled or provided.
		var jar http.CookieJar
		if tpt.cookiejar != nil {
			jar = tpt.cookiejar
		} else if !tpt.config.NoCookieJar {
			// TODO:
			jar = cookiejar.New()
		}

		// ddb.teacup.SetHeader("User-Agent", ddb.teapot.UserAgent())
		// ddb.teapot = teapot.NewClient(teapot.ClientConfig{
		// 	Header:                ddb.headers,
		// 	Timeout:               30 * time.Second,
		// 	TLSHandshakeTimeout:   3 * time.Second,
		// 	InsecureSkipVerify:    true,
		// 	ResponseHeaderTimeout: 5 * time.Second,
		// 	ExpectContinueTimeout: 20 * time.Second,
		// 	IdleConnTimeout:       30 * time.Second,
		// 	MaxIdleConns:          10,
		// 	MaxIdleConnsPerHost:   5,
		// 	MaxConnsPerHost:       5,
		// })

		if tpt.tlsconfig == nil {
			tpt.tlsconfig = &tls.Config{
				// NOTE: `gosec` linter complains if this is not explicitly
				// set even though TLS 1.2 is already the default it uses.
				MinVersion: tls.VersionTLS12,

				// NOTE: `gosec` linter complains if this is not set to false
				InsecureSkipVerify: false,
			}
			if tpt.config != nil && tpt.config.TLS != nil {
				tpt.config.TLS.Apply(tpt.tlsconfig)
			}

			// TODO
			// // NOTE: at this time this wrapper doesn't support less than TLS 1.2
			// if tpt.TLS.TLSMinVersion > tls.VersionTLS12 {
			// 	tpt.TLS.MinVersion = tpt.TLS.TLSMinVersion
			// }
		}

		if tpt.transport == nil {
			tpt.transport = http.DefaultTransport.(*http.Transport).Clone()

			// set these to more sane defaults
			tpt.transport.MaxIdleConns = 100
			tpt.transport.MaxConnsPerHost = 100
			tpt.transport.MaxIdleConnsPerHost = 100

			// not attempted by default when TLSClientConfig is set
			tpt.transport.ForceAttemptHTTP2 = true

			if tpt.tlsconfig != nil {
				tpt.transport.TLSClientConfig = tpt.tlsconfig
			}

			if tpt.config != nil && tpt.config.Transport != nil {
				tpt.config.Transport.Apply(tpt.transport)
			}
		}

		tpt.httpclient = &http.Client{Transport: tpt.transport, Timeout: tpt.config.Timeout, Jar: jar}
	}

	return tpt.httpclient
}
