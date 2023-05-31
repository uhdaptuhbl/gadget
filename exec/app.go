package exec

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"gadget/halt"
)

const ExitCodeError = 42

type ConfigChangeFunc func() error

// Program interface provides the API contract for applications.
type Program interface {
	Flags() *flag.FlagSet
	Viper() *viper.Viper
	Load() error
	OnConfigChange() ConfigChangeFunc
	ProfileMode() string
	HandleSignals(ctx context.Context, cancel context.CancelFunc, sigch chan os.Signal) error
	Run(ctx context.Context) error
}

// ApplicationBase can be embedded in your Program struct to avoid having to reimplement functions not used.
type ProgramBase struct{}

func (p *ProgramBase) Flags() *flag.FlagSet             { return nil }
func (p *ProgramBase) Viper() *viper.Viper              { return nil }
func (p *ProgramBase) Load() error                      { return nil }
func (p *ProgramBase) OnConfigChange() ConfigChangeFunc { return nil }
func (p *ProgramBase) ProfileMode() string              { return "" }
func (p *ProgramBase) HandleSignals(ctx context.Context, cancel context.CancelFunc, sigch chan os.Signal) error {
	signal.Notify(sigch, syscall.SIGHUP, syscall.SIGTERM, os.Interrupt)

	return halt.HandleInterrupts(halt.InterruptOptions{
		Signalch:        sigch,
		Context:         ctx,
		Shutdown:        cancel,
		GracefulTimeout: time.Second * 5,
		TimeoutExitCode: ExitCodeError,
		// TODO: is there a good way to pass this a logger instance?
		// Log:             app.Logger(),
	})
}

// Application interface provides the API contract for applications.
type Application interface {
	// Command should return the root cobra command of the application
	Command() (*cobra.Command, error)

	// Load should load the application configuration.
	// It should have access to the cmd and args via the invocation
	Load(cmd *cobra.Command, args []string) error

	// ProfileMode returns the string representation of the profiling mode.
	ProfileMode() string
	HandleSignals(ctx context.Context, cancel context.CancelFunc, sigch chan os.Signal) error
}

// ApplicationBase can be embedded in your Application struct to avoid having to reimplement functions not used.
type ApplicationBase struct{
	// Root   *cobra.Command
	// Flags  *flag.FlagSet
	// Snek   *viper.Viper
	// Invoke *Invocation
	// Log    logging.Logger
}

func (app *ApplicationBase) Command() (*cobra.Command, error)             { return nil, nil }
func (app *ApplicationBase) Load(cmd *cobra.Command, args []string) error { return nil }
func (app *ApplicationBase) ProfileMode() string                          { return "" }
func (app *ApplicationBase) HandleSignals(ctx context.Context, cancel context.CancelFunc, sigch chan os.Signal) error {
	signal.Notify(sigch, syscall.SIGHUP, syscall.SIGTERM, os.Interrupt)

	return halt.HandleInterrupts(halt.InterruptOptions{
		Signalch:        sigch,
		Context:         ctx,
		Shutdown:        cancel,
		GracefulTimeout: time.Second * 5,
		TimeoutExitCode: ExitCodeError,
		// TODO: is there a good way to pass this a logger instance?
		// Log:             app.Logger(),
	})
}
