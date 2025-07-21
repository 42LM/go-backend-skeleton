package logging

import (
	"log/slog"
	"time"
)

func Log(l *slog.Logger, method string, params map[string]any) {
	defer func(begin time.Time) {
		l.Info(
			"Find",
			"params", params,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
}
