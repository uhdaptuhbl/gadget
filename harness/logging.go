package harness

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// FormatMessage is a utility for log messages prior to instantiating a logger.
func FormatMessage(level string, msg string) string {
	_, fn, lineNo, _ := runtime.Caller(2)
	return fmt.Sprintf(
		"%s\t%s\t%s:%d\t%s",
		UTCnowRFC3339(),
		level,
		filepath.Join(filepath.Base(filepath.Dir(fn)), filepath.Base(fn)),
		lineNo,
		msg,
	)
}

func LogDebug(msg string) {
	_, _ = fmt.Fprint(os.Stderr, (FormatMessage("DEBUG", msg) + "\n"))
}
func LogDebugf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, (FormatMessage("DEBUG", msg) + "\n"), args...)
}

func LogInfo(msg string) {
	_, _ = fmt.Fprint(os.Stderr, (FormatMessage("INFO", msg) + "\n"))
}
func LogInfof(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, (FormatMessage("INFO", msg) + "\n"), args...)
}

func LogFatal(exitCode int, msg string) {
	_, _ = fmt.Fprint(os.Stderr, (FormatMessage("ERROR", "FATAL: "+msg) + "\n"))
	os.Exit(exitCode)
}
func LogFatalf(exitCode int, msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, (FormatMessage("ERROR", "FATAL: "+msg) + "\n"), args...)
	os.Exit(exitCode)
}

func PrettyJSON(content []byte) (string, error) {
	var err error
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, content, "", "  "); err != nil {
		return prettyJSON.String(), err
	}
	return prettyJSON.String(), err
}
