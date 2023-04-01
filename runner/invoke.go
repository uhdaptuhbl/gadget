package harness

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/pkg/profile"
	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"

	"gadget/settings"
)

const ProfileCPU = "cpu"
const ProfileMem = "mem"
const ProfileMutex = "mutex"
const ProfileBlock = "block"
const ProfileTrace = "trace"

func ProfileModes() []string {
	return []string{
		string(ProfileCPU),
		string(ProfileMem),
		string(ProfileMutex),
		string(ProfileBlock),
		string(ProfileTrace),
	}
}

func PrettyProfileModes() string {
	return strings.Join(ProfileModes(), ",")
}

/*
Execute runs the main program with interrupt handlers and config loading.

TODO: convert this all to Cobra since it auto-handles help messages and such
*/
func Invoke(initialize InitializeFunc, invokeArgs InvokeArgs) {
	var err error
	var program Program

	var interruptch = make(chan os.Signal, 1)

	defer func() {
		// TODO: is recover() actually needed to get a stack trace?
		if r := recover(); r != nil {
			var buf = make([]byte, 1<<10)
			runtime.Stack(buf, false)
			_, _ = fmt.Fprintf(os.Stderr, "[ERROR] recovered in Invoke() from panic:%v\n%s\n", r, string(buf))
		}
	}()

	if program, err = initialize(invokeArgs); err != nil {
		LogFatalf(fmt.Sprintf("%v", err))
	}

	if flags := program.Flags(); flags != nil && !invokeArgs.NoParseFlags {
		settings.Flags.IgnoreUnknown(false)(flags)
		if err = flags.Parse(invokeArgs.Args); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				os.Exit(0)
			}
			LogFatalf("parsing runtime options failed: %v", err)
		}
	}

	if err = program.Load(); err != nil {
		LogFatalf("program.Load() failed: %v", err)
	}

	if invokeArgs.HelpOnEmptyArgs && len(invokeArgs.Args) == 0 {
		if flags := program.Flags(); flags != nil {
			if flags.Usage != nil {
				flags.Usage()
			} else {
				flags.PrintDefaults()
			}
		}
		os.Exit(0)
	}

	if invokeArgs.CreateMissingConfigFile {
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			LogFatal(fmt.Sprintf("%v", err))
		}
		// TODO: default config file writing
		// var createdPath string
		// if createdPath, err = config.WriteDefaultConfigFile(); err != nil {
		// 	if !errors.Is(err, os.ErrExist) {
		// 		LogFatal(fmt.Sprintf("%v", err))
		// 	}
		// }
		// if errors.Is(err, os.ErrExist) {
		// 	LogInfof("config path already exists: %s", createdPath)
		// } else {
		// 	LogInfof("created config file: %s", createdPath)
		// }

		// // this can still error if --config is specified and is different from default
		// if err = conf.Load(); err != nil {
		// 	LogFatal(fmt.Sprintf("%v", err))
		// }
	}

	var updateConfig = program.OnConfigChange()
	if snek := program.Viper(); snek != nil && updateConfig != nil && snek.ConfigFileUsed() != "" {
		snek.OnConfigChange(func(e fsnotify.Event) {
			if err = updateConfig(); err != nil {
				interruptch <- syscall.SIGTERM
			}
		})
		snek.WatchConfig()
	} else if updateConfig != nil && snek == nil {
		LogFatalf("OnConfigChange() specified but Viper instance is nil")
	}

	switch program.ProfileMode() {
	case "cpu":
		defer profile.Start(profile.CPUProfile).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile).Stop()
	case "mutex":
		defer profile.Start(profile.MutexProfile).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile).Stop()
	case "trace":
		defer profile.Start(profile.TraceProfile).Stop()
	}

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var g, gctx = errgroup.WithContext(ctx)
	g.Go(func() error {
		return program.HandleSignals(ctx, cancel, interruptch)
	})
	g.Go(func() error {
		var runerr error
		// TODO: is this the right way to go about terminating the signal
		// handler? Would it be more idiomatic to close the channel if that works?
		defer cancel()

		if runerr = program.Run(gctx); runerr != nil {
			runerr = errors.Wrap(runerr, "program.Run()")
		}
		return runerr
	})
	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		LogFatalf("%v", err)
		os.Exit(ExitCodeError)
	}
}
