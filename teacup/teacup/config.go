package teacup

import (
	"crypto/tls"
	"net/http"
	"time"
)

const DefaultUserAgent = "teacup/v0.1.0"

// TLSConfig mirrors most of the tls.Config struct but with tags and serializable.
//
// https://cs.opensource.google/go/go/+/refs/tags/go1.20.3:src/crypto/tls/common.go;l=521
type TLSConfig struct {
	// InsecureSkipVerify controls whether a client verifies the server's
	// certificate chain and host name. If InsecureSkipVerify is true, crypto/tls
	// accepts any certificate presented by the server and any host name in that
	// certificate. In this mode, TLS is susceptible to machine-in-the-middle
	// attacks unless custom verification is used. This should be used only for
	// testing or in combination with VerifyConnection or VerifyPeerCertificate.
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify" json:"insecure_skip_verify,omitempty"`

	// MinVersion contains the minimum TLS version that is acceptable.
	//
	// By default, TLS 1.2 is currently used as the minimum when acting as a
	// client, and TLS 1.0 when acting as a server. TLS 1.0 is the minimum
	// supported by this package, both as a client and as a server.
	//
	// The client-side default can temporarily be reverted to TLS 1.0 by
	// including the value "x509sha1=1" in the GODEBUG environment variable.
	// Note that this option will be removed in Go 1.19 (but it will still be
	// possible to set this field to VersionTLS10 explicitly).
	MinVersion uint16 `mapstructure:"min_version" json:"min_version,omitempty"`
}

func (cfg TLSConfig) Apply(config *tls.Config) {
	if config == nil {
		return
	}
	config.InsecureSkipVerify = cfg.InsecureSkipVerify
	config.MinVersion = cfg.MinVersion
}

// TransportConfig mirrors most of the http.Transport struct but with tags and serializable.
//
// https://cs.opensource.google/go/go/+/refs/tags/go1.20.3:src/net/http/transport.go;l=95
type TransportConfig struct {
	// TLSClientConfig specifies the TLS configuration to use with tls.Client.
	// If nil, the default configuration is used.
	// If non-nil, HTTP/2 support may not be enabled by default.
	TLSClientConfig *tls.Config `mapstructure:"-" json:"-"`

	// ForceAttemptHTTP2 controls whether HTTP/2 is enabled when a non-zero
	// Dial, DialTLS, or DialContext func or TLSClientConfig is provided.
	// By default, use of any those fields conservatively disables HTTP/2.
	// To use a custom dialer or TLS config and still attempt HTTP/2
	// upgrades, set this to true.
	ForceAttemptHTTP2 bool `mapstructure:"force_attempt_http2" json:"force_attempt_http2,omitempty"`

	// TLSHandshakeTimeout specifies the maximum amount of time waiting to
	// wait for a TLS handshake. Zero means no timeout.
	TLSHandshakeTimeout time.Duration `mapstructure:"tls_handshake_timeout" json:"tls_handshake_timeout,omitempty"`

	// DisableKeepAlives, if true, disables HTTP keep-alives and
	// will only use the connection to the server for a single
	// HTTP request.
	//
	// This is unrelated to the similarly named TCP keep-alives.
	DisableKeepAlives bool `mapstructure:"disable_keep_alives" json:"disable_keep_alives,omitempty"`

	// DisableCompression, if true, prevents the Transport from
	// requesting compression with an "Accept-Encoding: gzip"
	// request header when the Request contains no existing
	// Accept-Encoding value. If the Transport requests gzip on
	// its own and gets a gzipped response, it's transparently
	// decoded in the Response.Body. However, if the user
	// explicitly requested gzip it is not automatically
	// uncompressed.
	DisableCompression bool `mapstructure:"disable_compression" json:"disable_compression,omitempty"`

	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConns int `mapstructure:"max_idle_conns" json:"max_idle_conns,omitempty"`

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used.
	MaxIdleConnsPerHost int `mapstructure:"max_idle_conns_per_host" json:"max_idle_conns_per_host,omitempty"`

	// MaxConnsPerHost optionally limits the total number of
	// connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, dials will block.
	//
	// Zero means no limit.
	MaxConnsPerHost int `mapstructure:"max_conns_per_host" json:"max_conns_per_host,omitempty"`

	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleConnTimeout time.Duration `mapstructure:"idle_conn_timeout" json:"idle_conn_timeout,omitempty"`

	// ResponseHeaderTimeout, if non-zero, specifies the amount of
	// time to wait for a server's response headers after fully
	// writing the request (including its body, if any). This
	// time does not include the time to read the response body.
	ResponseHeaderTimeout time.Duration `mapstructure:"response_header_timeout" json:"response_header_timeout,omitempty"`

	// ExpectContinueTimeout, if non-zero, specifies the amount of
	// time to wait for a server's first response headers after fully
	// writing the request headers if the request has an
	// "Expect: 100-continue" header. Zero means no timeout and
	// causes the body to be sent immediately, without
	// waiting for the server to approve.
	// This time does not include the time to send the request header.
	ExpectContinueTimeout time.Duration `mapstructure:"expect_continue_timeout" json:"expect_continue_timeout,omitempty"`

	// ProxyConnectHeader optionally specifies headers to send to
	// proxies during CONNECT requests.
	// To set the header dynamically, see GetProxyConnectHeader.
	ProxyConnectHeader http.Header `mapstructure:"proxy_connect_header" json:"proxy_connect_header,omitempty"`

	// MaxResponseHeaderBytes specifies a limit on how many
	// response bytes are allowed in the server's response
	// header.
	//
	// Zero means to use a default limit.
	MaxResponseHeaderBytes int64 `mapstructure:"max_response_header_bytes" json:"max_response_header_bytes,omitempty"`

	// WriteBufferSize specifies the size of the write buffer used
	// when writing to the transport.
	// If zero, a default (currently 4KB) is used.
	WriteBufferSize int `mapstructure:"write_buffer_size" json:"write_buffer_size,omitempty"`

	// ReadBufferSize specifies the size of the read buffer used
	// when reading from the transport.
	// If zero, a default (currently 4KB) is used.
	ReadBufferSize int
}

func (cfg TransportConfig) Apply(transport *http.Transport) {
	if transport == nil {
		return
	}

	if cfg.TLSClientConfig != nil {
		transport.TLSClientConfig = cfg.TLSClientConfig
	}

	if cfg.ForceAttemptHTTP2 {
		transport.ForceAttemptHTTP2 = cfg.ForceAttemptHTTP2
	}

	if cfg.TLSHandshakeTimeout != 0 {
		transport.TLSHandshakeTimeout = cfg.TLSHandshakeTimeout
	}

	if cfg.DisableKeepAlives {
		transport.DisableKeepAlives = cfg.DisableKeepAlives
	}

	if cfg.DisableCompression {
		transport.DisableCompression = cfg.DisableCompression
	}

	if cfg.MaxIdleConns != 0 {
		transport.MaxIdleConns = cfg.MaxIdleConns
	}

	if cfg.MaxIdleConnsPerHost != 0 {
		transport.MaxIdleConnsPerHost = cfg.MaxIdleConnsPerHost
	}

	if cfg.MaxConnsPerHost != 0 {
		transport.MaxConnsPerHost = cfg.MaxConnsPerHost
	}

	if cfg.IdleConnTimeout != 0 {
		transport.IdleConnTimeout = cfg.IdleConnTimeout
	}

	if cfg.ResponseHeaderTimeout != 0 {
		transport.ResponseHeaderTimeout = cfg.ResponseHeaderTimeout
	}

	if cfg.ExpectContinueTimeout != 0 {
		transport.ExpectContinueTimeout = cfg.ExpectContinueTimeout
	}

	if len(cfg.ProxyConnectHeader) != 0 {
		transport.ProxyConnectHeader = cfg.ProxyConnectHeader
	}

	if cfg.MaxResponseHeaderBytes != 0 {
		transport.MaxResponseHeaderBytes = cfg.MaxResponseHeaderBytes
	}

	if cfg.WriteBufferSize != 0 {
		transport.WriteBufferSize = cfg.WriteBufferSize
	}

	if cfg.ReadBufferSize != 0 {
		transport.ReadBufferSize = cfg.ReadBufferSize
	}
}
