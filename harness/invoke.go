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
	"gadget/settings"
)

// Invoke runs cobra commands along with boilerplate.
func Invoke(initialize func(*exec.Invocation) (exec.Application, error), options ...exec.Option) {
	var err error
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

	var invoke = new(exec.Invocation)
	for _, option := range options {
		option(invoke)
	}
	if dirs, err := settings.GetUserDirs(invoke.Name); err != nil {
		logging.Fatalf(invoke.ExitCodeError, "unable to get user directories: %v", err)
	} else {
		invoke.Configure(exec.WithUserDirs(dirs))
	}

	if app, err = initialize(invoke); err != nil {
		logging.Fatalf(invoke.ExitCodeError, "error initializing program: %v", err)
	}

	if cmd, err = app.Command(); cmd == nil || err != nil {
		logging.Fatalf(invoke.ExitCodeError, "command init: %v", app)
	}

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

		// https://github.com/carolynvs/stingoftheviper/blob/main/main.go
		var ppre = cmd.PersistentPreRun
		cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			// NOTE: cobra and viper can be bound in a few locations,
			// but PersistencePreRunE on the root command works well.
			// NOTE: app.Load() should have access to the cmd and args
			// from the invocation object it got when initialized.
			if err = app.Load(cmd, args); err != nil {
				logging.Fatalf(invoke.ExitCodeError, "Load() failed: %v", err)
			}

			if ppre != nil {
				ppre(cmd, args)
			}
		}

		if runerr = cmd.ExecuteContext(gctx); runerr != nil {
			runerr = errors.Wrap(runerr, "app.Run()")
		}
		return runerr
	})
	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		logging.Fatalf(invoke.ExitCodeError, "%v", err)
	}
}
