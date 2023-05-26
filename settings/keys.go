package settings

import (
	"strings"

	"gadget/logging"
)

const KeyConfigPath = "config"
const HelpConfigPath = "Specify `<path>` to config file."

const KeyEnvPrefix = "env-prefix"
const HelpEnvPrefix = "Set a `<prefix>` for environment variables."

var KeyProfileMode = "profile-mode"
var HelpProfileMode = "Set the `<profile>` mode: {" + PrettyProfileModes() + "}."

const KeyVerbose = "verbose"
const DefaultVerbose = false
const HelpVerbose = "Enable more detailed output."

const KeyDebug = "debug"
const DefaultDebug = false
const HelpDebug = "Enable debug output."

const KeyForce = "force"
const DefaultForce = false
const HelpForce = "Allow potentially destructive actions."

var KeyLogLevel = "log-level"
var DefaultLogLevel = string(logging.LogLevelDebug)
var HelpLogLevel = "Set logging `<level>`: {" + logging.PrettyLogLevels() + "}."

var KeyLogFormat = "log-format"
var DefaultLogFormat = string(logging.LogFormatJSON)
var HelpLogFormat = "Set logging `<format>`: {" + logging.PrettyLogFormats() + "}."

var KeyLogVerbosity = "log-verbosity"
var DefaultLogVerbosity = string(logging.LogVerbosityBare)
var HelpLogVerbosity = "Set logging `<verbosity>`: {" + logging.PrettyLogVerbosities() + "}."

var KeyLogOutputs = "log-outputs"
var DefaultLogOutputs = []string{"stdout"}
var HelpLogOutputs = "Set logging file paths to write to: {stdout} (typically a local absolute file path, but when using the zap logging package, there are some additional options; see https://pkg.go.dev/go.uber.org/zap#Open)."

var DefaultPFlagsXform = map[string]string{
	KeyConfigPath:   "",
	KeyEnvPrefix:    "",
	KeyProfileMode:  "profile_mode",
	KeyLogVerbosity: "logging.verbosity",
	KeyLogFormat:    "logging.format",
	KeyLogLevel:     "logging.level",
	KeyLogOutputs:   "logging.outputpaths",
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
