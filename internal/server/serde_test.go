package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/glass-cms/glasscms/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerializeJSONResponse(t *testing.T) {
	t.Parallel()

	type response struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name       string
		statusCode int
		data       response
		expected   string
	}{
		{
			name:       "valid response",
			statusCode: http.StatusOK,
			data:       response{Message: "success"},
			expected:   `{"message":"success"}`,
		},
		{
			name:       "internal server error",
			statusCode: http.StatusInternalServerError,
			data:       response{Message: "error"},
			expected:   `{"message":"error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			server.SerializeJSONResponse(rr, tt.statusCode, tt.data)

			assert.Equal(t, tt.statusCode, rr.Code)
			assert.JSONEq(t, tt.expected, rr.Body.String())
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		})
	}
}
func TestDeserializeJSONRequestBody(t *testing.T) {
	t.Parallel()

	type request struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name        string
		body        string
		expected    *request
		expectError bool
	}{
		{
			name:        "valid request body",
			body:        `{"message":"hello"}`,
			expected:    &request{Message: "hello"},
			expectError: false,
		},
		{
			name:        "invalid request body",
			body:        `{"message":}`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "empty request body",
			body:        ``,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))

			var result *request
			result, err := server.DeserializeJSONRequestBody[request](req)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
