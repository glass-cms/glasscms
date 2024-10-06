package fieldmask

import (
	"reflect"
	"strings"
)

// FieldMasksForType returns a map of JSON field names for the given struct type.
func FieldMasksForType(item any, parent string) map[string]struct{} {
	fields := make(map[string]struct{})

	val := reflect.ValueOf(item)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)

		jsonTag := field.Tag.Get("json")

		// If there's a comma (like `json:"name,omitempty"`), strip out the options
		jsonField := strings.Split(jsonTag, ",")[0]

		// Skip fields with json tag "-"
		if jsonField == "" || jsonField == "-" {
			continue
		}

		fieldType := field.Type
		kind := fieldType.Kind()
		if kind == reflect.Struct {
			// Recursively add fields from nested structs
			for nestedField := range FieldMasksForType(reflect.New(fieldType).Elem().Interface(), jsonField) {
				fields[nestedField] = struct{}{}
			}
		} else {
			// If the field has a parent, add it to the field name
			if parent != "" {
				jsonField = parent + "." + jsonField
			}
			fields[jsonField] = struct{}{}
		}
	}

	return fields
}
