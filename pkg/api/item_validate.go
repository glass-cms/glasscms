// Code generated by field mask validation generator; DO NOT EDIT.

package api

import "errors"

var ErrInvalidItemField = errors.New("invalid field")

// ValidateItemFieldMask validates a field mask for the Item struct
func ValidateItemFieldMask(fieldmask []string) error {
	validFields := map[string]struct{}{
		"content":      {},
		"create_time":  {},
		"delete_time":  {},
		"display_name": {},
		"hash":         {},
		"metadata":     {},
		"name":         {},
		"properties":   {},
		"update_time":  {},
	}

	for _, field := range fieldmask {
		if _, exists := validFields[field]; !exists {
			return ErrInvalidItemField
		}
	}
	return nil
}
