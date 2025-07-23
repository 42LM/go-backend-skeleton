//go:generate sh proto/protoc_generate.sh

// Package grpc provides a function to create a grpc gateway.
package grpc

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	pb "go-backend-skeleton/app/internal/transport/grpc/pb"
)

// metadataAnnotator creates metadata for an outgoing request
// by extracting the fields from the context.
func metadataAnnotator(ctx context.Context, req *http.Request) metadata.MD {
	// return metadata.New(map[string]string{
	// 	"x": ctx.Value("x").(string),
	// })
	return metadata.New(map[string]string{
		"foo": "bar",
	})
}

// GRPCServeMux creates a grpc client,
// creates a gateway multiplexer
// and returns the multiplexer.
func GRPCServeMux(srv pb.MessageServer) *runtime.ServeMux {
	// create in memory listener
	// avoids openening a network port
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	s := grpc.NewServer()
	pb.RegisterMessageServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve internal gRPC: %v", err)
		}
	}()

	// create custom dialer that connects to the in-memory listener
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return lis.Dial()
	}

	// the gateway needs a client connection to the gRPC service
	// use custom dialer to connect to the in-memory buffer
	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		log.Fatalf("failed to dial internal gRPC server: %v", err)
	}

	// create gateway multiplexer
	// this will handle all requests for our gRPC service
	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(metadataAnnotator),
	)
	err = pb.RegisterMessageHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("failed to register gateway handler: %v", err)
	}

	return gwmux
}
