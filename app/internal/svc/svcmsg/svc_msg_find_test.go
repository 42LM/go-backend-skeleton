package svcmsg_test

import (
	"context"
	"testing"

	dbmock "go-backend-skeleton/app/internal/db/mock"
	"go-backend-skeleton/app/internal/svc/svcmsg"

	// TODO: fix imports
	testifymock "github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func Test_FindMsg(t *testing.T) {
	t.Parallel()

	mockMsgRepo := dbmock.NewMockMsgRepo(t)

	ctx := context.Background()

	testCases := map[string]struct {
		findMsgServiceCall *testifymock.Call
		setupMock          func(mockSvc *dbmock.MockMsgRepo)
	}{
		"ok": {
			setupMock: func(mockSvc *dbmock.MockMsgRepo) {
				mockSvc.On("Find", testifymock.Anything, "777").Return("hello world").Once()
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.setupMock != nil {
				tc.setupMock(mockMsgRepo)
			}
			svc := svcmsg.New(&svcmsg.MsgSvcConfig{
				MsgRepo: mockMsgRepo,
			})
			res := svc.FindMsg(ctx, "777")
			assert.Equal(t, "hello world", res)
			assert.True(t, mockMsgRepo.AssertCalled(t, "Find", testifymock.Anything, "777"))
		})
	}
	assert.True(t, mockMsgRepo.AssertExpectations(t))
	assert.True(t, mockMsgRepo.AssertNumberOfCalls(t, "Find", 1))
}
