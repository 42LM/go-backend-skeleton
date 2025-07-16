package none

import (
	"context"

	"go-backend-skeleton/app/internal/db"
)

type NoneRepository struct{}

// NewNoneRepository constructs a none repository with middleware wired in.
func NewNoneRepository() db.NoneRepository {
	var repo db.NoneRepository = &NoneRepository{}
	return repo
}

// Check if struct implements interface explicitly.
var _ db.NoneRepository = &NoneRepository{}

func (r *NoneRepository) Find(ctx context.Context) string {
	return "NONE"
}
