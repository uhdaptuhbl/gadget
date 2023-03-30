package halt

import (
	"context"
	"os"
	"syscall"
	"time"

	"bitbucket.services.cymru.com/voltron/logging"
)

type HandlerFunc func(signal os.Signal) error

type InterruptOptions struct {
	Signalch        <-chan os.Signal
	Context         context.Context
	Shutdown        context.CancelFunc
	GracefulTimeout time.Duration
	TimeoutExitCode int
	Log             logging.Logger
	Callback        HandlerFunc
}

func HandleInterrupts(options InterruptOptions) error {
	var err error
	var sig os.Signal

	var log = options.Log

	// TODO: add interrupt -> message map to options allowing caller to
	// override the log messages used when handling interrupts.
	var msgSighup = "received syscall.SIGHUP - shutting down (pid=%d)."
	var msgSigterm = "received syscall.SIGTERM - shutting down (pid=%d)."
	var msgSigint = "received os.Interrupt - shutting down (pid=%d)."

	// The defer ensures that the Graceful shutdown deadline will be started,
	// even if the context is canceled for another reason than a signal
	// and running the timeout in a goroutine allows the handler to exit.
	defer func() {
		if options.GracefulTimeout <= (0 * time.Second) {
			return
		}

		// TODO: should this use time.Sleep(), time.After()/time.AfterFunc(),
		// or context.WithTimeout()?
		// time.AfterFunc(options.GracefulTimeout, func() {
		// 	log.Errorf("Graceful shutdown timeout limit of %.2f reached - now exiting", options.GracefulTimeout)
		// 	os.Exit(options.TimeoutExitCode)
		// })
		go func() {
			time.Sleep(options.GracefulTimeout)
			if log != nil {
				log.Errorf("Graceful shutdown timeout limit of %.2f reached - now exiting", options.GracefulTimeout)
			}
			os.Exit(options.TimeoutExitCode)
		}()
	}()

	for {
		select {
		case <-options.Context.Done():
			return options.Context.Err()
		case sig = <-options.Signalch:
			if sig == syscall.SIGHUP {
				if log != nil {
					log.Infof(msgSighup, os.Getpid())
				}
				options.Shutdown()
				return err
			} else if sig == syscall.SIGTERM {
				if log != nil {
					log.Infof(msgSigterm, os.Getpid())
				}
				options.Shutdown()
				return err
			} else if sig == os.Interrupt {
				if log != nil {
					log.Infof(msgSigint, os.Getpid())
				}
				options.Shutdown()
				return err
			}
		}
	}
}
