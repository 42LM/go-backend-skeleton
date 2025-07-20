package httpnone_test

import (
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	svcmock "go-backend-skeleton/app/internal/svc/mock"
	"go-backend-skeleton/app/internal/transport/http"

	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TODO: more love <3
func Test_NoneHandler_FindNone(t *testing.T) {
	t.Parallel()

	mockSvc := svcmock.NewMockNoneSvc(t)

	testCases := map[string]struct {
		findNoneServiceCall *testifymock.Call
		setupMock           func(mockSvc *svcmock.MockNoneSvc)

		httpMethod string

		expResponseCode int
		expResponse     string
	}{
		"ok": {
			setupMock: func(mockSvc *svcmock.MockNoneSvc) {
				mockSvc.On("FindNone", testifymock.Anything).Return("none").Once()
			},
			httpMethod:      nethttp.MethodGet,
			expResponseCode: nethttp.StatusOK,
			expResponse:     "none",
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

			handler := http.NewHandler(http.HandlerConfig{
				NoneSvc: mockSvc,
			})

			srv := httptest.NewServer(handler)
			defer srv.Close()

			cli := nethttp.DefaultClient
			req, err := nethttp.NewRequest(tc.httpMethod, srv.URL+"/none", nil)
			require.NoError(t, err)

			resp, err := cli.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expResponseCode, resp.StatusCode, "http status code")
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tc.expResponse, string(body))
			if tc.setupMock != nil {
				assert.True(t, mockSvc.AssertCalled(t, "FindNone", testifymock.Anything))
			}
		})
	}
	assert.True(t, mockSvc.AssertExpectations(t))
	assert.True(t, mockSvc.AssertNumberOfCalls(t, "FindNone", 1))
}
