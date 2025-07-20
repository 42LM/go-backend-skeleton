package svcnone

import (
	"context"
)

// FindNone finds something in the none repository.
func (s *NoneSvc) FindNone(ctx context.Context) string {
	return s.noneRepo.Find(ctx)
}
