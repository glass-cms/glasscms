package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/item/repository"
	"github.com/glass-cms/glasscms/internal/item/repository/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertQueryItem(t *testing.T) {
	tests := []struct {
		name     string
		dbItem   query.Item
		expected *item.Item
		wantErr  bool
	}{
		{
			name: "Valid conversion",
			dbItem: query.Item{
				Name:        "test_name",
				DisplayName: "Test Name",
				Content:     sql.NullString{String: "test content", Valid: true},
				Hash:        sql.NullString{String: "test hash", Valid: true},
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
				DeleteTime:  sql.NullTime{Time: time.Now(), Valid: true},
				Properties:  `{"key1": "value1"}`,
				Metadata:    `{"meta1": "data1"}`,
			},
			expected: &item.Item{
				Name:        "test_name",
				DisplayName: "Test Name",
				Content:     "test content",
				Hash:        "test hash",
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
				DeleteTime:  func() *time.Time { t := time.Now(); return &t }(),
				Properties:  map[string]any{"key1": "value1"},
				Metadata:    map[string]any{"meta1": "data1"},
			},
			wantErr: false,
		},
		{
			name: "Invalid Properties JSON",
			dbItem: query.Item{
				Name:        "test_name",
				DisplayName: "Test Name",
				Content:     sql.NullString{String: "test content", Valid: true},
				Hash:        sql.NullString{String: "test hash", Valid: true},
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
				DeleteTime:  sql.NullTime{Time: time.Now(), Valid: true},
				Properties:  `{"key1": "value1"`,
				Metadata:    `{"meta1": "data1"}`,
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Invalid Metadata JSON",
			dbItem: query.Item{
				Name:        "test_name",
				DisplayName: "Test Name",
				Content:     sql.NullString{String: "test content", Valid: true},
				Hash:        sql.NullString{String: "test hash", Valid: true},
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
				DeleteTime:  sql.NullTime{Time: time.Now(), Valid: true},
				Properties:  `{"key1": "value1"}`,
				Metadata:    `{"meta1": "data1"`,
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repository.ConvertQueryItem(tt.dbItem)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Name, got.Name)
				assert.Equal(t, tt.expected.DisplayName, got.DisplayName)
				assert.Equal(t, tt.expected.Content, got.Content)
				assert.Equal(t, tt.expected.Hash, got.Hash)
				assert.WithinDuration(t, tt.expected.CreateTime, got.CreateTime, time.Second)
				assert.WithinDuration(t, tt.expected.UpdateTime, got.UpdateTime, time.Second)
				assert.WithinDuration(t, *tt.expected.DeleteTime, *got.DeleteTime, time.Second)
				assert.Equal(t, tt.expected.Properties, got.Properties)
				assert.Equal(t, tt.expected.Metadata, got.Metadata)
			}
		})
	}
}
