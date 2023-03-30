package teapot

import (
	"time"
)

// TODO: is this a use-case for functional options? seemse like it might be
// since there's quite a few things that could be specified for http client
// type Config struct {
// 	Client Client `mapstructure:"client" json:"client"`
// 	Server Server `mapstructure:"server" json:"server"`
// }

type ClientConfig struct {
	UserAgent             string        `mapstructure:"user_agent" json:"user_agent"`
	Headers               http.Header   `mapstructure:"headers" json:"headers"`
	Timeout               time.Duration `mapstructure:"timeout" json:"timeout"`
	TLSHandshakeTimeout   time.Duration `mapstructure:"tls_handshake_timeout" json:"tls_handshake_timeout"`
	InsecureSkipVerify    bool          `mapstructure:"insecure_skip_verify" json:"insecure_skip_verify"`
	TLSMinVersion         uint16        `mapstructure:"tls_min_version" json:"tls_min_version"`
	ResponseHeaderTimeout time.Duration `mapstructure:"response_header_timeout" json:"response_header_timeout"`
	ExpectContinueTimeout time.Duration `mapstructure:"continue_timeout" json:"continue_timeout"`
	IdleConnTimeout       time.Duration `mapstructure:"idle_conn_timeout" json:"idle_conn_timeout"`
	MaxIdleConns          int           `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	MaxIdleConnsPerHost   int           `mapstructure:"max_idle_conns_per_host" json:"max_idle_conns_per_host"`
	MaxConnsPerHost       int           `mapstructure:"max_conns_per_host" json:"max_conns_per_host"`
}
// TODO: untangle aliases
type Config = ClientConfig

type ServerConfig struct {
	Bindport     int           `mapstructure:"bind_port" json:"bind_port"`
	BindAddress  string        `mapstructure:"bind_address" json:"bind_address"`
	CAPath       string        `mapstructure:"capath" json:"capath"`
	CertPath     string        `mapstructure:"certpath" json:"certpath"`
	KeyPath      string        `mapstructure:"keypath" json:"keypath"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout"`
}
