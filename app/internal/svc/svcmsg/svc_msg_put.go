package svcmsg

import (
	"context"
)

// PutMsg creates or updates a message.
func (s *MsgSvc) PutMsg(ctx context.Context, id, msg string) error {
	return s.msgRepo.Put(ctx, id, msg)
}
