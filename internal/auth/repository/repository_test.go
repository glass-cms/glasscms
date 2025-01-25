package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/glass-cms/glasscms/internal/auth"
	"github.com/glass-cms/glasscms/internal/auth/repository"
	"github.com/glass-cms/glasscms/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GetTestDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	if err = database.MigrateDatabase(db, database.Config{
		Driver: "sqlite3",
	}); err != nil {
		panic(err)
	}
	return db
}

func TestRepositoryImpl_CreateToken(t *testing.T) {
	t.Parallel()

	db := GetTestDatabase()
	repo := repository.NewRepository(db, &database.SqliteErrorHandler{})

	token := auth.Token{
		ID:         "1",
		Suffix:     "suffix",
		Hash:       "hash",
		ExpireTime: time.Now().Add(24 * time.Hour),
	}

	tx, err := db.Begin()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback())
	}()

	err = repo.CreateToken(context.Background(), tx, token)
	assert.NoError(t, err)
}

func TestRepositoryImpl_GetToken(t *testing.T) {
	t.Parallel()

	createTestToken := func(
		t testing.TB,
		repo *repository.TokenRepository,
		db *sql.DB,
		id, suffix, hash string,
	) auth.Token {
		token := auth.Token{
			ID:         id,
			Suffix:     suffix,
			Hash:       hash,
			ExpireTime: time.Now().Add(24 * time.Hour),
		}
		tx, err := db.Begin()
		require.NoError(t, err)

		err = repo.CreateToken(context.Background(), tx, token)
		require.NoError(t, err)
		require.NoError(t, tx.Commit())

		return token
	}

	tests := []struct {
		name        string
		withoutTx   bool
		createToken func(testing.TB, *repository.TokenRepository, *sql.DB) auth.Token
		wantErr     error
	}{
		{
			name:      "with transaction - token exists",
			withoutTx: false,
			createToken: func(_ testing.TB, repo *repository.TokenRepository, db *sql.DB) auth.Token {
				return createTestToken(t, repo, db, "1", "suffix", "hash")
			},
		},
		{
			name:      "without transaction - token exists",
			withoutTx: true,
			createToken: func(_ testing.TB, repo *repository.TokenRepository, db *sql.DB) auth.Token {
				return createTestToken(t, repo, db, "2", "suffix2", "hash2")
			},
		},
		{
			name:      "token not found",
			withoutTx: false,
			createToken: func(_ testing.TB, _ *repository.TokenRepository, _ *sql.DB) auth.Token {
				return auth.Token{
					Hash: "nonexistent-hash",
				}
			},
			wantErr: database.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := GetTestDatabase()
			repo := repository.NewRepository(db, &database.SqliteErrorHandler{})

			token := tt.createToken(t, repo, db)

			var fetchedToken *auth.Token
			var getErr error

			if tt.withoutTx {
				fetchedToken, getErr = repo.GetToken(context.Background(), nil, token.Hash)
			} else {
				tx, txErr := db.Begin()
				require.NoError(t, txErr)
				defer func() {
					require.NoError(t, tx.Rollback())
				}()
				fetchedToken, getErr = repo.GetToken(context.Background(), tx, token.Hash)
			}

			if tt.wantErr != nil {
				require.ErrorIs(t, getErr, tt.wantErr)
				require.Nil(t, fetchedToken)
				return
			}

			require.NoError(t, getErr)
			require.NotNil(t, fetchedToken, "fetched token should not be nil")

			assert.Equal(t, token.ID, fetchedToken.ID)
			assert.Equal(t, token.Suffix, fetchedToken.Suffix)
			assert.Equal(t, token.Hash, fetchedToken.Hash)
		})
	}
}

func TestRepositoryImpl_DeleteToken(t *testing.T) {
	t.Parallel()
	db := GetTestDatabase()
	repo := repository.NewRepository(db, &database.SqliteErrorHandler{})

	token := auth.Token{
		ID:         "1",
		Suffix:     "suffix",
		Hash:       "hash",
		ExpireTime: time.Now().Add(24 * time.Hour),
	}

	tx, err := db.Begin()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tx.Rollback())
	}()

	err = repo.CreateToken(context.Background(), tx, token)
	require.NoError(t, err)

	err = repo.DeleteToken(context.Background(), tx, token.ID)
	require.NoError(t, err)

	_, err = repo.GetToken(context.Background(), tx, token.Hash)
	assert.Error(t, err)
}
