package svcnone

import (
	"context"
)

type NoneRepo interface {
	Find(ctx context.Context) string
}

type NoneSvc struct {
	noneRepo NoneRepo
}

// NoneSvcConfig contains the configuration params of the none service.
type NoneSvcConfig struct {
	NoneRepo NoneRepo
}

// New returns a none service.
func New(config *NoneSvcConfig) *NoneSvc {
	svc := &NoneSvc{
		noneRepo: config.NoneRepo,
	}
	return svc
}
