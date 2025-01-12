package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
)

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

	// TODO: Handle tokens not found.
	// TODO: Handle tokens expired.

	dbToken, err := a.repo.GetToken(ctx, tx, hash)
	if err != nil {
		return false, err
	}
	return dbToken != nil, nil
}
