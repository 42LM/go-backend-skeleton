// Package loggingnone implements the decorator pattern
// and defines the decorator structs that wrap our concrete layer types.
package loggingnone

import (
	"context"

	"go-backend-skeleton/app/internal/logging"
	"go-backend-skeleton/app/internal/svc/svcnone"
	"go-backend-skeleton/app/internal/transport/http/httpnone"
)

// db level logging

type loggingRepo struct {
	next   svcnone.NoneRepo
	logger *logging.LoggerWrapper
}

func NewLoggingRepo(next svcnone.NoneRepo, logger *logging.LoggerWrapper) svcnone.NoneRepo {
	return &loggingRepo{next: next, logger: logger}
}

func (l *loggingRepo) Find(ctx context.Context) (res string) {
	l.logger.Log("Put", nil, map[string]any{"results.none": res}, nil)
	return l.next.Find(ctx)
}

// svc level logging

type loggingSvc struct {
	next   httpnone.NoneSvc
	logger *logging.LoggerWrapper
}

func NewLoggingSvc(next httpnone.NoneSvc, logger *logging.LoggerWrapper) httpnone.NoneSvc {
	return &loggingSvc{next: next, logger: logger}
}

func (l *loggingSvc) FindNone(ctx context.Context) (res string) {
	l.logger.Log("Put", nil, map[string]any{"results.none": res}, nil)
	return l.next.FindNone(ctx)
}
