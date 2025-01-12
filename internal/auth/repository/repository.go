package repository

import (
	"context"
	"database/sql"

	"github.com/glass-cms/glasscms/internal/auth"
	"github.com/glass-cms/glasscms/internal/auth/repository/query"
	"github.com/glass-cms/glasscms/internal/database"
)

var _ auth.Repository = &TokenRepository{}

type TokenRepository struct {
	db           *sql.DB
	queries      *query.Queries
	errorHandler database.ErrorHandler
}

func NewRepository(db *sql.DB, errorHandler database.ErrorHandler) *TokenRepository {
	return &TokenRepository{
		db:           db,
		queries:      query.New(db),
		errorHandler: errorHandler,
	}
}

func (r *TokenRepository) CreateToken(ctx context.Context, tx *sql.Tx, token auth.Token) error {
	q := r.queries.WithTx(tx)

	params := query.CreateTokenParams{
		ID:         token.ID,
		Suffix:     token.Suffix,
		Hash:       token.Hash,
		ExpireTime: token.ExpireTime,
	}

	if err := q.CreateToken(ctx, params); err != nil {
		return r.errorHandler.HandleError(ctx, err)
	}

	return nil
}

func (r *TokenRepository) GetToken(ctx context.Context, tx *sql.Tx, hash string) (*auth.Token, error) {
	q := r.queries.WithTx(tx)

	// Handle tokens not found.

	dbToken, err := q.GetToken(ctx, hash)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	return &auth.Token{
		ID:         dbToken.ID,
		Suffix:     dbToken.Suffix,
		Hash:       dbToken.Hash,
		CreateTime: dbToken.CreateTime,
		ExpireTime: dbToken.ExpireTime,
	}, nil
}

func (r *TokenRepository) DeleteToken(ctx context.Context, tx *sql.Tx, id string) error {
	q := r.queries.WithTx(tx)

	if err := q.DeleteToken(ctx, id); err != nil {
		return r.errorHandler.HandleError(ctx, err)
	}

	return nil
}
