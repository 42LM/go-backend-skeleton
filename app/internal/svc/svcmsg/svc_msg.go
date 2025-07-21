// Package svcmsg injects the message repository into the service layer.
package svcmsg

import (
	"context"

	"go-backend-skeleton/app/internal/db/dynamodb"
)

// MsgRepo represents the message repository dependency that provides data for the MsgSvc.
// This interface allows for decoupling the service from a concrete data source,
// making it easy to swap implementations for testing or other purposes.
type MsgRepo interface {
	// Find loads the message and returns it.
	Find(ctx context.Context, id string) string
}

// Prove that the message repositroy implements the MsgRepo interface
var _ MsgRepo = &dynamodb.MsgRepository{}

// MsgSvc is the concrete struct of the message service.
// It wraps the repository interface.
type MsgSvc struct {
	msgRepo MsgRepo
}

// MsgSvcConfig contains the configuration params of the message service.
type MsgSvcConfig struct {
	MsgRepo MsgRepo
}

// New returns a message service.
func New(config *MsgSvcConfig) *MsgSvc {
	svc := &MsgSvc{
		msgRepo: config.MsgRepo,
	}
	return svc
}
