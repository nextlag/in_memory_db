package errs

import (
	"errors"
)

// config.
var (
	ErrEmptyConfigPath   = errors.New("config path is empty")
	ErrIncorrectReader   = errors.New("reader is nil")
	ErrFailedReadBuffer  = errors.New("failed to read config content")
	ErrFailedParseConfig = errors.New("failed to parse config")
	ErrFailedOpenConfig  = errors.New("failed to open config file")
	ErrFailedInitConfig  = errors.New("failed to init config")
)

// app.
var (
	ErrServerIsNil        = errors.New("launcher is nil")
	ErrConfigIsNil        = errors.New("config is nil")
	ErrLoggerIsNil        = errors.New("logger is nil")
	ErrUnknownServerType  = errors.New("unknown server type")
	ErrFailedInitLauncher = errors.New("failed to init launcher")
	ErrFailedInitUseCase  = errors.New("failed to init use case")
	ErrComputeIsNil       = errors.New("compute is nil")
	ErrStorageIsNil       = errors.New("storage is nil")
)

// compute.
var (
	ErrInvalidQuery          = errors.New("invalid query")
	ErrInvalidCommand        = errors.New("invalid command")
	ErrInvalidQueryArguments = errors.New("invalid query arguments")
	ErrEngineTypeIsIncorrect = errors.New("engine type is incorrect")
	ErrFailedInitEngine      = errors.New("failed to init engine")
	ErrFailedInitCompute     = errors.New("failed to init compute")
	ErrParse                 = errors.New("parse")
	ErrAnalyze               = errors.New("analyze")
)

// storage.
var (
	ErrNotFound          = errors.New("not found")
	ErrEngineIsNil       = errors.New("engine is nil")
	ErrFailedInitStorage = errors.New("failed to init storage")
)

// IsNotFound checks if an error is the NotFound error or wraps it
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
