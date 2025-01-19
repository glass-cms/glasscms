package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/glass-cms/glasscms/internal/database"
)

// ErrTokenNotFound is returned when a token cannot be found in the database.
var ErrTokenNotFound = errors.New("token not found")

// ErrTokenExpired is returned when a token's expiration time has passed.
var ErrTokenExpired = errors.New("token expired")

// Auth is the service that handles token generation and validation.
type Auth struct {
	db     *sql.DB
	repo   Repository
	logger *slog.Logger
}

// NewAuth creates a new Auth service.
func NewAuth(db *sql.DB, repo Repository, logger *slog.Logger) *Auth {
	return &Auth{db: db, repo: repo, logger: logger}
}

// ValidateToken validates a token and returns true if it is valid, false otherwise.
func (a *Auth) ValidateToken(ctx context.Context, token string) (bool, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimPrefix(token, "sk_")

	hash := tokenHash(token)

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			a.logger.Error("failed to rollback transaction", "error", rollbackErr)
		}
	}()

	dbToken, err := a.repo.GetToken(ctx, tx, hash)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return false, ErrTokenNotFound
		}

		return false, err
	}

	if dbToken.ExpireTime.Before(time.Now()) {
		return false, ErrTokenExpired
	}

	return dbToken != nil, nil
}
