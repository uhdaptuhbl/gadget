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

// Application interface provides the API contract for applications.
type Application interface {
	Flags() *flag.FlagSet
	Viper() *viper.Viper
	Command() (*cobra.Command, error)
	Load() error
	OnConfigChange() ConfigChangeFunc
	ProfileMode() string
	HandleSignals(ctx context.Context, cancel context.CancelFunc, sigch chan os.Signal) error
	Run(ctx context.Context) error
}

// ApplicationBase can be embedded in your Application struct to avoid having to reimplement functions not used.
type ApplicationBase struct{}

func (app *ApplicationBase) Flags() *flag.FlagSet             { return nil }
func (app *ApplicationBase) Viper() *viper.Viper              { return nil }
func (app *ApplicationBase) Command() (*cobra.Command, error) { return nil, nil }
func (app *ApplicationBase) Load() error                      { return nil }
func (app *ApplicationBase) OnConfigChange() ConfigChangeFunc { return nil }
func (app *ApplicationBase) ProfileMode() string              { return "" }
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
