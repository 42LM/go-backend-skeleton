// Package svcnone injects the none repository into the service layer.
package svcnone

import (
	"context"

	"go-backend-skeleton/app/internal/db/none"
)

// NoneRepo represents the none repository dependency that provides data for the NoneSvc.
type NoneRepo interface {
	// Find returns a static string.
	Find(ctx context.Context) string
}

// Prove that the none repositroy implements the NoneRepo interface
var _ NoneRepo = &none.NoneRepository{}

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
