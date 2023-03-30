package settings

import (
	flag "github.com/spf13/pflag"
)

var Flags = new(FlagFactory)

type FlagFunc func(flags *flag.FlagSet)

type FlagFactory struct{}

func (factory FlagFactory) Build(name string, errHandling flag.ErrorHandling, opts ...FlagFunc) *flag.FlagSet {
	var flags = flag.NewFlagSet(name, errHandling)
	for _, flagfunc := range opts {
		flagfunc(flags)
	}
	return flags
}

func (factory FlagFactory) Configure(flags *flag.FlagSet, opts ...FlagFunc) *flag.FlagSet {
	for _, flagfunc := range opts {
		flagfunc(flags)
	}
	return flags
}

func (factory FlagFactory) Init(name string, errHandling flag.ErrorHandling) FlagFunc {
	return func(flags *flag.FlagSet) {
		flags.Init(name, errHandling)
	}
}

func (factory FlagFactory) Sort(sortFlags bool) FlagFunc {
	return func(flags *flag.FlagSet) {
		flags.SortFlags = sortFlags
	}
}

func (factory FlagFactory) IgnoreUnknown(ignore bool) FlagFunc {
	return func(flags *flag.FlagSet) {
		flags.ParseErrorsWhitelist = flag.ParseErrorsWhitelist{UnknownFlags: ignore}
	}
}

func (factory FlagFactory) Usage(usageFunc func(*flag.FlagSet)) FlagFunc {
	return func(flags *flag.FlagSet) {
		flags.Usage = func() { usageFunc(flags) }
	}
}

func (factory FlagFactory) BoolOption(key string, defaultVal bool, helpStr string) FlagFunc {
	return func(flags *flag.FlagSet) {
		flags.Bool(key, defaultVal, helpStr)
	}
}

func (factory FlagFactory) StringOption(key string, defaultVal string, helpStr string) FlagFunc {
	return func(flags *flag.FlagSet) {
		flags.String(key, defaultVal, helpStr)
	}
}
