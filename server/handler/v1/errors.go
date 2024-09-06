package v1

import (
	"errors"
	"fmt"

	v1 "github.com/glass-cms/glasscms/api/v1"
	"github.com/glass-cms/glasscms/lib/resource"
)

// ErrorMapperAlreadyExistsError maps a resource.AlreadyExistsError to an API error response.
func ErrorMapperAlreadyExistsError(err error) *v1.Error {
	var alreadyExistsErr *resource.AlreadyExistsError
	if !errors.As(err, &alreadyExistsErr) {
		panic("error is not a resource.AlreadyExistsError")
	}

	return &v1.Error{
		Code:    v1.ResourceAlreadyExists,
		Message: fmt.Sprintf("An %s with the name already exists", alreadyExistsErr.Resource),
		Type:    v1.ApiError,
		Details: map[string]interface{}{
			"resource": alreadyExistsErr.Resource,
			"name":     alreadyExistsErr.Name,
		},
	}
}
