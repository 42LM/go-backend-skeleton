// Package logging provides functionality to log.
package logging

import (
	"log/slog"
	"os"
	"time"
)

// LoggerWrapper is a wrapper for the logger of your choice.
type LoggerWrapper struct {
	logger *slog.Logger
}

// New returns a logger wrapper.
func New(l *slog.Logger) *LoggerWrapper {
	return &LoggerWrapper{
		logger: l,
	}
}

// Log wraps the logging functionality of the underlaying logger
// and will actually log things.
func (l *LoggerWrapper) Log(
	method string,
	params map[string]any,
	results map[string]any,
	err error,
) {
	defer func(begin time.Time) {
		l.logger.Info(
			method,
			"params", params,
			"results", results,
			"results.err", err,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
}

// newSlogger creates a customized *slog.Logger.
//
// * level omitted
// * message key renamed to method
// * change time format
func NewSlogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				return slog.Attr{}
			case slog.MessageKey:
				return slog.Attr{}
			case slog.TimeKey:
				t := a.Value.Time()
				newFormat := "2006-01-02 15:04:05"
				a.Value = slog.StringValue(t.Format(newFormat))
				return a
			default:
				return a
			}
		},
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	return logger
}
