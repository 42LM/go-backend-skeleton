// Package logging provides functionality to log.
package logging

import (
	"log/slog"
	"os"
	"time"
)

// TODO: to be able to change the logger easily provide struct that wraps logger and method that wraps logging. Than params does not need to be changed later.

// NewSlogger creates a customized *slog.Logger.
func NewSlogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				return slog.Attr{}
			case slog.MessageKey:
				a.Key = "method"
				return a
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

func Log(l *slog.Logger, method string, params map[string]any) {
	defer func(begin time.Time) {
		l.Info(
			method,
			"params", params,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
}
