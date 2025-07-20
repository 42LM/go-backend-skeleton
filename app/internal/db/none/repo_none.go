package none

import (
	"context"
)

type NoneRepository struct{}

// NewNoneRepository constructs a none repository with middleware wired in.
func NewNoneRepository() *NoneRepository {
	return &NoneRepository{}
}

func (r *NoneRepository) Find(ctx context.Context) string {
	return "NONE"
}
