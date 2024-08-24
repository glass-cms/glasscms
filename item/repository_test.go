package item_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/glass-cms/glasscms/database"
	"github.com/glass-cms/glasscms/item"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
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

func SeedDatabase(db *sql.DB, items ...*item.Item) error {
	repo := item.NewRepository(db, &database.SqliteErrorHandler{})
	for _, i := range items {
		if err := repo.CreateItem(context.Background(), i); err != nil {
			return err
		}
	}
	return nil
}

func getTestItem(name string) *item.Item {
	return &item.Item{
		UID:         "1234",
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		Hash:        "hash",
		Name:        name,
		DisplayName: "DisplayName",
		Content:     "Content",
		Properties:  map[string]interface{}{"key": "value"},
	}
}

func TestRepository_CreateItem(t *testing.T) {
	t.Parallel()

	type fields struct {
		db   *sql.DB
		seed func(*sql.DB)
	}

	type args struct {
		ctx  context.Context
		item *item.Item
	}

	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
		err     error
	}{
		"successful creation": {
			fields: fields{
				db: GetTestDatabase(),
			},
			args: args{
				ctx:  context.Background(),
				item: getTestItem("items/name"),
			},
			wantErr: false,
			err:     nil,
		},
		"returns an err when context is canceled": {
			fields: fields{
				db: GetTestDatabase(),
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				item: getTestItem("items/name"),
			},
			wantErr: true,
			err:     database.ErrOperationFailed,
		},
		"returns an error when properties cannot be marshalled": {
			fields: fields{
				db: GetTestDatabase(),
			},
			args: args{
				ctx: context.Background(),
				item: &item.Item{
					UID:         "1234",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
					Hash:        "hash",
					Name:        "items/name",
					DisplayName: "DisplayName",
					Content:     "Content",
					Properties:  map[string]interface{}{"key": make(chan int)},
				},
			},
			wantErr: true,
			err:     database.ErrOperationFailed,
		},
		"returns an error when metadata cannot be marshalled": {
			fields: fields{
				db: GetTestDatabase(),
			},
			args: args{
				ctx: context.Background(),
				item: &item.Item{
					UID:         "1234",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
					Hash:        "hash",
					Name:        "items/name",
					DisplayName: "DisplayName",
					Content:     "Content",
					Properties:  map[string]interface{}{"key": "value"},
					Metadata:    map[string]interface{}{"key": make(chan int)},
				},
			},
			wantErr: true,
			err:     database.ErrOperationFailed,
		},
		"returns an error when key already exists": {
			fields: fields{
				db: GetTestDatabase(),
				seed: func(db *sql.DB) {
					if err := SeedDatabase(db, getTestItem("items/name")); err != nil {
						t.Error(err)
					}
				},
			},
			args: args{
				ctx:  context.Background(),
				item: getTestItem("items/name2"),
			},
			wantErr: true,
			err:     database.ErrDuplicatePrimaryKey,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			if tt.fields.seed != nil {
				tt.fields.seed(tt.fields.db)
			}
			r := item.NewRepository(tt.fields.db, &database.SqliteErrorHandler{})

			// Act
			err := r.CreateItem(tt.args.ctx, tt.args.item)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestRepository_GetItem(t *testing.T) {
	t.Parallel()

	type fields struct {
		db   *sql.DB
		seed func(*sql.DB)
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    *item.Item
		wantErr bool
	}{
		"Successful retrieval": {
			fields: fields{
				db: GetTestDatabase(),
				seed: func(db *sql.DB) {
					if err := SeedDatabase(db, getTestItem("items/name")); err != nil {
						t.Error(err)
					}
				},
			},
			args: args{
				ctx: context.Background(),
				uid: "1234",
			},
			want:    getTestItem("items/name"),
			wantErr: false,
		},
		"Context canceled": {
			fields: fields{
				db: GetTestDatabase(),
				seed: func(db *sql.DB) {
					if err := SeedDatabase(db, getTestItem("items/name")); err != nil {
						t.Error(err)
					}
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				uid: "1234",
			},
			want:    nil,
			wantErr: true,
		},
		"Item not found": {
			fields: fields{
				db: GetTestDatabase(),
			},
			args: args{
				ctx: context.Background(),
				uid: "nonexistent",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := item.NewRepository(tt.fields.db, &database.SqliteErrorHandler{})
			if tt.fields.seed != nil {
				tt.fields.seed(tt.fields.db)
			}
			got, err := r.GetItem(tt.args.ctx, tt.args.uid)
			assert.Equal(t, tt.wantErr, err != nil, "Repository.GetItem() error = %v, wantErr %v", err, tt.wantErr)
			if tt.want != nil && got != nil {
				assert.Equal(t, tt.want.UID, got.UID)
				assert.WithinDuration(t, tt.want.CreateTime, got.CreateTime, time.Second)
				assert.WithinDuration(t, tt.want.UpdateTime, got.UpdateTime, time.Second)
				assert.Equal(t, tt.want.Hash, got.Hash)
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.DisplayName, got.DisplayName)
				assert.Equal(t, tt.want.Content, got.Content)
				assert.Equal(t, tt.want.Properties, got.Properties)
				assert.Equal(t, tt.want.Metadata, got.Metadata)
			}
		})
	}
}

func TestRepository_UpdateItem(t *testing.T) {
	t.Parallel()

	type fields struct {
		db   *sql.DB
		seed func(*sql.DB)
	}
	type args struct {
		ctx  context.Context
		item *item.Item
	}
	tests := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"Successful update": {
			fields: fields{
				db: GetTestDatabase(),
				seed: func(db *sql.DB) {
					if err := SeedDatabase(db, getTestItem("items/name")); err != nil {
						t.Error(err)
					}
				},
			},
			args: args{
				ctx: context.Background(),
				item: &item.Item{
					UID:         "1234",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
					Hash:        "newhash",
					Name:        "NewName",
					DisplayName: "NewDisplayName",
					Content:     "NewContent",
					Properties:  map[string]interface{}{"newkey": "newvalue"},
				},
			},
			wantErr: false,
		},
		"Context canceled": {
			fields: fields{
				db: GetTestDatabase(),
				seed: func(db *sql.DB) {
					if err := SeedDatabase(db, getTestItem("items/name")); err != nil {
						t.Error(err)
					}
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				}(),
				item: &item.Item{
					UID:         "1234",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
					Hash:        "newhash",
					Name:        "NewName",
					DisplayName: "NewDisplayName",
					Content:     "NewContent",
					Properties:  map[string]interface{}{"newkey": "newvalue"},
				},
			},
			wantErr: true,
		},
		"Update non-existent item": {
			fields: fields{
				db: GetTestDatabase(),
			},
			args: args{
				ctx: context.Background(),
				item: &item.Item{
					UID:         "nonexistent",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
					Hash:        "newhash",
					Name:        "NewName",
					DisplayName: "NewDisplayName",
					Content:     "NewContent",
					Properties:  map[string]interface{}{"newkey": "newvalue"},
				},
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := item.NewRepository(tt.fields.db, &database.SqliteErrorHandler{})
			if tt.fields.seed != nil {
				tt.fields.seed(tt.fields.db)
			}
			err := r.UpdateItem(tt.args.ctx, tt.args.item)
			assert.Equal(t, tt.wantErr, err != nil, "Repository.UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
