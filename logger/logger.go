package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func New(level zerolog.Level) zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level)
}
