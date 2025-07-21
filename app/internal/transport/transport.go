package transport

import (
	"context"

	"go-backend-skeleton/app/internal/svc/svcmsg"
)

// MsgSvc represents the message service dependency that provides data for the MsgHandler.
type MsgSvc interface {
	// FindMsg finds something in the msg repository by given id.
	FindMsg(ctx context.Context, id string) string
	// PutMsg creates or updates something in the msg repository.
	PutMsg(ctx context.Context, id, msg string) error
}

// Prove that the message service implements the MsgSvc interface
var _ MsgSvc = &svcmsg.MsgSvc{}
