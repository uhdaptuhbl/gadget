package harness

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"bitbucket.services.cymru.com/FST/gadget/halt"
	"bitbucket.services.cymru.com/FST/gadget/settings"
	"bitbucket.services.cymru.com/voltron/logging"
)

const KeyConfigPath = "config"
const KeyEnvPrefix = "env-prefix"
const KeyProfileMode = "profile-mode"
const KeyVerbosity = "verbosity"
const KeyDebug = "debug"
const KeyForce = "force"

// const KeyLogLevel = "loglevel"
// const KeyLogFormat = "logformat"
// const KeyLogOutput = "logoutput"

const HelpConfigPath = "Specify `<path>` to config file."
const HelpEnvPrefix = "Set a `<prefix>` for environment variables."

var HelpProfileMode = "Set the `<profile>` mode: {" + PrettyProfileModes() + "}."
var HelpVerbosity = "Set `<verbosity>` of output: {" + logging.PrettyLogVerbosities() + "}."

const HelpDebug = "Enable debug output."
const HelpForce = "Perform potentially destructive actions."

// const HelpLogLevel = "minimum logging level"
// const HelpLogFormat = "logging format"
// const HelpLogOutput = "logging output file paths"

var DefaultDebug = false
var DefaultForce = false
var DefaultVerbosity = string(logging.LogVerbositySimple)

// var DefaultLogLevel = string(logging.LogLevelInfo)
// var DefaultLogFormat = string(logging.LogFormatJSON)
// var DefaultLogOutputs = []string{"stdout"}

const ExitCodeError = 42
const buildsep = "#"

var DefaultPFlagsXform = map[string]string{
	KeyConfigPath:  "",
	KeyEnvPrefix:   "",
	KeyProfileMode: "profile_mode",
	KeyVerbosity:   "logging.verbosity",
	// KeyLogLevel:     "logging.level",
	// KeyLogFormat:    "logging.format",
	// KeyLogOutput:    "logging.outputpaths",
}

type ConfigChangeFunc func() error

/*
Program interface provides the contract for long-running daemon applications.
*/
type Program interface {
	Flags() *flag.FlagSet
	Viper() *viper.Viper
	Load() error
	OnConfigChange() ConfigChangeFunc
	ProfileMode() string
	HandleSignals(ctx context.Context, cancel context.CancelFunc, sigch chan os.Signal) error
	Run(ctx context.Context) error
}

/*
ProgramBase can be embedded in your Program struct to avoid having to reimplement functions not used.
*/
type ProgramBase struct{}

func (program *ProgramBase) Flags() *flag.FlagSet             { return nil }
func (program *ProgramBase) Viper() *viper.Viper              { return nil }
func (program *ProgramBase) Load() error                      { return nil }
func (program *ProgramBase) OnConfigChange() ConfigChangeFunc { return nil }
func (program *ProgramBase) ProfileMode() string              { return "" }
func (program *ProgramBase) HandleSignals(ctx context.Context, cancel context.CancelFunc, sigch chan os.Signal) error {
	signal.Notify(sigch, syscall.SIGHUP, syscall.SIGTERM, os.Interrupt)

	return halt.HandleInterrupts(halt.InterruptOptions{
		Signalch:        sigch,
		Context:         ctx,
		Shutdown:        cancel,
		GracefulTimeout: time.Second * 5,
		TimeoutExitCode: ExitCodeError,
		// TODO: is there a good way to pass this a logger instance?
		// Log:             program.Logger(),
	})
}

/*
InitializeFunc initializes a new Program instance.
*/
type InitializeFunc func(args InvokeArgs) (Program, error)

/*
InvokeArgs provides the necessary runtime values to the initializer function.
*/
type InvokeArgs struct {
	Name      string
	Args      []string
	Version   string
	BuildID   string
	BuildDate string

	ExitCodeError   int
	ShutdownTimeout time.Duration

	ExitOnError     bool
	NoParseFlags    bool
	HelpOnEmptyArgs bool

	// TODO: where should the config ext be stored / passed from?
	ConfigExt               string
	CreateMissingConfigFile bool

	InterruptHandler halt.HandlerFunc
}

// TODO: implement default config file writing

func (iArgs InvokeArgs) Build() string {
	return (iArgs.BuildDate + buildsep + iArgs.BuildID)
}

func (iArgs InvokeArgs) BuildFlags(opts ...settings.FlagFunc) *flag.FlagSet {
	var errBehavior = flag.ContinueOnError

	if iArgs.ExitOnError {
		errBehavior = flag.ExitOnError
	}
	if len(opts) == 0 {
		opts = iArgs.defaultFlagFuncs()
	}

	return settings.Flags.Build(iArgs.Name, errBehavior, opts...)
}

func (iArgs InvokeArgs) BuildViper(flags *flag.FlagSet, ext string, opts ...settings.ViperFunc) (*viper.Viper, error) {
	var err error
	var snek *viper.Viper

	if len(opts) == 0 {
		if opts, err = iArgs.defaultViperFuncs(flags, ext); err != nil {
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

func (iArgs InvokeArgs) NewLogger(logconf logging.Config) (logging.Logger, error) {
	if log, err := logging.NewZapLogger(logconf); err != nil {
		return log, err
	} else if log == nil {
		return log, fmt.Errorf("unknown error caused nil logger")
	} else {
		return log, err
	}
}

func (iArgs InvokeArgs) ParseFlags(flags *flag.FlagSet, ignoreUnknown bool) error {
	settings.Flags.IgnoreUnknown(ignoreUnknown)(flags)
	return flags.Parse(iArgs.Args)
}

func (iArgs InvokeArgs) defaultFlagFuncs() []settings.FlagFunc {
	// TODO: add a version option and help option if needed

	return []settings.FlagFunc{
		settings.Flags.Sort(false),
		settings.Flags.IgnoreUnknown(false),
		settings.Flags.StringOption(KeyConfigPath, "", HelpConfigPath),
		settings.Flags.StringOption(KeyEnvPrefix, "", HelpEnvPrefix),
		settings.Flags.StringOption(KeyProfileMode, "", HelpProfileMode),
		settings.Flags.StringOption(KeyVerbosity, DefaultVerbosity, HelpVerbosity),
		settings.Flags.BoolOption(KeyDebug, DefaultDebug, HelpDebug),
		settings.Flags.BoolOption(KeyForce, DefaultForce, HelpForce),
		// func(flags *flag.FlagSet) {
		// 	flags.String(KeyLogLevel, defaultLogLevel, helpLogLevel)
		// 	flags.String(KeyLogFormat, defaultLogFormat, helpLogFormat)
		// 	flags.StringSlice(KeyLogOutput, defaultLogOutputs, helpLogOutput)
		// },
		settings.Flags.Usage(func(flags *flag.FlagSet) {
			var subcommands []string

			// NOTE: Output() is the "correct" way to do it, however, the output
			// writer is not exported via Output() in the current version of pflag
			// so it will not be possible to do it this way until a new version.
			// var dest = flags.Output()
			var dest = os.Stderr

			var sep = "  "
			var name = iArgs.Name
			var ver = iArgs.Version
			var build = iArgs.Build()

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

func (iArgs InvokeArgs) defaultViperFuncs(flags *flag.FlagSet, confType string) ([]settings.ViperFunc, error) {
	var err error
	var opts = []settings.ViperFunc{
		settings.Viper.Name(iArgs.Name),
		settings.Viper.ConfigName(iArgs.Name),
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
			if configPath, err = flags.GetString(KeyConfigPath); err != nil {
				return nil, err
			}
			if envPrefix, err = flags.GetString(KeyEnvPrefix); err != nil {
				return nil, err
			}
		} else {
			var useFlags *flag.FlagSet
			var concrete = *flags
			useFlags = &concrete
			settings.Flags.IgnoreUnknown(true)(useFlags)
			if err = iArgs.ParseFlags(useFlags, true); err != nil {
				if !errors.Is(err, flag.ErrHelp) {
					return nil, err
				}
			}
			if configPath, err = useFlags.GetString(KeyConfigPath); err != nil {
				return nil, err
			}
			if envPrefix, err = useFlags.GetString(KeyEnvPrefix); err != nil {
				return nil, err
			}
		}

		opts = append(opts, settings.Viper.BindPFlags(flags, DefaultPFlagsXform))

		// config file specified on the command line overrides defaults
		if configPath != "" {
			// NOTE: using SetConfigFile() means viper will
			// ignore any other paths added for config.
			opts = append(opts, settings.Viper.ConfigFile(configPath))
		} else {
			opts = append(opts, settings.Viper.ConfigPath(DefaultConfigDir(iArgs.Name)))
		}

		if envPrefix != "" {
			opts = append(opts, settings.Viper.EnvPrefix(envPrefix))
			opts = append(opts, settings.Viper.AutomaticEnv)
		}
	}

	return opts, err
}
