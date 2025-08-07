package grpcmsg

import (
	"context"

	"go-backend-skeleton/app/internal/transport/grpc"
	pb "go-backend-skeleton/app/internal/transport/grpc/pb"
)

// PutMsg implements the PutMsg RPC method.
func (s *Server) PutMsg(ctx context.Context, in *pb.PutMsgRequest) (*pb.PutMsgReply, error) {
	md := grpc.MetaDataFromContext(ctx)

	err := s.MsgSvc.PutMsg(ctx, in.Id, in.Msg+md.X)
	if err != nil {
		return nil, grpc.ConvertError2Pb(err, "PutMsg")
	}

	return &pb.PutMsgReply{
		Id:  in.Id,
		Msg: in.Msg,
	}, nil
}
