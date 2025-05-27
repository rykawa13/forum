package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func Setup(level string) {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)
}

func Info() *zerolog.Event {
	return log.Info()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}
