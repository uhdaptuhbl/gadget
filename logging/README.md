# Logging

The standard Logging interface for voltron packages.

## Implementations

### NOP
A no operation logger, typically used for unit tests or as a stub when a logger implementation has yet to be selected.

Usage:
```go
logger := logging.NewNopLogger()
```

### Zap
A logger implementation powered by the Uber zap library.

Usage:
```go
logger := logging.NewZapLogger()
```

## Custom Loggers
Situations will arise that require a custom logger. Voltron packages require acceptance of the logging.Logger interface.
This means that whenever there is a requirement for custom logging all one needs to do is implement the logging.Logger
interface.
