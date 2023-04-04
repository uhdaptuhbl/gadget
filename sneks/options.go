package sneks

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CobraOption func(cmd *cobra.Command)

type FlagOption func(flags *flag.FlagSet)

type ViperOption func(vip *viper.Viper)

type Sneks struct {
	Cobra *cobra.Command
	Viper *viper.Viper
}
