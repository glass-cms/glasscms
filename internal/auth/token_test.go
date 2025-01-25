package auth_test

import (
	"strings"
	"testing"
	"time"

	"github.com/glass-cms/glasscms/internal/auth"
)

func TestNewToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		expireTime time.Time
		check      func(*testing.T, *auth.Token, string)
	}{
		{
			name:       "token ID should not be empty",
			expireTime: time.Now().Add(24 * time.Hour),
			check: func(t *testing.T, token *auth.Token, _ string) {
				if token.ID == "" {
					t.Error("token ID should not be empty")
				}
			},
		},
		{
			name:       "token hash should not be empty",
			expireTime: time.Now().Add(24 * time.Hour),
			check: func(t *testing.T, token *auth.Token, _ string) {
				if token.Hash == "" {
					t.Error("token hash should not be empty")
				}
			},
		},
		{
			name:       "token suffix should not be empty",
			expireTime: time.Now().Add(24 * time.Hour),
			check: func(t *testing.T, token *auth.Token, _ string) {
				if token.Suffix == "" {
					t.Error("token suffix should not be empty")
				}
			},
		},
		{
			name:       "pretty value should start with sk_",
			expireTime: time.Now().Add(24 * time.Hour),
			check: func(t *testing.T, _ *auth.Token, prettyValue string) {
				if !strings.HasPrefix(prettyValue, "sk_") {
					t.Errorf("pretty value should start with 'sk_', got %s", prettyValue)
				}
			},
		},
		{
			name:       "token suffix should match pretty value suffix",
			expireTime: time.Now().Add(24 * time.Hour),
			check: func(t *testing.T, token *auth.Token, prettyValue string) {
				if prettyValue[len(prettyValue)-4:] != token.Suffix {
					t.Error("token suffix should match pretty value suffix")
				}
			},
		},
		{
			name:       "create time should not be in the future",
			expireTime: time.Now().Add(24 * time.Hour),
			check: func(t *testing.T, token *auth.Token, _ string) {
				if token.CreateTime.After(time.Now()) {
					t.Error("create time should not be in the future")
				}
			},
		},
		{
			name:       "hash should be 64 characters long",
			expireTime: time.Now().Add(24 * time.Hour),
			check: func(t *testing.T, token *auth.Token, _ string) {
				if len(token.Hash) != 64 {
					t.Errorf("hash should be 64 characters long, got %d", len(token.Hash))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token, prettyValue := auth.NewToken(tt.expireTime)
			tt.check(t, token, prettyValue)
		})
	}
}
