package svc_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"go-backend-skeleton/app/internal/svc"

	dbmock "go-backend-skeleton/app/internal/db/mock"

	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func Test_FindNone(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mockNoneRepo := dbmock.NewMockNoneRepository(t)

	testCases := map[string]struct {
		noneRepoCall *testifymock.Call
		expRes       string
		expErr       error
	}{
		"ok": {
			noneRepoCall: mockNoneRepo.On(
				"Find", ctx,
			).Return(
				"none",
			).Once(),
			expRes: "none",
		},
	}
	for tname, tc := range testCases {
		t.Run(tname, func(t *testing.T) {
			s := svc.New(&svc.ServiceConfig{
				Logger:   slog.New(slog.NewTextHandler(io.Discard, nil)),
				NoneRepo: mockNoneRepo,
			})
			res := s.FindNone(ctx)
			assert.True(t, mockNoneRepo.AssertCalled(t, "Find", ctx))
			assert.Equal(t, tc.expRes, res)
		})
	}

	assert.True(t, mockNoneRepo.AssertExpectations(t))
	assert.True(t, mockNoneRepo.AssertNumberOfCalls(t, "Find", 1))
}
