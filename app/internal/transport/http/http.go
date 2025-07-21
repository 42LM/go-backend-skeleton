package http

import (
	"log/slog"
	"net/http"

	"go-backend-skeleton/app/internal/transport/grpc"
	"go-backend-skeleton/app/internal/transport/grpc/grpcmsg"
	"go-backend-skeleton/app/internal/transport/http/httpmsg"
	"go-backend-skeleton/app/internal/transport/http/httpnone"

	"github.com/42LM/muxify"
)

// HandlerConfig defines the config for the HTTP handler.
type HandlerConfig struct {
	NoneSvc    httpnone.NoneSvc
	MsgSvc     httpmsg.MsgSvc
	GRPCMsgSvc grpcmsg.MsgSvc
	Logger     *slog.Logger
}

// NewHandler creates a default *http.ServeMux and defines routes.
// It returns an HTTP handler with middleware wired in.
func NewHandler(config HandlerConfig) http.Handler {
	mux := muxify.NewMux()

	grpcServer := grpcmsg.New(config.GRPCMsgSvc)
	grpcMux := grpc.GRPCServeMux(grpcServer)

	// grpc gateway
	mux.Handle("/rpc/", grpcMux)

	// TODO: Create logging mw
	// mux.Use(LoggingMiddleware)
	mux.Prefix("/v1")
	mux.HandleFunc("GET /none", httpnone.New(config.NoneSvc).HandlerFunc)
	mux.HandleFunc("GET /msg/{id}", httpmsg.New(config.MsgSvc).HandlerFunc)

	return mux
}
