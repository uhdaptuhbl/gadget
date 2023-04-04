package logging

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetZapLogger(z *ZapLogger) *zap.Logger {
	return z.logger.Desugar().WithOptions(zap.AddCallerSkip(-1))
}

/*
ZapLogger is the Logger implementation with Uber's zap as its core.
*/
type ZapLogger struct {
	logger *zap.SugaredLogger
	cfg    zap.Config
	debug  bool
}

/*
NewZapLogger creates and configures a new ZapLogger instance.
*/
func NewZapLogger(config Config) (*ZapLogger, error) {
	var zl = ZapLogger{}
	var err = zl.Configure(config)
	if err != nil {
		return nil, err
	}

	return &zl, nil
}

/*
FromZapConfig builds a logger from provided zap.Config directly.
*/
func FromZapConfig(config zap.Config) (*ZapLogger, error) {
	var err error
	var zapLogger *zap.Logger
	var zl = ZapLogger{}

	if zapLogger, err = config.Build(); err != nil {
		return &zl, &InitializeError{err: err}
	}

	zl.logger = zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	zl.cfg = config

	return &zl, err
}

func (z *ZapLogger) clone(newLogger *zap.SugaredLogger) *ZapLogger {
	if newLogger == nil {
		newLogger = z.logger
	}
	return &ZapLogger{logger: newLogger, cfg: z.cfg, debug: z.debug}
}

func (z *ZapLogger) IsDebug() bool {
	return z.debug
}

func (z *ZapLogger) Config() zap.Config {
	return z.cfg
}

/*
Configure sets up or reconfigures a ZapLogger instance.
*/
func (z *ZapLogger) Configure(config Config) error {
	var err error
	var zapLogger *zap.Logger

	switch config.Format {
	case LogFormatHuman:
		z.cfg = zap.NewDevelopmentConfig()
		// z.cfg.EncoderConfig.ConsoleSeparator = " "
		// z.cfg.EncoderConfig.EncodeCaller = StandardCallerEncoder
		// z.cfg.EncoderConfig.EncodeTime = StandardTimeEncoder
		// z.cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case LogFormatJSON:
		z.cfg = zap.NewProductionConfig()
		// z.cfg.EncoderConfig.EncodeTime = RFC3339UTCTimeEncoder
		// z.cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	default:
		return &InvalidLogFormatError{Input: string(config.Format)}
	}

	switch config.Level {
	case LogLevelError:
		z.cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case LogLevelWarn:
		z.cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case LogLevelInfo:
		z.cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case LogLevelDebug:
		z.cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		z.debug = true
	default:
		return &InvalidLogLevelError{Input: string(config.Level)}
	}

	switch config.Verbosity {
	case LogVerbosityBare:
		z.cfg.DisableCaller = true
		z.cfg.DisableStacktrace = true
	case LogVerbositySimple:
		// z.cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		z.cfg.DisableCaller = false
		z.cfg.DisableStacktrace = true
	case LogVerbosityVerbose:
		// z.cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		z.cfg.DisableCaller = false
		z.cfg.DisableStacktrace = false
		// z.cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

		// NOTE: sadly there is no hook or encoder that can be used to override
		// and shorten the func key output which is the full path by default.
		z.cfg.EncoderConfig.FunctionKey = "func"

		z.cfg.InitialFields = make(map[string]interface{})
		z.cfg.InitialFields["pid"] = os.Getpid()
		z.cfg.InitialFields["runtime"] = runtime.Version()
		z.cfg.InitialFields["app"] = path.Base(os.Args[0])
		z.cfg.InitialFields["version"] = config.Version
		z.cfg.InitialFields["build"] = config.Build

		var name string
		if name, err = os.Hostname(); err == nil && name != "" {
			z.cfg.InitialFields["host"] = name
		}
	default:
		return &InvalidVerbosityError{Input: string(config.Verbosity)}
	}

	z.cfg.OutputPaths = config.OutputPaths
	z.cfg.EncoderConfig.EncodeTime = RFC3339UTCTimeEncoder

	if zapLogger, err = z.cfg.Build(); err != nil {
		return &InitializeError{err: err}
	}
	z.logger = zapLogger.WithOptions(zap.AddCallerSkip(1)).Sugar()

	return err
}

// StandardTimeEncoder serializes a time.Time
//
// If enc supports AppendTimeLayout(t time.Time,layout string), it's used
// instead of appending a pre-formatted string value.
func StandardTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, "2006-01-02 15:04:05.000000", enc)
}

func encodeTimeLayout(t time.Time, layout string, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

// StandardCallerEncoder s
func StandardCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// TODO: consider using a byte-oriented API to save an allocation.
	enc.AppendString(filepath.Base(caller.TrimmedPath()))
}

/*
HandleError checks if err has already been logged, otherwise logs it, wraps, and returns it.

This allows errors to be logged once, but errors which have not already
been logged can still be logged further up the call stack.
*/
func (z ZapLogger) HandleError(err error) error {
	var errcheck *LoggingHandledError

	if errors.As(err, &errcheck) {
		return err
	}
	z.logger.Error(err)

	return &LoggingHandledError{err: err}
}

func (z ZapLogger) Traced(ctx context.Context) Logger {
	var newLogger = z.logger
	if ctx != nil {
		if ctxRqID, ok := ctx.Value(RequestIDKey).(string); ok {
			newLogger = newLogger.With(zap.String("request_id", ctxRqID))
		}
	}

	return z.clone(newLogger)
}

func (z ZapLogger) WithExtraFields(fields map[string]string) Logger {
	var newLogger = z.logger
	for k, v := range fields {
		newLogger = newLogger.With(zap.String(k, v))
	}

	return z.clone(newLogger)
}

func (z ZapLogger) WithOptions(opts ...zap.Option) Logger {
	return z.clone(z.logger.Desugar().WithOptions(opts...).Sugar())
}

func (z ZapLogger) AddCallerSkip(skip int) Logger {
	return z.clone(z.logger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar())
}

func (z ZapLogger) Error(args ...interface{}) {
	z.logger.Error(args...)
}

func (z ZapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Errorf(template, args...)
}

func (z ZapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.logger.Errorw(msg, keysAndValues...)
}

func (z ZapLogger) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}

func (z ZapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Warnf(template, args...)
}

func (z ZapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	z.logger.Warnw(msg, keysAndValues...)
}

func (z ZapLogger) Info(args ...interface{}) {
	z.logger.Info(args...)
}

func (z ZapLogger) Infof(template string, args ...interface{}) {
	z.logger.Infof(template, args...)
}

func (z ZapLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, keysAndValues...)
}

func (z ZapLogger) Debug(args ...interface{}) {
	z.logger.Debug(args...)
}

func (z ZapLogger) Debugf(template string, args ...interface{}) {
	z.logger.Debugf(template, args...)
}

func (z ZapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	z.logger.Debugw(msg, keysAndValues...)
}
