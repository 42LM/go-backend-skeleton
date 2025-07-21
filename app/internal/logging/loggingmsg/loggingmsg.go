// Package loggingmsg implements the decorator pattern
// and defines the decorator structs that wrap our concrete layer types.
package loggingmsg

import (
	"context"
	"log/slog"

	"go-backend-skeleton/app/internal/logging"
	"go-backend-skeleton/app/internal/svc/svcmsg"
	"go-backend-skeleton/app/internal/transport"
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
	logging.Log(l.logger, "Find", map[string]any{"params.id": id})
	return l.next.Find(ctx, id)
}

func (l *loggingRepo) Put(ctx context.Context, id, msg string) error {
	logging.Log(l.logger, "Put", map[string]any{"params.id": id, "params.msg": msg})
	return l.next.Put(ctx, id, msg)
}

// svc level logging

type loggingSvc struct {
	next   transport.MsgSvc
	logger *slog.Logger
}

func NewLoggingSvc(next transport.MsgSvc, logger *slog.Logger) transport.MsgSvc {
	return &loggingSvc{next: next, logger: logger}
}

func (l *loggingSvc) FindMsg(ctx context.Context, id string) string {
	logging.Log(l.logger, "FindMsg", map[string]any{"params.id": id})
	return l.next.FindMsg(ctx, id)
}

func (l *loggingSvc) PutMsg(ctx context.Context, id, msg string) error {
	logging.Log(l.logger, "PutMsg", map[string]any{"params.id": id, "params.msg": msg})
	return l.next.PutMsg(ctx, id, msg)
}
