package grpcmsg

import (
	"context"
	"errors"

	pb "go-backend-skeleton/app/internal/transport/grpc/pb"

	"google.golang.org/grpc/metadata"
)

// PutMsg implements the PutMsg RPC method.
func (s *Server) PutMsg(ctx context.Context, in *pb.PutMsgRequest) (*pb.PutMsgReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("")
	}
	getX := md.Get("x")
	var x string
	if getX != nil {
		x = getX[0]
	}

	err := s.MsgSvc.PutMsg(ctx, in.Id, in.Msg+x)
	if err != nil {
		return nil, err
	}

	return &pb.PutMsgReply{
		Id:  in.Id,
		Msg: in.Msg,
	}, nil
}
