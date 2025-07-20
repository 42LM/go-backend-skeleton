// Package loggingnone implements the decorator pattern
// and defines the decorator structs that wrap our concrete layer types.
package loggingnone

import (
	"context"
	"log/slog"
	"time"

	"go-backend-skeleton/app/internal/svc/svcnone"
	"go-backend-skeleton/app/internal/transport/http/httpnone"
)

// db level logging

type loggingRepo struct {
	next   svcnone.NoneRepo
	logger *slog.Logger
}

func NewLoggingRepo(next svcnone.NoneRepo, logger *slog.Logger) svcnone.NoneRepo {
	return &loggingRepo{next: next, logger: logger}
}

func (l *loggingRepo) Find(ctx context.Context) string {
	defer func(begin time.Time) {
		l.logger.Info(
			"Find",
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return l.next.Find(ctx)
}

// svc level logging

type loggingSvc struct {
	next   httpnone.NoneSvc
	logger *slog.Logger
}

func NewLoggingSvc(next httpnone.NoneSvc, logger *slog.Logger) httpnone.NoneSvc {
	return &loggingSvc{next: next, logger: logger}
}

func (l *loggingSvc) FindNone(ctx context.Context) string {
	defer func(begin time.Time) {
		l.logger.Info(
			"FindNone",
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return l.next.FindNone(ctx)
}
