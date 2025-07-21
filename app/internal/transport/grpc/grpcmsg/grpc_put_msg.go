package grpcmsg

import (
	"context"

	pb "go-backend-skeleton/app/internal/transport/grpc/pb"
)

// PutMsg implements the PutMsg RPC method.
func (s *Server) PutMsg(ctx context.Context, in *pb.PutMsgRequest) (*pb.PutMsgReply, error) {
	err := s.MsgSvc.PutMsg(ctx, in.Id, in.Msg)
	if err != nil {
		return nil, err
	}

	return &pb.PutMsgReply{
		Id:  in.Id,
		Msg: in.Msg,
	}, nil
}
