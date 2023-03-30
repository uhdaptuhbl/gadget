package settings

import (
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Viper = new(ViperFactory)

type ViperFunc func(snek *viper.Viper) error

type ViperFactory struct{}

// Build creates a new viper.Viper instance with any provided options applied.
func (factory ViperFactory) Build(opts ...ViperFunc) (*viper.Viper, error) {
	var snek = viper.New()
	for _, snekfunc := range opts {
		if err := snekfunc(snek); err != nil {
			return nil, err
		}
	}
	return snek, nil
}

func (factory ViperFactory) Configure(snek *viper.Viper, opts ...ViperFunc) (*viper.Viper, error) {
	for _, snekfunc := range opts {
		if err := snekfunc(snek); err != nil {
			return nil, err
		}
	}
	return snek, nil
}

func (factory ViperFactory) AutomaticEnv(snek *viper.Viper) error {
	snek.AutomaticEnv()
	return nil
}

func (factory ViperFactory) Name(name string) ViperFunc {
	return func(snek *viper.Viper) error {
		snek.Set("name", name)
		return nil
	}
}

func (factory ViperFactory) EnvPrefix(prefix string) ViperFunc {
	return func(snek *viper.Viper) error {
		snek.SetEnvPrefix(prefix)
		return nil
	}
}

func (factory ViperFactory) BindPFlags(flags *flag.FlagSet, xform map[string]string) ViperFunc {
	// TODO: instead of the xform map could RegisterAlias be used?
	// https://pkg.go.dev/github.com/spf13/viper#Viper.RegisterAlias
	return func(snek *viper.Viper) error {
		var err error

		if xform == nil {
			err = snek.BindPFlags(flags)
		} else {
			// NOTE: `snek.BindPFlags(flags)` uses the existing names in the
			// flags which can cause a mismatch between value keys from the CLI
			// args vs a config file when loading from both; the xform map
			// allows the caller to specify a mapping of flag keys to config
			// file keys.
			flags.VisitAll(func(fl *flag.Flag) {
				var confkey string
				var ok bool
				if err != nil {
					return
				}
				if _, ok = xform[fl.Name]; ok {
					confkey = xform[fl.Name]
				} else {
					confkey = fl.Name
				}
				// NOTE: this allows the mapping to use empty strings as values
				// to indicate the flag should not be bound at all.
				if confkey != "" {
					if err = snek.BindPFlag(confkey, fl); err != nil {
						err = errors.Wrapf(err, "cannot bind flag: %s <=> %s", confkey, fl.Name)
					}
				}
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (factory ViperFactory) ConfigName(configName string) ViperFunc {
	return func(snek *viper.Viper) error {
		snek.SetConfigName(configName)
		return nil
	}
}

func (factory ViperFactory) ConfigType(configType string) ViperFunc {
	return func(snek *viper.Viper) error {
		// the SetConfigType() func should ignore empty string
		snek.SetConfigType(configType)
		return nil
	}
}

func (factory ViperFactory) ConfigPath(path string) ViperFunc {
	return func(snek *viper.Viper) error {
		snek.AddConfigPath(path)
		return nil
	}
}

func (factory ViperFactory) ConfigFile(configFilePath string) ViperFunc {
	return func(snek *viper.Viper) error {
		snek.SetConfigFile(configFilePath)
		return nil
	}
}
