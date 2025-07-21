// Package grpcmsg provides a grpc server / client.
// It exposes HTTP via grpc gateway.
package grpcmsg

import (
	"go-backend-skeleton/app/internal/transport"
	pb "go-backend-skeleton/app/internal/transport/grpc/pb"
)

// Server is the gRPC server implementation. It is an INTERNAL component
// that the HTTP gateway will call. It is not exposed to the public.
type Server struct {
	MsgSvc transport.MsgSvc

	pb.UnimplementedMessageServer
}

// New returns a grpc server.
func New(s transport.MsgSvc) *Server {
	return &Server{
		MsgSvc: s,
	}
}
