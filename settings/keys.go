package settings

import (
	"strings"

	"gadget/logging"
)

const KeyConfigPath = "config"
const KeyEnvPrefix = "env-prefix"
const KeyProfileMode = "profile-mode"
const KeyVerbose = "verbose"
const KeyDebug = "debug"
const KeyForce = "force"

const DefaultVerbose = false
const DefaultDebug = false
const DefaultForce = false
const DefaultVerbosity = string(logging.LogVerbositySimple)

const HelpConfigPath = "Specify `<path>` to config file."
const HelpEnvPrefix = "Set a `<prefix>` for environment variables."
var HelpProfileMode = "Set the `<profile>` mode: {" + PrettyProfileModes() + "}."
const HelpVerbose = "Enable more detailed output."
const HelpDebug = "Enable debug output."
const HelpForce = "Perform potentially destructive actions."

const KeyLogFormat = "logformat"
const KeyLogLevel = "loglevel"
const KeyLogVerbosity = "verbosity"
const KeyLogOutput = "logoutput"

const DefaultLogFormat = string(logging.LogFormatJSON)
const DefaultLogLevel = string(logging.LogLevelInfo)
const DefaultLogVerbosity = string(logging.LogVerbositySimple)
var DefaultLogOutputs = []string{"stdout"}

const HelpLogFormat = "logging format"
const HelpLogLevel = "minimum logging level"
var HelpLogVerbosity = "Set `<verbosity>` of output: {" + logging.PrettyLogVerbosities() + "}."
const HelpLogOutput = "logging output file paths"

var DefaultPFlagsXform = map[string]string{
	KeyConfigPath:  "",
	KeyEnvPrefix:   "",
	KeyProfileMode: "profile_mode",
	KeyLogFormat:   "logging.format",
	KeyLogLevel:    "logging.level",
	KeyLogVerbosity:   "logging.verbosity",
	KeyLogOutput:   "logging.outputpaths",
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
