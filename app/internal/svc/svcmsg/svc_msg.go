package svcmsg

import (
	"context"
)

type MsgRepo interface {
	Find(
		ctx context.Context,
		id string,
	) string
}

type MsgSvc struct {
	msgRepo MsgRepo
}

// MsgSvcConfig contains the configuration params of the msg service.
type MsgSvcConfig struct {
	MsgRepo MsgRepo
}

// New returns a msg service.
func New(config *MsgSvcConfig) *MsgSvc {
	svc := &MsgSvc{
		msgRepo: config.MsgRepo,
	}
	return svc
}
