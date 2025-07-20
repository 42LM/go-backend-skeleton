// Package none provides the most basic example of a database repository.
// It returns fake data and serves as example for seeing the path of structure of the onion.
package none

import (
	"context"
)

type NoneRepository struct{}

// NewNoneRepository constructs a none repository with middleware wired in.
func NewNoneRepository() *NoneRepository {
	return &NoneRepository{}
}

// Find returns a static string.
func (r *NoneRepository) Find(ctx context.Context) string {
	return "NONE"
}
