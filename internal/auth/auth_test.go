package auth_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/glass-cms/glasscms/internal/auth"
	"github.com/glass-cms/glasscms/internal/auth/repository"
	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	t.Parallel()
	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewRepository(db, &database.SqliteErrorHandler{})
	a := auth.NewAuth(db, repo, log.NoopLogger())

	token, tokenValue := auth.NewToken(time.Now().Add(24 * time.Hour))

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

	valid, validateErr := a.ValidateToken(context.Background(), fmt.Sprintf("Bearer %s", tokenValue))
	if validateErr != nil {
		t.Fatal(validateErr)
	}

	assert.True(t, valid)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	t.Parallel()

	db, err := database.NewTestDB()
	if err != nil {
		t.Fatal(err)
	}

	repo := repository.NewRepository(db, &database.SqliteErrorHandler{})
	a := auth.NewAuth(db, repo, log.NoopLogger())

	valid, err := a.ValidateToken(context.Background(), "invalid_token")
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, valid)
}
