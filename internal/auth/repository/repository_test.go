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

	fetchedToken, err := repo.GetToken(context.Background(), tx, token.Hash)
	require.NoError(t, err)
	require.NotNil(t, fetchedToken, "fetched token should not be nil")
	assert.Equal(t, token.ID, fetchedToken.ID)
	assert.Equal(t, token.Suffix, fetchedToken.Suffix)
	assert.Equal(t, token.Hash, fetchedToken.Hash)
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
