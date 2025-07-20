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

// NewHandler returns an HTTP handler with middleware wired in.
func NewHandler(config HandlerConfig) http.Handler {
	mux := muxify.NewMux()

	// TODO: Create logging mw
	// mux.Use(LoggingMiddleware)
	mux.Prefix("/v1")
	mux.HandleFunc("GET /none", httpnone.NewNoneHandler(config.NoneSvc, config.Logger).HandlerFunc)
	mux.HandleFunc("GET /msg/{id}", httpmsg.NewMsgHandler(config.MsgSvc, config.Logger).HandlerFunc)

	return mux
}
