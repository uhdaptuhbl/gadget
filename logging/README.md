# Logging

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
