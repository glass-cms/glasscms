package fieldmask

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const (
	QueryParamFieldMask = "fields"
)

var (
	ErrFieldMaskMissing = errors.New("field mask query param is not set")
	ErrFieldMaskEmpty   = errors.New("field mask is empty")
	ErrInvalidFieldMask = errors.New("field mask is invalid: ")

	FieldMaskRe = regexp.MustCompile(`^(\*|(\w+\.)*\*?|\w+(\.\w+)*\*?)$`)
)

// ParseFieldMask parses and validates the field mask from the HTTP request's query parameters.
// It returns a slice of field masks or an error if the field mask is missing, empty, or invalid.
func ParseFieldMask(r *http.Request) ([]string, error) {
	query := r.URL.Query()

	fields := query.Get(QueryParamFieldMask)
	if fields == "" {
		if !query.Has(QueryParamFieldMask) {
			return nil, ErrFieldMaskMissing
		}
		return nil, ErrFieldMaskEmpty
	}

	fieldMasks := strings.Split(fields, ",")
	invalidFieldMasks := []string{}

	for _, fm := range fieldMasks {
		if !FieldMaskRe.MatchString(fm) {
			invalidFieldMasks = append(invalidFieldMasks, fm)
		}
	}

	if len(invalidFieldMasks) > 0 {
		return nil, fmt.Errorf("%w%s", ErrInvalidFieldMask, strings.Join(invalidFieldMasks, ", "))
	}

	return fieldMasks, nil
}
