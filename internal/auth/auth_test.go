package auth_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/glass-cms/glasscms/internal/auth"
	"github.com/glass-cms/glasscms/internal/auth/repository"
	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestAuth(t *testing.T) (*auth.Auth, *sql.DB, auth.Repository) {
	t.Helper()
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewRepository(db, &database.SqliteErrorHandler{})
	a := auth.NewAuth(db, repo, log.NoopLogger())
	return a, db, repo
}

func createToken(t *testing.T, db *sql.DB, repo auth.Repository, expireTime time.Time) string {
	t.Helper()
	token, tokenValue := auth.NewToken(expireTime)

	tx, txErr := db.Begin()
	if txErr != nil {
		t.Fatal(txErr)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			t.Fatal(rollbackErr)
		}
	}()

	if createErr := repo.CreateToken(context.Background(), tx, *token); createErr != nil {
		t.Fatal(createErr)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		t.Fatal(commitErr)
	}

	return tokenValue
}

func TestValidateToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupToken    func(_ *testing.T, db *sql.DB, repo auth.Repository) string
		expectedValid bool
		expectedErr   error
	}{
		{
			name: "valid token",
			setupToken: func(_ *testing.T, db *sql.DB, repo auth.Repository) string {
				tokenValue := createToken(t, db, repo, time.Now().Add(24*time.Hour))
				return fmt.Sprintf("Bearer %s", tokenValue)
			},
			expectedValid: true,
			expectedErr:   nil,
		},
		{
			name: "invalid token",
			setupToken: func(_ *testing.T, _ *sql.DB, _ auth.Repository) string {
				return "invalid_token"
			},
			expectedValid: false,
			expectedErr:   auth.ErrTokenNotFound,
		},
		{
			name: "expired token",
			setupToken: func(_ *testing.T, db *sql.DB, repo auth.Repository) string {
				tokenValue := createToken(t, db, repo, time.Now().Add(-24*time.Hour))
				return fmt.Sprintf("Bearer %s", tokenValue)
			},
			expectedValid: false,
			expectedErr:   auth.ErrTokenExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a, db, repo := setupTestAuth(t)

			tokenString := tt.setupToken(t, db, repo)
			valid, err := a.ValidateToken(context.Background(), tokenString)

			assert.Equal(t, tt.expectedValid, valid)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		expireTime  time.Time
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "successful token creation",
			expireTime:  time.Now().Add(24 * time.Hour),
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "token with past expiration",
			expireTime:  time.Now().Add(-24 * time.Hour),
			wantErr:     true,
			expectedErr: auth.ErrInvalidExpireTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a, _, _ := setupTestAuth(t)

			token, prettyValue, err := a.CreateToken(context.Background(), tt.expireTime)

			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedErr != nil {
					require.ErrorIs(t, err, tt.expectedErr)
				}

				assert.Nil(t, token)
				assert.Empty(t, prettyValue)

				return
			}

			require.NoError(t, err)

			assert.NotNil(t, token)
			assert.NotEmpty(t, prettyValue)
			assert.True(t, strings.HasPrefix(prettyValue, "sk_"))
			assert.Equal(t, tt.expireTime.Unix(), token.ExpireTime.Unix())
		})
	}
}
