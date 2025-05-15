package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var l *zerolog.Logger

func GetInstance() *zerolog.Logger {
	return l
}

// InitLogger инициализирует Logger
func InitLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		CallerWithSkipFrameCount(2).
		Logger()
	l = &logger
	return l
}
