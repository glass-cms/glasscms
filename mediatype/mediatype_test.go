// Package mediatype provides types, constant and parsing functions for media types (MIME).
package mediatype_test

import (
	"reflect"
	"testing"

	"github.com/glass-cms/glasscms/mediatype"
)

func stringPtr(s string) *string {
	return &s
}

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    *mediatype.MediaType
		wantErr bool
	}{
		{
			name:  "application/json",
			input: "application/json",
			want: &mediatype.MediaType{
				MediaType:  "application/json",
				Type:       "application",
				Subtype:    stringPtr("json"),
				Parameters: map[string]string{},
			},
			wantErr: false,
		},
		{
			name:  "application/json with charset",
			input: "application/json; charset=utf-8",
			want: &mediatype.MediaType{
				MediaType:  "application/json",
				Type:       "application",
				Subtype:    stringPtr("json"),
				Parameters: map[string]string{"charset": "utf-8"},
			},
			wantErr: false,
		},
		{
			name:  "without subtype",
			input: "text",
			want: &mediatype.MediaType{
				MediaType:  "text",
				Type:       "text",
				Subtype:    nil,
				Parameters: map[string]string{},
			},
			wantErr: false,
		},
		{
			name:    "invalid media type with slash",
			input:   "text/",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := mediatype.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
