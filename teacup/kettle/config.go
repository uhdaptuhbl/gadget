package teapot

import (
	"time"
)

type Config struct {
	Bindport     int           `mapstructure:"bind_port" json:"bind_port"`
	BindAddress  string        `mapstructure:"bind_address" json:"bind_address"`
	CAPath       string        `mapstructure:"capath" json:"capath"`
	CertPath     string        `mapstructure:"certpath" json:"certpath"`
	KeyPath      string        `mapstructure:"keypath" json:"keypath"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout"`
}
