package config

import (
	"strings"

	"gadget/logging"
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

var DefaultPFlagsXform = map[string]string{
	KeyConfigPath:  "",
	KeyEnvPrefix:   "",
	KeyProfileMode: "profile_mode",
	KeyVerbosity:   "logging.verbosity",
	// KeyLogLevel:     "logging.level",
	// KeyLogFormat:    "logging.format",
	// KeyLogOutput:    "logging.outputpaths",
}

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
