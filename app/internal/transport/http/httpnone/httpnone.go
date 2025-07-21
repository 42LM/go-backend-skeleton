// Package httpnone injects the none service into the handler layer.
package httpnone

import (
	"context"

	"go-backend-skeleton/app/internal/svc/svcnone"
)

// NoneSvc represents the none service dependency that provides data for the NoneHandler.
type NoneSvc interface {
	FindNone(ctx context.Context) string
}

// Prove that the none service implements the NoneSvc interface
var _ NoneSvc = &svcnone.NoneSvc{}

// NoneHandler is the concrete struct of the none handler.
// It wraps the service interface.
type NoneHandler struct {
	svc NoneSvc
}

// New returns a none handler.
func New(s NoneSvc) *NoneHandler {
	return &NoneHandler{
		svc: s,
	}
}
