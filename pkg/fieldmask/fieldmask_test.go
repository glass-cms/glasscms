package fieldmask_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/glass-cms/glasscms/pkg/fieldmask"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFieldMask(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		query         string
		expectedMasks []string
		expectedError error
	}{
		{
			name:          "missing field mask",
			query:         "",
			expectedMasks: nil,
			expectedError: fieldmask.ErrFieldMaskMissing,
		},
		{
			name:          "empty field mask",
			query:         "fields=",
			expectedMasks: nil,
			expectedError: fieldmask.ErrFieldMaskEmpty,
		},
		{
			name:          "valid field mask",
			query:         "fields=name,age,address.city",
			expectedMasks: []string{"name", "age", "address.city"},
			expectedError: nil,
		},
		{
			name:          "invalid field mask",
			query:         "fields=name,age,address..city",
			expectedMasks: nil,
			expectedError: fmt.Errorf("%w%s", fieldmask.ErrInvalidFieldMask, "address..city"),
		},
		{
			name:          "mixed valid and invalid field masks",
			query:         "fields=name,age,address..city,address.street",
			expectedMasks: nil,
			expectedError: fmt.Errorf("%w%s", fieldmask.ErrInvalidFieldMask, "address..city"),
		},
		{
			name:          "valid wildcard field mask",
			query:         "fields=*,name,age,address.*",
			expectedMasks: []string{"*", "name", "age", "address.*"},
			expectedError: nil,
		},
		{
			name:          "invalid wildcard field mask",
			query:         "fields=*,name,age,address..*",
			expectedMasks: nil,
			expectedError: fmt.Errorf("%w%s", fieldmask.ErrInvalidFieldMask, "address..*"),
		},
		{
			name:          "illegal characters in field mask",
			query:         "fields=name,age,address.city,addre$$.street",
			expectedMasks: nil,
			expectedError: fmt.Errorf("%w%s", fieldmask.ErrInvalidFieldMask, "addre$$.street"),
		},
		{
			name:          "valid mixed wildcard and specific fields",
			query:         "fields=name,age,address.*,address.city",
			expectedMasks: []string{"name", "age", "address.*", "address.city"},
			expectedError: nil,
		},
		{
			name:          "invalid wildcard placement",
			query:         "fields=name,age,address.city.*.street",
			expectedMasks: nil,
			expectedError: fmt.Errorf("%w%s", fieldmask.ErrInvalidFieldMask, "address.city.*.street"),
		},
		{
			name:          "multiple invalid field masks",
			query:         "fields=name,age,address..city,addre$$.street",
			expectedMasks: nil,
			expectedError: fmt.Errorf("%w%s", fieldmask.ErrInvalidFieldMask, "address..city, addre$$.street"),
		},
		{
			name:          "multiple invalid field masks with valid ones",
			query:         "fields=name,age,address..city,addre$$.street,address.city",
			expectedMasks: nil,
			expectedError: fmt.Errorf("%w%s", fieldmask.ErrInvalidFieldMask, "address..city, addre$$.street"),
		},
		{
			name:          "valid nested fields",
			query:         "fields=user.name,user.age,user.address.city",
			expectedMasks: []string{"user.name", "user.age", "user.address.city"},
			expectedError: nil,
		},
		{
			name:          "invalid nested fields",
			query:         "fields=user.name,user..age,user.address.city",
			expectedMasks: nil,
			expectedError: fmt.Errorf("%w%s", fieldmask.ErrInvalidFieldMask, "user..age"),
		},
		{
			name:          "valid and invalid nested fields",
			query:         "fields=user.name,user..age,user.address.city,user.address..street",
			expectedMasks: nil,
			expectedError: fmt.Errorf("%w%s", fieldmask.ErrInvalidFieldMask, "user..age, user.address..street"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/test?%s", tt.query), nil)
			require.NoError(t, err, "failed to create request")

			masks, err := fieldmask.ParseFieldMaskRequest(req)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedError.Error(), "expected error %v, got %v", tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}

			assert.ElementsMatch(t, tt.expectedMasks, masks, "expected masks %v, got %v", tt.expectedMasks, masks)
		})
	}
}
