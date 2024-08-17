package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glass-cms/glasscms/lib/mediatype"
	"github.com/glass-cms/glasscms/lib/middleware"
)

func Test_MediaType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		contentHeader string
		accepted      []string
		expected      int
	}{
		"application/json": {
			contentHeader: "application/json",
			accepted:      []string{mediatype.ApplicationJSON},
			expected:      http.StatusOK,
		},
		"application/json with charset": {
			contentHeader: "application/json; charset=utf-8",
			accepted:      []string{mediatype.ApplicationJSON},
			expected:      http.StatusOK,
		},
		"application/xml": {
			contentHeader: "application/xml",
			accepted:      []string{mediatype.ApplicationXML},
			expected:      http.StatusOK,
		},
		"unsupported": {
			contentHeader: "text/plain",
			accepted:      []string{mediatype.ApplicationJSON},
			expected:      http.StatusUnsupportedMediaType,
		},
		"invalid": {
			contentHeader: "text/",
			accepted:      []string{mediatype.ApplicationJSON},
			expected:      http.StatusBadRequest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodPost, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", test.contentHeader)

			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			w := httptest.NewRecorder()
			middleware.MediaType(test.accepted...)(handler).ServeHTTP(w, req)

			if w.Code != test.expected {
				t.Errorf("expected %d, got %d", test.expected, w.Code)
			}
		})
	}
}
