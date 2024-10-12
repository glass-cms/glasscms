package log_test

import (
	"testing"

	"github.com/glass-cms/glasscms/pkg/log"
)

func TestFormat_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    log.Format
		wantErr bool
	}{
		{"valid TEXT", []byte("TEXT"), log.FormatText, false},
		{"valid text", []byte("text"), log.FormatText, false},
		{"valid JSON", []byte("JSON"), log.FormatJSON, false},
		{"valid json", []byte("json"), log.FormatJSON, false},
		{"invalid type", []byte("xml"), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var lt log.Format
			err := lt.UnmarshalText(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if lt != tt.want {
				t.Errorf("UnmarshalText() = %v, want %v", lt, tt.want)
			}
		})
	}
}
