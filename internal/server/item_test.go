package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/item/repository"
	"github.com/glass-cms/glasscms/internal/server"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestAPIHandler_ItemsCreate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		req      func() *http.Request
		expected int
	}{
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

			testdb, err := database.NewTestDB()
			if err != nil {
				t.Fatal(err)
			}
			defer testdb.Close()

			repo := repository.NewRepository(testdb, &database.SqliteErrorHandler{})

			handler, err := server.New(
				log.NoopLogger(),
				item.NewService(repo),
				[]func(http.Handler) http.Handler{},
			)
			if err != nil {
				t.Fatal(err)
				return
			}

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

			repo := repository.NewRepository(testdb, &database.SqliteErrorHandler{})

			server, err := server.New(
				log.NoopLogger(),
				item.NewService(repo),
				[]func(http.Handler) http.Handler{},
			)
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			request := tt.req()
			request.Header.Set("Accept", "application/json")

			// Make the request
			server.Handler().ServeHTTP(rr, request)

			// Assert
			assert.Equal(t, tt.expected, rr.Code)
		})
	}
}

//nolint:gocognit // This test is testing multiple cases.
func TestAPIHandler_ItemsList(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		req      func() *http.Request
		seed     func(*item.Service)
		expected int
	}{
		"returns a 200 status code and a list of items": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/v1/items", nil)
			},
			seed: func(svc *item.Service) {
				items := []item.Item{
					{Name: "items/name1", DisplayName: "Item 1"},
					{Name: "items/name2", DisplayName: "Item 2"},
				}
				for _, itm := range items {
					if _, err := svc.CreateItem(context.Background(), itm); err != nil {
						t.Error(err)
					}
				}
			},
			expected: http.StatusOK,
		},
		"returns a 200 status code and a list of items with fieldmask": {
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/v1/items?fields=name,display_name", nil)
				return req
			},
			seed: func(svc *item.Service) {
				items := []item.Item{
					{Name: "items/name1", DisplayName: "Item 1"},
					{Name: "items/name2", DisplayName: "Item 2"},
				}
				for _, itm := range items {
					if _, err := svc.CreateItem(context.Background(), itm); err != nil {
						t.Error(err)
					}
				}
			},
			expected: http.StatusOK,
		},
		"returns a 400 status code when fieldmask is invalid": {
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/v1/items?fields=invalid_field", nil)
				return req
			},
			seed:     nil,
			expected: http.StatusBadRequest,
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

			repo := repository.NewRepository(testdb, &database.SqliteErrorHandler{})
			svc := item.NewService(repo)

			if tt.seed != nil {
				tt.seed(svc)
			}

			handler, err := server.New(
				log.NoopLogger(),
				svc,
				[]func(http.Handler) http.Handler{},
			)
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			request := tt.req()
			request.Header.Set("Accept", "application/json")

			// Act
			fields := request.URL.Query().Get("fields")
			var params api.ItemsListParams
			if fields != "" {
				// Split the fields by comma
				splitFields := strings.Split(fields, ",")
				params.Fields = &splitFields
			}
			handler.ItemsList(rr, request, params)

			// Assert
			assert.Equal(t, tt.expected, rr.Code)
		})
	}
}

func TestAPIHandler_ItemsUpsert(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		req      func() *http.Request
		seed     func(*item.Service)
		expected int
	}{
		"returns a 200 status code when items are upserted successfully": {
			req: func() *http.Request {
				items := api.ItemsUpsertJSONRequestBody{
					{
						Name:        "items/name1",
						DisplayName: "Item 1",
						Content:     "content1",
					},
					{
						Name:        "items/name2",
						DisplayName: "Item 2",
						Content:     "content2",
					},
				}
				body, _ := json.Marshal(items)
				return httptest.NewRequest(http.MethodPost, "/v1/items/upsert", bytes.NewReader(body))
			},
			seed:     nil,
			expected: http.StatusOK,
		},
		"returns a 400 status code when request body is invalid": {
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/v1/items/upsert", bytes.NewReader([]byte("invalid body")))
			},
			seed:     nil,
			expected: http.StatusBadRequest,
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

			repo := repository.NewRepository(testdb, &database.SqliteErrorHandler{})
			svc := item.NewService(repo)

			if tt.seed != nil {
				tt.seed(svc)
			}

			handler, err := server.New(
				log.NoopLogger(),
				svc,
				[]func(http.Handler) http.Handler{},
			)
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			request := tt.req()
			request.Header.Set("Accept", "application/json")

			// Act
			handler.ItemsUpsert(rr, request)

			// Assert
			assert.Equal(t, tt.expected, rr.Code)
		})
	}
}
