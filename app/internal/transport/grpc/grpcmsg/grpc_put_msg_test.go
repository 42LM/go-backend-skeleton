package grpcmsg_test

import (
	"context"
	"net"
	"testing"

	svcmock "go-backend-skeleton/app/internal/svc/mock"
	"go-backend-skeleton/app/internal/transport/grpc/grpcmsg"
	"go-backend-skeleton/app/internal/transport/grpc/pb"

	testifymock "github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func Test_GRPC_PutMsg(t *testing.T) {
	mockSvc := svcmock.NewMockMsgSvc(t)

	testCases := map[string]struct {
		setupMock func(mockSvc *svcmock.MockMsgSvc)
	}{
		"ok - no input struct set": {
			setupMock: func(mockSvc *svcmock.MockMsgSvc) {
				mockSvc.On("PutMsg", testifymock.Anything, "1", "test-msg").Return(nil).Once()
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.setupMock != nil {
				tc.setupMock(mockSvc)
			}

			s := grpc.NewServer()
			srv := grpcmsg.New(mockSvc)
			pb.RegisterMessageServer(s, srv)
			conn := grpcClientConn(t, s)
			client := pb.NewMessageClient(conn)

			reply, err := client.PutMsg(context.Background(), &pb.PutMsgRequest{
				Id:  "1",
				Msg: "test-msg",
			})
			require.NoError(t, err)
			assert.Equal(t, "1", reply.Id)
			assert.Equal(t, "test-msg", reply.Msg)

			if tc.setupMock != nil {
				assert.True(t, mockSvc.AssertCalled(t, "PutMsg", testifymock.Anything, "1", "test-msg"))
			}
		})
	}
	assert.True(t, mockSvc.AssertExpectations(t))
	assert.True(t, mockSvc.AssertNumberOfCalls(t, "PutMsg", 1))
}

// grpcClientConn constructs an in-memory gRPC connection which has all networking involved.
func grpcClientConn(t *testing.T, grpcServer *grpc.Server) *grpc.ClientConn {
	t.Helper()

	const bufSize = 1024 * 1024
	grpcListener := bufconn.Listen(bufSize)
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return grpcListener.Dial()
	}

	go func() {
		err := grpcServer.Serve(grpcListener)
		require.NoError(t, err)
	}()

	dialOptions := []grpc.DialOption{
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("passthrough:///bufnet", dialOptions...)
	require.NoError(t, err, "failed to dial bufnet")

	t.Cleanup(func() {
		conn.Close()
	})

	return conn
}
