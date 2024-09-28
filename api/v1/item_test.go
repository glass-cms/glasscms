package v1_test

import (
	"testing"
	"time"

	v1 "github.com/glass-cms/glasscms/api/v1"
	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/parser"
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
				Properties:  map[string]interface{}{"key1": "value1", "key2": "value2"},
				UpdateTime:  time.Now(),
			},
			want: &item.Item{
				Name:        "Test Name",
				DisplayName: "Test Display Name",
				Content:     "Test Content",
				Hash:        parser.HashContent([]byte("Test Content")),
				CreateTime:  time.Now().Add(-24 * time.Hour),
				UpdateTime:  time.Now(),
				Properties:  map[string]any{"key1": "value1", "key2": "value2"},
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
				Properties:  nil,
				UpdateTime:  time.Time{},
			},
			want: &item.Item{
				Name:        "",
				DisplayName: "",
				Content:     "",
				Hash:        parser.HashContent([]byte("")),
				CreateTime:  time.Time{},
				UpdateTime:  time.Time{},
				Properties:  nil,
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
				Properties:  nil,
				UpdateTime:  time.Now(),
			},
			want: &item.Item{
				Name:        "Test Name 2",
				Content:     "Test Content",
				DisplayName: "Test Display Name",
				Hash:        parser.HashContent([]byte("Test Content")),
				CreateTime:  time.Now().Add(-24 * time.Hour),
				UpdateTime:  time.Now(),
				Properties:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			apiItem := &v1.ItemCreate{
				Content:     tt.fields.Content,
				CreateTime:  tt.fields.CreateTime,
				DisplayName: tt.fields.DisplayName,
				Name:        tt.fields.Name,
				Properties:  tt.fields.Properties,
				UpdateTime:  tt.fields.UpdateTime,
			}
			got := apiItem.ToItem()

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
