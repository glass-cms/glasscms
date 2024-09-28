package database_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/glass-cms/glasscms/internal/database"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestSqliteErrorHandler_HandleError(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}

	tests := map[string]struct {
		args args
		want error
	}{
		"ErrNoRows": {
			args: args{
				err: sql.ErrNoRows,
			},
			want: database.ErrNotFound,
		},
		"ErrConstraintPrimaryKey": {
			args: args{
				err: sqlite3.Error{
					Code:         sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintPrimaryKey,
				},
			},
			want: database.ErrDuplicatePrimaryKey,
		},
		"ErrConstraintUnique": {
			args: args{
				err: sqlite3.Error{
					Code:         sqlite3.ErrConstraint,
					ExtendedCode: sqlite3.ErrConstraintUnique,
				},
			},
			want: database.ErrUniqueConstraint,
		},
		"ErrOperationFailed": {
			args: args{
				err: errors.New("some error"),
			},
			want: database.ErrOperationFailed,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &database.SqliteErrorHandler{}

			got := e.HandleError(context.TODO(), tt.args.err)
			assert.ErrorIs(t, got, tt.want)
		})
	}
}
