package db

import "context"

type MsgRepository interface {
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) string
	Put(ctx context.Context, id, message string) error
}
