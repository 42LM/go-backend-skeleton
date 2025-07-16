package svc

import (
	"context"
	"time"
)

// FindNone finds something in the none repository.
func (s *service) FindNone(ctx context.Context) string {
	return s.noneRepo.Find(ctx)
}

func (mw loggingMiddleware) FindNone(ctx context.Context) string {
	defer func(begin time.Time) {
		mw.logger.Info(
			"FindNone",
			"took", float64(time.Since(begin))/1e6,
		)
	}(time.Now())
	return mw.next.FindNone(ctx)
}
