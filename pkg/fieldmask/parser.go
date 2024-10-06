package fieldmask

import (
	"reflect"
	"strings"
)

// FieldMasksForType returns a map of JSON field names for the given struct type.
func FieldMasksForType(item interface{}) map[string]struct{} {
	fields := make(map[string]struct{})

	val := reflect.ValueOf(item)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		jsonTag := field.Tag.Get("json")

		// If there's a comma (like `json:"name,omitempty"`), strip out the options
		jsonField := strings.Split(jsonTag, ",")[0]

		// Skip fields with json tag "-"
		if jsonField != "" && jsonField != "-" {
			fields[jsonField] = struct{}{}
		}
	}

	return fields
}
