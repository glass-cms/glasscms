package repository_test

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/item/repository"
	"github.com/glass-cms/glasscms/internal/item/repository/query"
)

func TestConvertDatabaseItem(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name   string
		dbItem query.Item
		want   item.Item
	}{
		{
			name: "all fields valid",
			dbItem: query.Item{
				Name:        "test_name",
				DisplayName: "Test Name",
				Content:     sql.NullString{String: "test content", Valid: true},
				Hash:        sql.NullString{String: "test hash", Valid: true},
				CreateTime:  now,
				UpdateTime:  now,
				DeleteTime:  sql.NullTime{Time: now, Valid: true},
				Properties:  map[string]any{"key1": "value1"},
				Metadata:    map[string]any{"meta1": "data1"},
			},
			want: item.Item{
				Name:        "test_name",
				DisplayName: "Test Name",
				Content:     "test content",
				Hash:        "test hash",
				CreateTime:  now,
				UpdateTime:  now,
				DeleteTime:  &now,
				Properties:  map[string]any{"key1": "value1"},
				Metadata:    map[string]any{"meta1": "data1"},
			},
		},
		{
			name: "null fields",
			dbItem: query.Item{
				Name:        "test_name",
				DisplayName: "Test Name",
				Content:     sql.NullString{Valid: false},
				Hash:        sql.NullString{Valid: false},
				CreateTime:  now,
				UpdateTime:  now,
				DeleteTime:  sql.NullTime{Valid: false},
				Properties:  nil,
				Metadata:    nil,
			},
			want: item.Item{
				Name:        "test_name",
				DisplayName: "Test Name",
				Content:     "",
				Hash:        "",
				CreateTime:  now,
				UpdateTime:  now,
				DeleteTime:  nil,
				Properties:  nil,
				Metadata:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := repository.ConvertQueryItem(tt.dbItem); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertDatabaseItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
