package grpcmsg

import (
	"go-backend-skeleton/app/internal/svc/svcmsg"
	"go-backend-skeleton/app/internal/transport"
	pb "go-backend-skeleton/app/internal/transport/grpc/pb"
)

// curl localhost:8081/ping
// curl -X POST localhost:8081/v1/greeter/say_hello -d '{"name":"luke"}'

// protoc -I ./proto \
//     --go_out=./proto --go_opt=paths=source_relative \
//     --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
//     --grpc-gateway_out=./proto --grpc-gateway_opt=paths=source_relative \
//     ./proto/helloworld.proto;

// server is the gRPC server implementation. It is an INTERNAL component
// that the HTTP gateway will call. It is not exposed to the public.
type Server struct {
	MsgSvc transport.MsgSvc

	pb.UnimplementedGreeterServer
}

// Prove that the message service implements the MsgSvc interface
var _ transport.MsgSvc = &svcmsg.MsgSvc{}

// New returns a grpc server.
func New(s transport.MsgSvc) *Server {
	return &Server{
		MsgSvc: s,
	}
}
