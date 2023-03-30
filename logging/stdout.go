package logging

// NewStdoutLogger - create a new stdout logger instance
func NewStdoutLogger(level LogLevel, format LogFormat) (Logger, error) {
	zl := ZapLogger{}
	err := zl.Configure(Config{
		Format:      format,
		Level:       level,
		OutputPaths: []string{"stdout"},
	})

	if err != nil {
		return nil, err
	}

	return &zl, nil
}
