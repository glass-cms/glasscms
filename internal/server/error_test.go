package server_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/glass-cms/glasscms/internal/server"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/resource"
	"github.com/stretchr/testify/require"
)

func TestErrorMapperAlreadyExistsError(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}
	tests := map[string]struct {
		args         args
		want         *api.Error
		expectPanics bool
	}{
		"maps resource.AlreadyExistsError to an API error response": {
			args: args{
				err: resource.NewAlreadyExistsError("item1", "item", errors.New("underlying error")),
			},
			want: &api.Error{
				Code:    api.ResourceAlreadyExists,
				Message: "An item with the name already exists",
				Type:    api.ApiError,
				Details: map[string]interface{}{
					"resource": "item",
					"name":     "item1",
				},
			},
		},
		"panics if error is not a resource.AlreadyExistsError": {
			args: args{
				err: errors.New("some error"),
			},
			expectPanics: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.expectPanics {
				require.Panics(t, func() {
					server.ErrorMapperAlreadyExistsError(tt.args.err)
				})
				return
			}

			require.Equal(t, tt.want, server.ErrorMapperAlreadyExistsError(tt.args.err))
		})
	}
}

func TestErrorMapperNotFoundError(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}
	tests := map[string]struct {
		args         args
		want         *api.Error
		expectPanics bool
	}{
		"maps resource.NotFoundError to an API error response": {
			args: args{
				err: resource.NewNotFoundError("item1", "item", errors.New("underlying error")),
			},
			want: &api.Error{
				Code:    api.ResourceMissing,
				Message: "The item with the name was not found",
				Type:    api.ApiError,
				Details: map[string]interface{}{
					"resource": "item",
					"name":     "item1",
				},
			},
		},
		"panics if error is not a resource.NotFoundError": {
			args: args{
				err: errors.New("some error"),
			},
			expectPanics: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.expectPanics {
				require.Panics(t, func() {
					server.ErrorMapperNotFoundError(tt.args.err)
				})
				return
			}

			require.Equal(t, tt.want, server.ErrorMapperNotFoundError(tt.args.err))
		})
	}
}

func TestErrorHandler_HandleError(t *testing.T) {
	t.Parallel()

	errorHandler := server.NewErrorHandler()
	errorHandler.RegisterErrorMapper(reflect.TypeOf(&resource.AlreadyExistsError{}), func(_ error) *api.Error {
		return &api.Error{
			Code:    api.ResourceAlreadyExists,
			Message: "ResourceType already exists",
			Type:    api.ApiError,
		}
	})

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   string
		expectedMsg    string
	}{
		{
			name:           "No Error",
			err:            nil,
			expectedStatus: http.StatusOK,
			expectedCode:   "",
			expectedMsg:    "",
		},
		{
			name:           "Mapped Error",
			err:            resource.NewAlreadyExistsError("ResourceType", "ResourceID", errors.New("test error")),
			expectedStatus: http.StatusConflict,
			expectedCode:   string(api.ResourceAlreadyExists),
			expectedMsg:    "ResourceType already exists",
		},
		{
			name:           "Unmapped Error",
			err:            errors.New("test error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   string(api.ProcessingError),
			expectedMsg:    "An error occurred while processing the request.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept", "application/json")
			rr := httptest.NewRecorder()

			// Handle the error
			errorHandler.HandleError(rr, req, tt.err)

			// Assert the response status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// If no error is expected, skip further checks
			if tt.err == nil {
				return
			}

			// Assert the response body
			var errResp api.Error
			if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
				t.Fatal(err)
			}

			if string(errResp.Code) != tt.expectedCode {
				t.Errorf("expected code %s, got %s", tt.expectedCode, errResp.Code)
			}

			if errResp.Message != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, errResp.Message)
			}
		})
	}
}
