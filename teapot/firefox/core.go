package firefox

import (
	"path/filepath"

	"gadget/logging"
)

type Option func(core *firefoxCore)

func UseLogger(log logging.Logger) Option {
	return func(core *firefoxCore) {
		core.log = log
	}
}

func UseDataPath(path string) Option {
	return func(core *firefoxCore) {
		core.datapath = path
	}
}

func UseProfile(profile string) Option {
	return func(core *firefoxCore) {
		core.profile = profile
	}
}

type firefoxCore struct {
	log      logging.Logger
	datapath string
	profile  string
}

// TODO: should this be exported?
func newFirefoxCore(options ...Option) *firefoxCore {
	var core = new(firefoxCore)
	for _, option := range options {
		option(core)
	}
	if core.log == nil {
		core.log = logging.NewNoopLogger()
	}
	return core
}

func (core *firefoxCore) valid() error {
	var err error

	if err = core.load(); err != nil {
		return err
	}

	return err
}

func (core *firefoxCore) load() error {
	if core.datapath == "" {
		// TODO: load from OS and return proper error
		return new(EmptySqlitePathError)
	}
	if core.profile == "" {
		// TODO: load from file and return proper error
		return new(EmptySqlitePathError)
	}
	return nil
}

func (core *firefoxCore) cookieDBPath() string {
	return "file:" + filepath.Join(core.datapath, "Profiles", core.profile, "cookies.sqlite") + "?immutable=1"
}
