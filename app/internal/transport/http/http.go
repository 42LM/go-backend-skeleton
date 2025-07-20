package http

import (
	"log/slog"
	"net/http"

	"go-backend-skeleton/app/internal/transport/http/httpmsg"
	"go-backend-skeleton/app/internal/transport/http/httpnone"

	"github.com/42LM/muxify"
)

// HandlerConfig defines the config for the HTTP handler.
type HandlerConfig struct {
	NoneSvc httpnone.NoneSvc
	MsgSvc  httpmsg.MsgSvc
	Logger  *slog.Logger
}

// NewHandler creates a default *http.ServeMux and defines routes.
// It returns an HTTP handler with middleware wired in.
func NewHandler(config HandlerConfig) http.Handler {
	mux := muxify.NewMux()

	// TODO: Create logging mw
	// mux.Use(LoggingMiddleware)
	mux.Prefix("/v1")
	mux.HandleFunc("GET /none", httpnone.New(config.NoneSvc, config.Logger).HandlerFunc)
	mux.HandleFunc("GET /msg/{id}", httpmsg.New(config.MsgSvc, config.Logger).HandlerFunc)

	return mux
}
