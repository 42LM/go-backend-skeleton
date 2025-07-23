package grpcmsg_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	svcmock "go-backend-skeleton/app/internal/svc/mock"
	"go-backend-skeleton/app/internal/transport/grpc/grpcmsg"
	pb "go-backend-skeleton/app/internal/transport/grpc/pb"
	transporthttp "go-backend-skeleton/app/internal/transport/http"

	testifymock "github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

func Test_GRPC_PutMsg(t *testing.T) {
	mockSvc := svcmock.NewMockMsgSvc(t)

	testCases := map[string]struct {
		setupMock func(mockSvc *svcmock.MockMsgSvc)
	}{
		"ok - no input struct set": {
			setupMock: func(mockSvc *svcmock.MockMsgSvc) {
				mockSvc.On("PutMsg", testifymock.Anything, "1", "test-msgX").Return(nil).Twice()
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.setupMock != nil {
				tc.setupMock(mockSvc)
			}

			// http
			handler := transporthttp.NewHandler(&transporthttp.HandlerConfig{
				MsgSvc: mockSvc,
				Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
			})
			httpSrv := httptest.NewServer(handler)
			defer httpSrv.Close()

			requestBodyString := `{"msg":"test-msg"}`
			bodyReader := strings.NewReader(requestBodyString)
			// ⚠️ context set in miggleware_logging.go
			// ctx := context.WithValue(context.Background(), "x", "X")

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/msg/1", httpSrv.URL), bodyReader)
			require.NoError(t, err)
			// req.Header.Set("Authorization", "Bearer the api secret")

			httpResp, err := httpSrv.Client().Do(req)
			require.NoError(t, err)
			defer httpResp.Body.Close()

			body, err := io.ReadAll(httpResp.Body)
			require.NoError(t, err)

			var resp pb.PutMsgReply
			err = json.Unmarshal(body, &resp)
			require.NoError(t, err)
			assert.Equal(t, "1", resp.Id)
			assert.Equal(t, "test-msg", resp.Msg)

			// grpc contract
			s := grpc.NewServer()
			srv := grpcmsg.New(mockSvc)
			pb.RegisterMessageServer(s, srv)
			conn := grpcClientConn(t, s)
			client := pb.NewMessageClient(conn)

			md := metadata.New(map[string]string{"x": "X"})
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			reply, err := client.PutMsg(ctx, &pb.PutMsgRequest{
				Id:  "1",
				Msg: "test-msg",
			})
			require.NoError(t, err)
			assert.Equal(t, "1", reply.Id)
			assert.Equal(t, "test-msg", reply.Msg)

			if tc.setupMock != nil {
				assert.True(t, mockSvc.AssertCalled(t, "PutMsg", testifymock.Anything, "1", "test-msgX"))
			}
		})
	}
	assert.True(t, mockSvc.AssertExpectations(t))
	assert.True(t, mockSvc.AssertNumberOfCalls(t, "PutMsg", 2))
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
