package harness

import (
	"context"
	"fmt"
	"gadget/exec"
	"os"
	"runtime"

	"github.com/pkg/errors"
	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"gadget/logging"
)

// Invoke runs cobra commands along with boilerplate.
func Invoke(initialize func(*exec.Invocation) (exec.Application, error), options ...exec.Option) {
	var err error
	var invoke exec.Invocation
	var app exec.Application
	var cmd *cobra.Command

	defer func() {
		// TODO: is recover() actually needed to get a stack trace?
		// TODO: move this to helper in exec module
		if r := recover(); r != nil {
			var buf = make([]byte, 1<<10)
			runtime.Stack(buf, false)
			_, _ = fmt.Fprintf(os.Stderr, "[ERROR] recovered in Invoke() from panic:%v\n%s\n", r, string(buf))
		}
	}()

	for _, option := range options {
		option(&invoke)
	}

	if app, err = initialize(&invoke); err != nil {
		logging.Fatalf(invoke.ExitCodeError, "error initializing program: %v", err)
	}

	if cmd = app.Command(); cmd == nil {
		logging.Fatalf(invoke.ExitCodeError, "nil command: %v", app)
	}

	// if flags := cmd.Flags(); flags != nil && !invoke.NoParseFlags {
	// 	settings.Flags.IgnoreUnknown(false)(flags)
	// 	if err = flags.Parse(invoke.Args); err != nil {
	// 		if errors.Is(err, flag.ErrHelp) {
	// 			os.Exit(0)
	// 		}
	// 		logging.Fatalf(invoke.ExitCodeError, "parsing runtime options failed: %v", err)
	// 	}
	// }

	// https://github.com/carolynvs/stingoftheviper/blob/main/main.go
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// NOTE: cobra and viper can be bound in a few locations,
		// but PersistencePreRunE on the root command works well.
		// NOTE: app.Load() should have access to the cmd and args
		// from the invocation object it got when initialized.
		if err = app.Load(cmd, args); err != nil {
			logging.Fatalf(invoke.ExitCodeError, "app.Load() failed: %v", err)
		}
		return nil
	}

	// if invoke.HelpOnEmptyArgs && len(invoke.Args) == 0 {
	// 	if flags := app.Flags(); flags != nil {
	// 		if flags.Usage != nil {
	// 			flags.Usage()
	// 		} else {
	// 			flags.PrintDefaults()
	// 		}
	// 	}
	// 	os.Exit(0)
	// }

	// if invoke.CreateMissingConfigFile {
	// 	if err != nil && !errors.Is(err, os.ErrNotExist) {
	// 		LogFatal(invoke.ExitCodeError, fmt.Sprintf("%v", err))
	// 	}
	// 	// TODO: default config file writing
	// 	// var createdPath string
	// 	// if createdPath, err = config.WriteDefaultConfigFile(); err != nil {
	// 	// 	if !errors.Is(err, os.ErrExist) {
	// 	// 		LogFatal(fmt.Sprintf("%v", err))
	// 	// 	}
	// 	// }
	// 	// if errors.Is(err, os.ErrExist) {
	// 	// 	LogInfof("config path already exists: %s", createdPath)
	// 	// } else {
	// 	// 	LogInfof("created config file: %s", createdPath)
	// 	// }

	// 	// // this can still error if --config is specified and is different from default
	// 	// if err = conf.Load(); err != nil {
	// 	// 	LogFatal(fmt.Sprintf("%v", err))
	// 	// }
	// }

	switch app.ProfileMode() {
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

	var interruptch = make(chan os.Signal, 1)

	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var g, gctx = errgroup.WithContext(ctx)
	g.Go(func() error {
		return app.HandleSignals(ctx, cancel, interruptch)
	})
	g.Go(func() error {
		var runerr error
		// TODO: is this the right way to go about terminating the signal
		// handler? Would it be more idiomatic to close the channel if that works?
		defer cancel()

		if runerr = cmd.ExecuteContext(gctx); runerr != nil {
			runerr = errors.Wrap(runerr, "app.Run()")
		}
		return runerr
	})
	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		logging.Fatalf(invoke.ExitCodeError, "%v", err)
	}
}
