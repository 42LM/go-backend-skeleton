package httpmsg

import (
	"context"
	"log/slog"
	// "net/http"
)

// MsgSvc is an interface defined by the handler layer (the consumer).
// It specifies the methods the handler needs from the service.
type MsgSvc interface {
	FindMsg(ctx context.Context, id string) string
}

type MsgHandler struct {
	svc    MsgSvc
	logger *slog.Logger
}

// NewMsgHandler is a constructor that returns a *MsgHandler.
func NewMsgHandler(s MsgSvc, l *slog.Logger) *MsgHandler {
	return &MsgHandler{
		svc:    s,
		logger: l,
	}
}
