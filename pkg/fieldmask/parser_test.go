package fieldmask_test

import (
	"reflect"
	"testing"

	"github.com/glass-cms/glasscms/pkg/fieldmask"
)

func TestFieldMasksForType(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Test  string // Should not be included.
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email,omitempty"`
		Phone string `json:"-"`
	}

	tests := []struct {
		name string
		item interface{}
		want map[string]struct{}
	}{
		{
			name: "basic struct",
			item: TestStruct{},
			want: map[string]struct{}{
				"name":  {},
				"age":   {},
				"email": {},
			},
		},
		{
			name: "struct with no json tags",
			item: struct {
				Field1 string
				Field2 int
			}{},
			want: map[string]struct{}{},
		},
		{
			name: "struct with all fields omitted",
			item: struct {
				Field1 string `json:"-"`
				Field2 int    `json:"-"`
			}{},
			want: map[string]struct{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := fieldmask.FieldMasksForType(tt.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FieldMasksForType() = %v, want %v", got, tt.want)
			}
		})
	}
}
