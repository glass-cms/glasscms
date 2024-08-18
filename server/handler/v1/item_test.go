package v1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/lib/log"
	"github.com/glass-cms/glasscms/lib/test"
	v1 "github.com/glass-cms/glasscms/server/handler/v1"
	"github.com/stretchr/testify/assert"
)

func TestAPIHandler_ItemsCreate(t *testing.T) {
	t.Parallel()

	type fields struct {
		repository *item.Repository
	}

	type args struct {
		req      func() *http.Request
		expected int
	}

	tests := map[string]struct {
		fields fields
		args   args
	}{
		"returns a 500 status code when the request body cannot be read": {
			fields: fields{
				repository: &item.Repository{},
			},
			args: args{
				req: func() *http.Request {
					return &http.Request{
						Body: &test.ErrorReadCloser{},
					}
				},
				expected: http.StatusInternalServerError,
			},
		},
		"returns a 400 status code when the buffer cannot be unmarshalled": {
			fields: fields{
				repository: &item.Repository{},
			},
			args: args{
				req: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/v1/items", nil)
				},
				expected: http.StatusBadRequest,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler := v1.NewAPIHandler(
				log.NoopLogger(),
				tt.fields.repository,
			)

			rr := httptest.NewRecorder()
			request := tt.args.req()

			// Act
			handler.ItemsCreate(rr, request)

			// Assert
			assert.Equal(t, tt.args.expected, rr.Code)
		})
	}
}
