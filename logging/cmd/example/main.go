package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"

	"gadget/logging"
)

func main() {
	l, err := logging.NewZapLogger(logging.Config{
		Build:       "1",
		Version:     "1",
		Format:      "json",
		Level:       "debug",
		OutputPaths: []string{"stdout"},
		Verbosity:   "simple",
	})

	if err != nil {
		fmt.Printf("unable to create logger: %v", err)
		os.Exit(1)
	}

	withExtraFields := l.WithExtraFields(map[string]string{"id": "1"})
	withExtraFields.Infow("testing", "field1", "1")

	// Gorm logger
	gl := logging.GormLogger(logging.GetZapLogger(l), "", false, "")
	gl.Info(context.Background(), "blah")

	// Pgx logger
	pl := logging.PgxLoggerFromZap(withExtraFields)
	pl.Log(context.Background(), pgx.LogLevelError, "test", map[string]interface{}{"key": "value"})
}
