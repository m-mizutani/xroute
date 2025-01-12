package logging

import (
	"log/slog"
	"sync"
)

var (
	defaultLogger = slog.Default()
	loggerMutex   sync.Mutex
)

func Default() *slog.Logger {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	return defaultLogger
}

func SetDefault(logger *slog.Logger) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	defaultLogger = logger
}
