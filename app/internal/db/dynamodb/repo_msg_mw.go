package dynamodb

import (
	"context"
	"log/slog"
	"time"

	"go-backend-skeleton/app/internal/db"
)

type MsgRepositoryMiddleware func(db.MsgRepository) db.MsgRepository

func NewMsgRepositoryLoggingMiddleware(logger *slog.Logger) MsgRepositoryMiddleware {
	return func(next db.MsgRepository) db.MsgRepository {
		return msgRepositoryLoggingMiddleware{
			logger: logger.With("repo", "Msg"),
			next:   next,
		}
	}
}

type msgRepositoryLoggingMiddleware struct {
	logger *slog.Logger
	next   db.MsgRepository
}

var _ db.MsgRepository = &msgRepositoryLoggingMiddleware{}

func (mw msgRepositoryLoggingMiddleware) Delete(
	ctx context.Context,
	id string,
) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"Delete",
			"params.id", id,
			"result.err", err,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return mw.next.Delete(ctx, id)
}

func (mw msgRepositoryLoggingMiddleware) Find(
	ctx context.Context,
	id string,
) (_ string) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"Find",
			"params.id", id,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return mw.next.Find(ctx, id)
}

func (mw msgRepositoryLoggingMiddleware) Put(
	ctx context.Context,
	id,
	msg string,
) (err error) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"Put",
			"params.id", id,
			"params.msg", msg,
			"result.err", err,
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return mw.next.Put(ctx, id, msg)
}
