package api_test

import (
	"testing"

	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		content    string
		properties map[string]interface{}
		metadata   map[string]interface{}
		wantErr    bool
	}{
		{
			name:       "valid input",
			content:    "example content",
			properties: map[string]interface{}{"key1": "value1"},
			metadata:   map[string]interface{}{"meta1": "data1"},
			wantErr:    false,
		},
		{
			name:       "empty content",
			content:    "",
			properties: map[string]interface{}{"key1": "value1"},
			metadata:   map[string]interface{}{"meta1": "data1"},
			wantErr:    false,
		},
		{
			name:       "nil properties",
			content:    "example content",
			properties: nil,
			metadata:   map[string]interface{}{"meta1": "data1"},
			wantErr:    false,
		},
		{
			name:       "nil metadata",
			content:    "example content",
			properties: map[string]interface{}{"key1": "value1"},
			metadata:   nil,
			wantErr:    false,
		},
		{
			name:       "serialization error in properties",
			content:    "example content",
			properties: map[string]interface{}{"key1": make(chan int)},
			metadata:   map[string]interface{}{"meta1": "data1"},
			wantErr:    true,
		},
		{
			name:       "serialization error in metadata",
			content:    "example content",
			properties: map[string]interface{}{"key1": "value1"},
			metadata:   map[string]interface{}{"meta1": make(chan int)},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := api.HashItem(tt.content, tt.properties, tt.metadata)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, got)
			}
		})
	}
}
