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

type loggingService struct {
	next   httpnone.NoneSvc
	logger *slog.Logger
}

func NewLoggingService(next httpnone.NoneSvc, logger *slog.Logger) httpnone.NoneSvc {
	return &loggingService{next: next, logger: logger}
}

func (l *loggingService) FindNone(ctx context.Context) string {
	defer func(begin time.Time) {
		l.logger.Info(
			"FindNone",
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return l.next.FindNone(ctx)
}
