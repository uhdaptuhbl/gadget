package logging

import (
	"context"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type PgxLogger struct {
	l Logger
}

func PgxLoggerFromZap(l Logger) *PgxLogger {
	return &PgxLogger{l}
}

func (pl *PgxLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	var fields = make([]interface{}, len(data))
	var index = 0
	for k, v := range data {
		fields[index] = zap.Any(k, v)
		index++
	}

	switch level {
	case pgx.LogLevelTrace:
		pl.l.Debugw(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	case pgx.LogLevelDebug:
		pl.l.Debugw(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	case pgx.LogLevelInfo:
		pl.l.Debugw(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	case pgx.LogLevelWarn:
		pl.l.Warnw(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	case pgx.LogLevelError:
		pl.l.Errorw(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	default:
		pl.l.Errorw(msg, append(fields, zap.Stringer("PGX_LOG_LEVEL", level))...)
	}
}
