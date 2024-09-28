package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/glass-cms/glasscms/api/v1"
	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/internal/item"
	v1 "github.com/glass-cms/glasscms/internal/server/handler/v1"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestAPIHandler_ItemsCreate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		req      func() *http.Request
		expected int
	}{
		"returns a 500 status code when the buffer cannot be unmarshalled": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/v1/items", nil)
			},
			expected: http.StatusInternalServerError,
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

			testdb, err := database.NewTestDB()
			if err != nil {
				t.Fatal(err)
			}
			defer testdb.Close()

			repo := item.NewRepository(testdb, &database.SqliteErrorHandler{})

			handler := v1.NewAPIHandler(
				log.NoopLogger(),
				item.NewService(repo),
			)

			rr := httptest.NewRecorder()
			request := tt.req()
			request.Header.Set("Accept", "application/json")

			// Act
			handler.ItemsCreate(rr, request)

			// Assert
			assert.Equal(t, tt.expected, rr.Code)
		})
	}
}

func TestAPIHandler_ItemsGet(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		req      func() *http.Request
		expected int
	}{
		"returns a 404 status code when the item cannot be found": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/v1/items/missing", nil)
			},
			expected: http.StatusNotFound,
		},
		// TODO: Add more tests.
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testdb, err := database.NewTestDB()
			if err != nil {
				t.Fatal(err)
			}
			defer testdb.Close()

			repo := item.NewRepository(testdb, &database.SqliteErrorHandler{})

			handler := v1.NewAPIHandler(
				log.NoopLogger(),
				item.NewService(repo),
			).Handler(http.NewServeMux(), []func(http.Handler) http.Handler{})

			rr := httptest.NewRecorder()
			request := tt.req()
			request.Header.Set("Accept", "application/json")

			// Make the request
			handler.ServeHTTP(rr, request)

			// Assert
			assert.Equal(t, tt.expected, rr.Code)
		})
	}
}
