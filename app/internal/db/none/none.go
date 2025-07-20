package none

import (
	"context"
)

type NoneRepository struct{}

// NewNoneRepository constructs a none repository.
func NewNoneRepository() *NoneRepository {
	return &NoneRepository{}
}

func (r *NoneRepository) Find(ctx context.Context) string {
	return "NONE"
}

func (r *NoneRepository) NotUsed(ctx context.Context) string {
	return "#"
}
