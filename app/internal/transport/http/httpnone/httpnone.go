package httpnone

import (
	"context"
	"log/slog"
	// "net/http"
)

// NoneSvc is an interface defined by the handler layer (the consumer).
// It specifies the methods the handler needs from the service.
type NoneSvc interface {
	FindNone(ctx context.Context) string
}

type NoneHandler struct {
	svc    NoneSvc
	logger *slog.Logger
}

// NewNoneHandler is a constructor that returns a *NoneHandler.
func NewNoneHandler(s NoneSvc, l *slog.Logger) *NoneHandler {
	return &NoneHandler{
		svc:    s,
		logger: l,
	}
}
