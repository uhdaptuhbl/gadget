package exec

import (
	"time"

	"gadget/settings"
)

// Option applies changes to the provided invocation.
type Option func(invoke *Invocation)

func WithLogging() {}

func WithSignalHandler() {}

func HelpOnEmptyArgs(invoke *Invocation) {
	// TODO
	panic("NOT YET IMPLEMENTED!")
}

func NoParseFlags(invoke *Invocation) {
	// TODO
	panic("NOT YET IMPLEMENTED!")
}

func CreateMissingConfigFile(invoke *Invocation) {
	// TODO
	panic("NOT YET IMPLEMENTED!")
}

func WithName(name string) Option {
	return func(invoke *Invocation) {
		invoke.Name = name
	}
}

func WithArgs(args []string) Option {
	return func(invoke *Invocation) {
		invoke.Args = args
	}
}

func WithVersion(version string) Option {
	return func(invoke *Invocation) {
		invoke.Version = version
	}
}

func WithBuildId(buildId string) Option {
	return func(invoke *Invocation) {
		invoke.BuildId = buildId
	}
}

func WithBuildDate(buildDate string) Option {
	return func(invoke *Invocation) {
		invoke.BuildDate = buildDate
	}
}

func WithShutdownTimeout(shutdownTimeout time.Duration) Option {
	return func(invoke *Invocation) {
		invoke.ShutdownTimeout = shutdownTimeout
	}
}

func WithUserDirs(dirs settings.UserDirs) Option {
	return func(invoke *Invocation) {
		invoke.UserDirs = dirs
	}
}
