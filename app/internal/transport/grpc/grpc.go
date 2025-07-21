package grpc

import (
	"context"
	"log"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "go-backend-skeleton/app/internal/transport/grpc/pb"
)

// curl localhost:8081/ping
// curl -X POST localhost:8081/v1/greeter/say_hello -d '{"name":"luke"}'

// protoc -I ./proto \
//     --go_out=./proto --go_opt=paths=source_relative \
//     --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
//     --grpc-gateway_out=./proto --grpc-gateway_opt=paths=source_relative \
//     ./proto/helloworld.proto;

// GRPCServeMux creates a grpc client,
// creates a gateway multiplexer
// and returns the multiplexer.
func GRPCServeMux(srv pb.GreeterServer) *runtime.ServeMux {
	// create in memory listener
	// avoids openening a network port
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, srv)

	log.Println("Starting internal-only gRPC service on in-memory buffer")
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
	gwmux := runtime.NewServeMux()
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("failed to register gateway handler: %v", err)
	}

	return gwmux
}
