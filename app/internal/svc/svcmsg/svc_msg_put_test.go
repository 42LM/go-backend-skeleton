package svcmsg_test

import (
	"context"
	"testing"

	dbmock "go-backend-skeleton/app/internal/db/mock"
	"go-backend-skeleton/app/internal/svc/svcmsg"

	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_MsgSvc_PutMsg(t *testing.T) {
	t.Parallel()

	mockMsgRepo := dbmock.NewMockMsgRepo(t)

	ctx := context.Background()

	testCases := map[string]struct {
		findMsgServiceCall *testifymock.Call
		setupMock          func(mockSvc *dbmock.MockMsgRepo)
	}{
		"ok": {
			setupMock: func(mockSvc *dbmock.MockMsgRepo) {
				mockSvc.On("Put", testifymock.Anything, "777", "test-message").Return(nil).Once()
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
			err := svc.PutMsg(ctx, "777", "test-message")
			require.NoError(t, err)
			assert.True(t, mockMsgRepo.AssertCalled(t, "Put", testifymock.Anything, "777", "test-message"))
		})
	}
	assert.True(t, mockMsgRepo.AssertExpectations(t))
	assert.True(t, mockMsgRepo.AssertNumberOfCalls(t, "Put", 1))
}
