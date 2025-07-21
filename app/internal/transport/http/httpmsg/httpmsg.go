// Package httpmsg injects the message service into the handler layer.
package httpmsg

import (
	"go-backend-skeleton/app/internal/svc/svcmsg"
	"go-backend-skeleton/app/internal/transport"
)

// Prove that the message service implements the MsgSvc interface
var _ transport.MsgSvc = &svcmsg.MsgSvc{}

// MsgHandler is the concrete struct of the message handler.
// It wraps the service interface.
type MsgHandler struct {
	svc transport.MsgSvc
}

// New returns a message handler.
func New(s transport.MsgSvc) *MsgHandler {
	return &MsgHandler{
		svc: s,
	}
}
