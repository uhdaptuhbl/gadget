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
func LogInfof(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, (FormatMessage("INFO", msg) + "\n"), args...)
}
func LogFatal(msg string) {
	_, _ = fmt.Fprint(os.Stderr, (FormatMessage("ERROR", "FATAL: "+msg) + "\n"))
	os.Exit(ExitCodeError)
}
func LogFatalf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, (FormatMessage("ERROR", "FATAL: "+msg) + "\n"), args...)
	os.Exit(ExitCodeError)
}

func PrettyJSON(content []byte) (string, error) {
	var err error
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, content, "", "  "); err != nil {
		return prettyJSON.String(), err
	}
	return prettyJSON.String(), err
}
