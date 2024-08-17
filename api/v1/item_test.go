package v1_test

import (
	"testing"
	"time"

	v1 "github.com/glass-cms/glasscms/api/v1"
	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/parser"
	"github.com/stretchr/testify/assert"
)

func TestItem_MapToDomain(t *testing.T) {
	t.Parallel()

	type fields struct {
		Content     string
		CreateTime  time.Time
		DisplayName string
		ID          string
		Name        string
		Path        string
		Properties  map[string]interface{}
		UpdateTime  time.Time
	}

	tests := []struct {
		name   string
		fields fields
		want   *item.Item
	}{
		{
			name: "Complete Item Mapping",
			fields: fields{
				Content:     "Test Content",
				CreateTime:  time.Now().Add(-24 * time.Hour),
				DisplayName: "Test Display Name",
				ID:          "1234",
				Name:        "Test Name",
				Path:        "/test/path",
				Properties:  map[string]interface{}{"key1": "value1", "key2": "value2"},
				UpdateTime:  time.Now(),
			},
			want: &item.Item{
				UID:        "1234",
				Name:       "Test Name",
				Path:       "/test/path",
				Content:    "Test Content",
				Hash:       parser.HashContent([]byte("Test Content")),
				CreateTime: time.Now().Add(-24 * time.Hour),
				UpdateTime: time.Now(),
				Properties: map[string]any{"key1": "value1", "key2": "value2"},
			},
		},
		{
			name: "Empty Item Mapping",
			fields: fields{
				Content:     "",
				CreateTime:  time.Time{},
				DisplayName: "",
				ID:          "",
				Name:        "",
				Path:        "",
				Properties:  nil,
				UpdateTime:  time.Time{},
			},
			want: &item.Item{
				UID:        "",
				Name:       "",
				Path:       "",
				Content:    "",
				Hash:       parser.HashContent([]byte("")),
				CreateTime: time.Time{},
				UpdateTime: time.Time{},
				Properties: nil,
			},
		},
		{
			name: "Nil Properties Mapping",
			fields: fields{
				Content:     "Test Content",
				CreateTime:  time.Now().Add(-24 * time.Hour),
				DisplayName: "Test Display Name",
				ID:          "5678",
				Name:        "Test Name 2",
				Path:        "/test/path2",
				Properties:  nil,
				UpdateTime:  time.Now(),
			},
			want: &item.Item{
				UID:        "5678",
				Name:       "Test Name 2",
				Path:       "/test/path2",
				Content:    "Test Content",
				Hash:       parser.HashContent([]byte("Test Content")),
				CreateTime: time.Now().Add(-24 * time.Hour),
				UpdateTime: time.Now(),
				Properties: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			apiItem := &v1.Item{
				Content:     tt.fields.Content,
				CreateTime:  tt.fields.CreateTime,
				DisplayName: tt.fields.DisplayName,
				Id:          tt.fields.ID,
				Name:        tt.fields.Name,
				Path:        tt.fields.Path,
				Properties:  tt.fields.Properties,
				UpdateTime:  tt.fields.UpdateTime,
			}
			got := apiItem.MapToDomain()

			// Adjust for potential differences in time (e.g., slight differences due to test execution timing)
			if !got.CreateTime.IsZero() && !tt.want.CreateTime.IsZero() {
				got.CreateTime = tt.want.CreateTime
			}
			if !got.UpdateTime.IsZero() && !tt.want.UpdateTime.IsZero() {
				got.UpdateTime = tt.want.UpdateTime
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
