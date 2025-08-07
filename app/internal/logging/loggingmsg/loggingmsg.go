// Package loggingmsg implements the decorator pattern
// and defines the decorator structs that wrap our concrete layer types.
package loggingmsg

import (
	"context"
	"time"

	"go-backend-skeleton/app/internal/logging"
	"go-backend-skeleton/app/internal/svc/svcmsg"
	"go-backend-skeleton/app/internal/transport"
)

// TODO: generate logging

// db level logging

type loggingRepo struct {
	next   svcmsg.MsgRepo
	logger *logging.LoggerWrapper
}

func NewLoggingRepo(next svcmsg.MsgRepo, logger *logging.LoggerWrapper) svcmsg.MsgRepo {
	return &loggingRepo{next: next, logger: logger}
}

func (l *loggingRepo) Find(ctx context.Context, id string) (res string) {
	defer l.logger.Log("Find", map[string]any{"params.id": id}, map[string]any{"results.id": res}, nil)(time.Now())
	return l.next.Find(ctx, id)
}

func (l *loggingRepo) Put(ctx context.Context, id, msg string) (err error) {
	defer l.logger.Log("Put", map[string]any{"params.id": id, "params.msg": msg}, nil, err)(time.Now())
	return l.next.Put(ctx, id, msg)
}

// svc level logging

type loggingSvc struct {
	next   transport.MsgSvc
	logger *logging.LoggerWrapper
}

func NewLoggingSvc(next transport.MsgSvc, logger *logging.LoggerWrapper) transport.MsgSvc {
	return &loggingSvc{next: next, logger: logger}
}

func (l *loggingSvc) FindMsg(ctx context.Context, id string) (res string) {
	defer l.logger.Log("FindMsg", map[string]any{"params.id": id}, map[string]any{"results.msg": res}, nil)(time.Now())
	return l.next.FindMsg(ctx, id)
}

func (l *loggingSvc) PutMsg(ctx context.Context, id, msg string) (err error) {
	defer l.logger.Log("PutMsg", map[string]any{"params.id": id, "params.msg": msg}, nil, err)(time.Now())
	return l.next.PutMsg(ctx, id, msg)
}
