// Package mediatype provides types, constant and parsing functions for media types (MIME).
package mediatype

import (
	"mime"
	"strings"
)

// MediaType represents a media type (MIME).
//
// A MIME type most commonly consists of just two parts: a type and a subtype,
// separated by a slash (/) — with no whitespace between:.
type MediaType struct {
	MediaType string

	Type    string
	Subtype *string

	// Parameters are optional key-value pairs that follow the type/subtype in a MIME type.
	Parameters map[string]string
}

// Parse parses a media type string and returns a MediaType.
// If the string is not a valid media type, an error is returned.
func Parse(s string) (*MediaType, error) {
	mt, params, err := mime.ParseMediaType(s)
	if err != nil {
		return nil, err
	}

	split := strings.Split(mt, "/")
	var subtype *string

	// If there is a subtype, set it.
	if len(split) == 2 {
		subtype = &split[1]
	}

	return &MediaType{
		MediaType:  mt,
		Type:       split[0],
		Subtype:    subtype,
		Parameters: params,
	}, nil
}
