package logging

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const pkgCheckTest = "_test.go"
const pkgCheckGorm = "gorm.io/gorm"
const pkgCheckGormAt = "gorm@"

// GormLoggerFromZap constructs a logger suitable for the gorm library.
func GormLogger(logger *zap.Logger, trace string, debug bool, encoding string) ZapGormLogger {
	return NewZapGormLogger(logger, ZapGormConfig{
		Debug:       debug,
		Encoding:    encoding,
		TracePrefix: trace,
	})
}

type ZapGormConfig struct {
	Debug                     bool
	Encoding                  string
	TracePrefix               string
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

type ZapGormLogger struct {
	ZapLogger *zap.Logger
	LogLevel  gormlogger.LogLevel
	cfg       ZapGormConfig
}

func NewZapGormLogger(zapLogger *zap.Logger, config ZapGormConfig) ZapGormLogger {
	var level = gormlogger.Warn
	if config.Debug {
		level = gormlogger.Info
	}
	if config.SlowThreshold <= 0 {
		config.SlowThreshold = 1000 * time.Millisecond
	}
	// TODO: should this set a default or respect an empty value?
	// if config.TracePrefix == "" {
	// 	config.TracePrefix = "query-trace"
	// }

	var newLogger = ZapGormLogger{
		ZapLogger: zapLogger,
		LogLevel:  level,
		cfg:       config,
	}

	return newLogger
}

func (l ZapGormLogger) SetAsDefault() {
	gormlogger.Default = l
}

func (l ZapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return ZapGormLogger{
		ZapLogger: l.ZapLogger,
		LogLevel:  level,
		cfg:       l.cfg,
	}
}

func (l ZapGormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.logger().Sugar().Debugf(str, args...)
}

func (l ZapGormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.logger().Sugar().Warnf(str, args...)
}

func (l ZapGormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.logger().Sugar().Errorf(str, args...)
}

func (l ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	var elapsed = time.Since(begin).Round(1 * time.Millisecond)
	var sql, rows = fc()

	switch {
	case l.checkErrorTrace(err):
		l.logger().Sugar().Errorw(
			l.cfg.TracePrefix,
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.checkElapsedTrace(elapsed):
		l.logger().Sugar().Warnw(
			l.cfg.TracePrefix,
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.LogLevel >= gormlogger.Info:
		l.logger().Sugar().Debugw(
			l.cfg.TracePrefix,
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}

func (l ZapGormLogger) checkErrorTrace(err error) bool {
	var errcheck = (err != nil && l.LogLevel >= gormlogger.Error)
	return (errcheck && (!l.cfg.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)))
}

func (l ZapGormLogger) checkElapsedTrace(elapsed time.Duration) bool {
	return (l.cfg.SlowThreshold != 0 && elapsed > l.cfg.SlowThreshold && l.LogLevel >= gormlogger.Warn)
}

func (l ZapGormLogger) logger() *zap.Logger {
	for index := 2; index < 15; index++ {
		_, file, _, ok := runtime.Caller(index)
		switch {
		case !ok:
		case strings.Contains(file, pkgCheckGormAt):
		case strings.Contains(file, pkgCheckGorm):
		case strings.HasSuffix(file, pkgCheckTest):
		default:
			// subtract one from index, otherwise it will also skip the expected file
			return l.ZapLogger.WithOptions(zap.AddCallerSkip(index - 1))
		}
	}
	return l.ZapLogger
}
