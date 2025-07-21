// Package loggingmsg implements the decorator pattern
// and defines the decorator structs that wrap our concrete layer types.
package loggingmsg

import (
	"context"
	"log/slog"
	"time"

	"go-backend-skeleton/app/internal/svc/svcmsg"
	"go-backend-skeleton/app/internal/transport/grpc/grpcmsg"
	"go-backend-skeleton/app/internal/transport/http/httpmsg"
)

// TODO: generate logging

// db level logging

type loggingRepo struct {
	next   svcmsg.MsgRepo
	logger *slog.Logger
}

func NewLoggingRepo(next svcmsg.MsgRepo, logger *slog.Logger) svcmsg.MsgRepo {
	return &loggingRepo{next: next, logger: logger}
}

func (l *loggingRepo) Find(ctx context.Context, id string) string {
	defer func(begin time.Time) {
		l.logger.Info(
			"Find",
			"params.id", id,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return l.next.Find(ctx, id)
}

func (l *loggingRepo) Put(ctx context.Context, id, msg string) error {
	defer func(begin time.Time) {
		l.logger.Info(
			"Put",
			"params.id", id,
			"params.msg", msg,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return l.next.Put(ctx, id, msg)
}

// svc level logging

type loggingSvc struct {
	next   httpmsg.MsgSvc
	logger *slog.Logger
}

func NewLoggingSvc(next httpmsg.MsgSvc, logger *slog.Logger) httpmsg.MsgSvc {
	return &loggingSvc{next: next, logger: logger}
}

func (l *loggingSvc) FindMsg(ctx context.Context, id string) string {
	defer func(begin time.Time) {
		l.logger.Info(
			"FindMsg",
			"params.id", id,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return l.next.FindMsg(ctx, id)
}

type loggingGrpcSvc struct {
	next   grpcmsg.MsgSvc
	logger *slog.Logger
}

func NewLoggingGrpcSvc(next grpcmsg.MsgSvc, logger *slog.Logger) grpcmsg.MsgSvc {
	return &loggingGrpcSvc{next: next, logger: logger}
}

func (l *loggingGrpcSvc) PutMsg(ctx context.Context, id, msg string) error {
	defer func(begin time.Time) {
		l.logger.Info(
			"PutMsg",
			"params.id", id,
			"params.msg", msg,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return l.next.PutMsg(ctx, id, msg)
}
