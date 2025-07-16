package dynamodb

import (
	"context"
	"errors"
)

// handleError handles errors that can happen during a repo action.
func handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		return err
	}

	return databaseError{Err: err}
}

// databaseError is usually a programmer error. It should be reported as server error.
type databaseError struct {
	Err error
}

func (e databaseError) Error() string {
	return "database error: " + e.Err.Error()
}
