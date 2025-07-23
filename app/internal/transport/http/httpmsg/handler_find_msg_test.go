package httpmsg_test

import (
	"io"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	svcmock "go-backend-skeleton/app/internal/svc/mock"
	"go-backend-skeleton/app/internal/transport/http"

	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_MsgHandler_FindMsg(t *testing.T) {
	t.Parallel()

	mockSvc := svcmock.NewMockMsgSvc(t)

	testCases := map[string]struct {
		findMsgServiceCall *testifymock.Call
		setupMock          func(mockSvc *svcmock.MockMsgSvc)

		httpMethod string

		expResponseCode int
		expResponse     string
	}{
		"ok": {
			setupMock: func(mockSvc *svcmock.MockMsgSvc) {
				mockSvc.On("FindMsg", testifymock.Anything, "777").Return("a wonderful msg").Once()
			},
			httpMethod:      nethttp.MethodGet,
			expResponseCode: nethttp.StatusOK,
			expResponse:     "a wonderful msg",
		},
		"internal server error: wrong request method": {
			httpMethod:      nethttp.MethodPatch,
			expResponseCode: nethttp.StatusMethodNotAllowed,
			expResponse:     "Method Not Allowed\n",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.setupMock != nil {
				tc.setupMock(mockSvc)
			}

			nopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := http.NewHandler(&http.HandlerConfig{
				MsgSvc: mockSvc,
				Logger: nopLogger,
			})

			srv := httptest.NewServer(handler)
			defer srv.Close()

			cli := nethttp.DefaultClient
			req, err := nethttp.NewRequest(tc.httpMethod, srv.URL+"/v1/msg/777", nil)
			require.NoError(t, err)

			resp, err := cli.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expResponseCode, resp.StatusCode, "http status code")
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tc.expResponse, string(body))
			if tc.setupMock != nil {
				assert.True(t, mockSvc.AssertCalled(t, "FindMsg", testifymock.Anything, "777"))
			}
		})
	}
	assert.True(t, mockSvc.AssertExpectations(t))
	assert.True(t, mockSvc.AssertNumberOfCalls(t, "FindMsg", 1))
}
