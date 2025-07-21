// Package httpmsg injects the message service into the handler layer.
package httpmsg

import (
	"context"
	"log/slog"

	"go-backend-skeleton/app/internal/svc/svcmsg"
)

// MsgSvc represents the message service dependency that provides data for the MsgHandler.
type MsgSvc interface {
	// FindMsg finds something in the msg repository by given id.
	FindMsg(ctx context.Context, id string) string
}

// Prove that the message service implements the MsgSvc interface
var _ MsgSvc = &svcmsg.MsgSvc{}

// MsgHandler is the concrete struct of the message handler.
// It wraps the service interface.
type MsgHandler struct {
	svc    MsgSvc
	logger *slog.Logger
}

// New returns a message handler.
func New(s MsgSvc, l *slog.Logger) *MsgHandler {
	return &MsgHandler{
		svc:    s,
		logger: l,
	}
}
