package svcnone_test

import (
	"context"
	"testing"

	dbmock "go-backend-skeleton/app/internal/db/mock"
	"go-backend-skeleton/app/internal/svc/svcnone"

	// TODO: fix imports
	testifymock "github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func Test_FindNone(t *testing.T) {
	t.Parallel()

	mockNoneRepo := dbmock.NewMockNoneRepo(t)

	ctx := context.Background()

	testCases := map[string]struct {
		findNoneServiceCall *testifymock.Call
		setupMock           func(mockSvc *dbmock.MockNoneRepo)
	}{
		"ok": {
			setupMock: func(mockSvc *dbmock.MockNoneRepo) {
				mockSvc.On("Find", testifymock.Anything).Return("none").Once()
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.setupMock != nil {
				tc.setupMock(mockNoneRepo)
			}
			svc := svcnone.New(&svcnone.NoneSvcConfig{
				NoneRepo: mockNoneRepo,
			})
			res := svc.FindNone(ctx)
			assert.Equal(t, "none", res)
			assert.True(t, mockNoneRepo.AssertCalled(t, "Find", testifymock.Anything))
		})
	}
	assert.True(t, mockNoneRepo.AssertExpectations(t))
	assert.True(t, mockNoneRepo.AssertNumberOfCalls(t, "Find", 1))
}
