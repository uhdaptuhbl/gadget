package exec

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"gadget/halt"
	"gadget/logging"
	"gadget/settings"
	"gadget/sneks"
)

const buildsep = "#"

// TODO: implement default config file writing

// MaxParallelism return conservative number of suggested max parallelism.
func MaxParallelism() int {
	var maxProcs = runtime.GOMAXPROCS(0)
	var numCPU = runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

// Invocation provides the necessary runtime values to the initializer function.
type Invocation struct {
	Name      string
	Version   string
	BuildId   string
	BuildDate string
	Args      []string

	HelpOnEmptyArgs         bool
	NoParseFlags            bool
	ConfigExt               string
	CreateMissingConfigFile bool

	InterruptHandler halt.HandlerFunc
	ShutdownTimeout  time.Duration
	ExitOnError      bool
	ExitCodeError    int

	Sneks sneks.Sneks
	UserDirs settings.UserDirs
}

func (invk *Invocation) Configure(options ...Option) {
	for _, option := range options {
		option(invk)
	}
}

func (invk *Invocation) Build() string {
	return (invk.BuildDate + buildsep + invk.BuildId)
}

func (invk *Invocation) BuildFlags(opts ...settings.FlagFunc) *flag.FlagSet {
	var errBehavior = flag.ContinueOnError

	if invk.ExitOnError {
		errBehavior = flag.ExitOnError
	}
	if len(opts) == 0 {
		opts = invk.DefaultFlagFuncs()
	}

	return settings.Flags.Build(invk.Name, errBehavior, opts...)
}

func (invk *Invocation) BuildViper(flags *flag.FlagSet, ext string, opts ...settings.ViperFunc) (*viper.Viper, error) {
	var err error
	var snek *viper.Viper

	if len(opts) == 0 {
		if opts, err = invk.DefaultViperFuncs(flags, ext); err != nil {
			return snek, err
		}
	}

	snek, err = settings.Viper.Build(opts...)
	if err != nil {
		return snek, err
	}
	if snek == nil {
		return snek, fmt.Errorf("unknown error initializing Viper instance")
	}

	return snek, err
}

// func (invoke Invocation) BuildCommand(opts ...CobraOption) (*cobra.Command, error) {
// 	return nil, nil
// }

func (invk *Invocation) NewLogger(logconf logging.Config) (logging.Logger, error) {
	if log, err := logging.NewZapLogger(logconf); err != nil {
		return log, err
	} else if log == nil {
		return log, fmt.Errorf("unknown error caused nil logger")
	} else {
		return log, err
	}
}

func (invk *Invocation) ParseFlags(flags *flag.FlagSet, ignoreUnknown bool) error {
	settings.Flags.IgnoreUnknown(ignoreUnknown)(flags)
	return flags.Parse(invk.Args)
}

func (invk *Invocation) DefaultFlagFuncs() []settings.FlagFunc {
	// TODO: add a version option and help option if needed

	return []settings.FlagFunc{
		settings.Flags.Sort(false),
		settings.Flags.IgnoreUnknown(false),
		settings.Flags.StringOption(settings.KeyConfigPath, "", settings.HelpConfigPath),
		settings.Flags.StringOption(settings.KeyEnvPrefix, "", settings.HelpEnvPrefix),
		settings.Flags.StringOption(settings.KeyProfileMode, "", settings.HelpProfileMode),
		settings.Flags.BoolOption(settings.KeyVerbose, settings.DefaultVerbose, settings.HelpVerbose),
		settings.Flags.BoolOption(settings.KeyDebug, settings.DefaultDebug, settings.HelpDebug),
		settings.Flags.BoolOption(settings.KeyForce, settings.DefaultForce, settings.HelpForce),
		func(flags *flag.FlagSet) {
			flags.String(settings.KeyLogFormat, settings.DefaultLogFormat, settings.HelpLogFormat)
			flags.String(settings.KeyLogLevel, settings.DefaultLogLevel, settings.HelpLogLevel)
			flags.String(settings.KeyLogVerbosity, settings.DefaultLogVerbosity, settings.HelpLogVerbosity)
			flags.StringSlice(settings.KeyLogOutput, settings.DefaultLogOutputs, settings.HelpLogOutput)
		},
		settings.Flags.Usage(func(flags *flag.FlagSet) {
			var subcommands []string

			// NOTE: Output() is the "correct" way to do it, however, the output
			// writer is not exported via Output() in the current version of pflag
			// so it will not be possible to do it this way until a new version.
			// var dest = flags.Output()
			var dest = os.Stderr

			var sep = "  "
			var name = invk.Name
			var ver = invk.Version
			var build = invk.Build()

			if len(ver) > 0 && ver[0] == 'v' && ver[:2] != "ver" {
				ver = ver[1:]
			}

			var lines = []string{
				"",
				name + sep + ver + sep + build,
				"",
				fmt.Sprintf("USAGE:\n\t%s ARGS [OPTION...]", name),
				"",
				"OPTIONS:",
				flags.FlagUsages(),
			}
			if len(subcommands) > 0 {
				lines = append(lines, "ACTIONS:")
				for _, action := range subcommands {
					lines = append(lines, fmt.Sprintf("\t%s", action))
				}
			}
			var output = strings.Join(lines, "\n") + "\n"

			fmt.Fprint(dest, output)
		}),
	}
}

func (invk *Invocation) DefaultViperFuncs(flags *flag.FlagSet, confType string) ([]settings.ViperFunc, error) {
	var err error
	var opts = []settings.ViperFunc{
		settings.Viper.Name(invk.Name),
		settings.Viper.ConfigName(invk.Name),
		settings.Viper.ConfigPath("."),
	}
	if confType != "" {
		opts = append(opts, settings.Viper.ConfigType(confType))
	}

	// these have to be accessed directly because full config has not been parsed yet
	if flags != nil {
		var configPath string
		var envPrefix string

		if flags.Parsed() {
			if configPath, err = flags.GetString(settings.KeyConfigPath); err != nil {
				return nil, err
			}
			if envPrefix, err = flags.GetString(settings.KeyEnvPrefix); err != nil {
				return nil, err
			}
		} else {
			var useFlags *flag.FlagSet
			var concrete = *flags
			useFlags = &concrete
			settings.Flags.IgnoreUnknown(true)(useFlags)
			if err = invk.ParseFlags(useFlags, true); err != nil {
				if !errors.Is(err, flag.ErrHelp) {
					return nil, err
				}
			}
			if configPath, err = useFlags.GetString(settings.KeyConfigPath); err != nil {
				return nil, err
			}
			if envPrefix, err = useFlags.GetString(settings.KeyEnvPrefix); err != nil {
				return nil, err
			}
		}

		opts = append(opts, settings.Viper.BindPFlags(flags, settings.DefaultPFlagsXform))

		// config file specified on the command line overrides defaults
		if configPath != "" {
			// NOTE: using SetConfigFile() means viper will
			// ignore any other paths added for config.
			opts = append(opts, settings.Viper.ConfigFile(configPath))
		} else {
			opts = append(opts, settings.Viper.ConfigPath(settings.DefaultConfigDir(invk.Name)))
		}

		if envPrefix != "" {
			opts = append(opts, settings.Viper.EnvPrefix(envPrefix))
			opts = append(opts, settings.Viper.AutomaticEnv)
		}
	}

	return opts, err
}
