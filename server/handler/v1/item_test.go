package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/glass-cms/glasscms/api/v1"
	"github.com/glass-cms/glasscms/database"
	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/lib/log"
	"github.com/glass-cms/glasscms/lib/test"
	v1 "github.com/glass-cms/glasscms/server/handler/v1"
	"github.com/stretchr/testify/assert"
)

func TestAPIHandler_ItemsCreate(t *testing.T) {
	t.Parallel()

	testdb, err := test.NewDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := item.NewRepository(testdb, &database.SqliteErrorHandler{})

	tests := map[string]struct {
		req      func() *http.Request
		expected int
	}{
		"returns a 500 status code when the request body cannot be read": {
			req: func() *http.Request {
				return &http.Request{
					Body: &test.ErrorReadCloser{},
				}
			},
			expected: http.StatusInternalServerError,
		},
		"returns a 400 status code when the buffer cannot be unmarshalled": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/v1/items", nil)
			},
			expected: http.StatusBadRequest,
		},
		"returns a 201 status code when the item is created successfully": {
			req: func() *http.Request {
				item := &api.ItemsCreateJSONRequestBody{
					Content: "content",
					Name:    "name",
				}
				body, _ := json.Marshal(item)
				return httptest.NewRequest(http.MethodPost, "/v1/items", bytes.NewReader(body))
			},
			expected: http.StatusCreated,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := v1.NewAPIHandler(
				log.NoopLogger(),
				item.NewService(repo),
			)

			rr := httptest.NewRecorder()
			request := tt.req()

			// Act
			handler.ItemsCreate(rr, request)

			// Assert
			assert.Equal(t, tt.expected, rr.Code)
		})
	}
}
