package slug_test

import (
	"testing"

	"github.com/glass-cms/glasscms/pkg/slug"
)

func TestSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		options  []slug.Option
		expected string
	}{
		{"Hello, World!", nil, "hello-world"},
		{"Hello, World!", []slug.Option{slug.AllowSlashesOption()}, "hello-world"},
		{"Hello/World", []slug.Option{slug.AllowSlashesOption()}, "hello/world"},
		{"Hello/World", nil, "hello-world"},
		{"Hello, World!", []slug.Option{slug.CustomSeparatorOption("_")}, "hello_world"},
		{"Hello/World", []slug.Option{slug.AllowSlashesOption(), slug.CustomSeparatorOption("_")}, "hello/world"},
		{"Café au lait", nil, "cafe-au-lait"},
		{"Über cool", nil, "uber-cool"},
		{"你好，世界", nil, "ni-hao-shi-jie"},
		{"Привет, мир", nil, "privet-mir"},
		{"¡Hola, mundo!", nil, "hola-mundo"},
		{"Hello---World", nil, "hello-world"},
		{"Hello___World", []slug.Option{slug.CustomSeparatorOption("_")}, "hello_world"},
		{"Hello---World", []slug.Option{slug.CustomSeparatorOption("_")}, "hello_world"},
		{"Hello/World/Again", []slug.Option{slug.AllowSlashesOption()}, "hello/world/again"},
		{"Hello/World/Again", []slug.Option{slug.AllowSlashesOption(), slug.CustomSeparatorOption("_")}, "hello/world/again"},
		{"Hello, 世界!", nil, "hello-shi-jie"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			result := slug.Slug(tt.input, tt.options...)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
