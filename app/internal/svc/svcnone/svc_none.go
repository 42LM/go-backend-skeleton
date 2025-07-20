// Package svcnone injects the none repository into the service layer.
package svcnone

import (
	"context"
)

// NoneRepo represents the none repository dependency that provides data for the NoneSvc.
type NoneRepo interface {
	// Find returns a static string.
	Find(ctx context.Context) string
}

// NoneSvc is the concrete struct of the none service.
// It wraps the repository interface.
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
