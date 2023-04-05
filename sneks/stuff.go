package sneks

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"gadget/logging"
)

type CobraOption func(cmd *cobra.Command)

type FlagOption func(flags *flag.FlagSet)

type ViperOption func(vip *viper.Viper)

type Sneks struct {
	Cobra *cobra.Command
	Viper *viper.Viper
}

func RunE(ec int, f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := f(cmd, args); err != nil {
			logging.Fatalf(ec, "%v", err)
		}
	}
}
