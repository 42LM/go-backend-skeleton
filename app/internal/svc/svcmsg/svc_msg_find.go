package svcmsg

import (
	"context"
)

// FindMsg finds something in the msg repository by given id.
func (s *MsgSvc) FindMsg(ctx context.Context, id string) string {
	return s.msgRepo.Find(ctx, id)
}
