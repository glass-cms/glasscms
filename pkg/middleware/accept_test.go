package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glass-cms/glasscms/pkg/middleware"
)

func Test_Accept(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		acceptHeader string
		accepted     []string
		expected     int
	}{
		"application/json": {
			acceptHeader: "application/json",
			accepted:     []string{"application/json"},
			expected:     http.StatusOK,
		},
		"application/json with charset": {
			acceptHeader: "application/json; charset=utf-8",
			accepted:     []string{"application/json"},
			expected:     http.StatusOK,
		},
		"application/xml": {
			acceptHeader: "application/xml",
			accepted:     []string{"application/xml"},
			expected:     http.StatusOK,
		},
		"unsupported": {
			acceptHeader: "text/plain",
			accepted:     []string{"application/json"},
			expected:     http.StatusNotAcceptable,
		},
		"invalid": {
			acceptHeader: "text/",
			accepted:     []string{"application/json"},
			expected:     http.StatusBadRequest,
		},
		"empty": {
			acceptHeader: "",
			accepted:     []string{"application/json"},
			expected:     http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Accept", test.acceptHeader)

			rr := httptest.NewRecorder()

			middleware.Accept(test.accepted...)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(rr, req)

			if rr.Code != test.expected {
				t.Errorf("expected %d; got %d", test.expected, rr.Code)
			}
		})
	}
}
