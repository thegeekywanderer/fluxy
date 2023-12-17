// Package logger configures zerolog to be used with fluxy
package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

func init() {
	logger = log.With().Timestamp().CallerWithSkipFrameCount(3).Stack().Logger()
}

// SetLogLevel sets the logging level based on params
func SetLogLevel(level string) {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)

}

// Info log level
func Info(message string, args ...interface{}) {
	log.Info().Msgf(message, args...)
}

// Debug log level
func Debug(message string, args ...interface{}) {
	logger.Debug().Msgf(message, args...)
}

// Warn log level
func Warn(message string, args ...interface{}) {
	log.Warn().Msgf(message, args...)
}

// Error log level
func Error(message string, args ...interface{}) {
	logger.Error().Msgf(message, args...)
}

// Fatal log level
func Fatal(message string, args ...interface{}) {
	logger.Fatal().Msgf(message, args...)
	os.Exit(1)
}

// Log ...
func Log(message string, args ...interface{}) {
	if len(args) == 0 {
		log.Info().Msg(message)
	} else {
		log.Info().Msgf(message, args...)
	}
}
