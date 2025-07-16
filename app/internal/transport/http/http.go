package http

import (
	"net/http"

	"go-backend-skeleton/app/internal/svc"
)

// HandlerConfig defines the config for the HTTP handler.
type HandlerConfig struct {
	Svc svc.Service
}

// NewHandler returns an HTTP handler with middleware wired in.
func NewHandler(config HandlerConfig) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/none", handlerFindNone(config.Svc))

	return mux
}
