package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/resource"
)

// ErrorMapperAlreadyExistsError maps a resource.AlreadyExistsError to an API error response.
func ErrorMapperAlreadyExistsError(err error) *api.Error {
	var alreadyExistsErr *resource.AlreadyExistsError
	if !errors.As(err, &alreadyExistsErr) {
		panic("error is not a resource.AlreadyExistsError")
	}

	return &api.Error{
		Code:    api.ResourceAlreadyExists,
		Message: fmt.Sprintf("An %s with the name already exists", alreadyExistsErr.Resource),
		Type:    api.ApiError,
		Details: map[string]interface{}{
			"resource": alreadyExistsErr.Resource,
			"name":     alreadyExistsErr.Name,
		},
	}
}

func ErrorMapperNotFoundError(err error) *api.Error {
	var notFoundErr *resource.NotFoundError
	if !errors.As(err, &notFoundErr) {
		panic("error is not a resource.NotFoundError")
	}

	return &api.Error{
		Code:    api.ResourceMissing,
		Message: fmt.Sprintf("The %s with the name was not found", notFoundErr.Resource),
		Type:    api.ApiError,
		Details: map[string]interface{}{
			"resource": notFoundErr.Resource,
			"name":     notFoundErr.Name,
		},
	}
}

var ErrorCodeMapping = map[api.ErrorCode]int{
	api.ParameterInvalid:      http.StatusBadRequest,
	api.ParameterMissing:      http.StatusBadRequest,
	api.ProcessingError:       http.StatusInternalServerError,
	api.ResourceAlreadyExists: http.StatusConflict,
	api.ResourceMissing:       http.StatusNotFound,
}
