package logging

// NewStdoutLogger - create a new stdout logger instance
func NewStdoutLogger(level string, format string) (Logger, error) {
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
