package server

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/resource"
)

var ErrorCodeMapping = map[api.ErrorCode]int{
	api.ParameterInvalid:      http.StatusBadRequest,
	api.ParameterMissing:      http.StatusBadRequest,
	api.ProcessingError:       http.StatusInternalServerError,
	api.ResourceAlreadyExists: http.StatusConflict,
	api.ResourceMissing:       http.StatusNotFound,
}

// ErrorMapper is a function that maps an error to an API error response.
type ErrorMapper func(error) *api.Error

type ErrorHandler struct {
	Mappers map[reflect.Type]ErrorMapper
}

// NewErrorHandler returns a new instance of ErrorHandler.
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		Mappers: make(map[reflect.Type]ErrorMapper),
	}
}

// RegisterErrorMapper registers an error mapper for a specific error type.
func (h *ErrorHandler) RegisterErrorMapper(errType reflect.Type, mapper ErrorMapper) {
	h.Mappers[errType] = mapper
}

// HandleError handles an error by writing an appropriate response to the client.
func (h *ErrorHandler) HandleError(w http.ResponseWriter, _ *http.Request, err error) {
	if err == nil {
		return
	}

	errType := reflect.TypeOf(err)
	if mapper, exists := h.Mappers[errType]; exists {
		errResp := mapper(err)

		statusCode, ok := ErrorCodeMapping[errResp.Code]
		if !ok {
			statusCode = http.StatusInternalServerError
		}

		SerializeJSONResponse(w, statusCode, errResp)
		return
	}

	// Fallback on generic error response if we don't have a specific error mapper.
	errResp := &api.Error{
		Code:    api.ProcessingError,
		Message: "An error occurred while processing the request.",
		Type:    api.ApiError,
	}
	SerializeJSONResponse(w, http.StatusInternalServerError, errResp)
}

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

// ErrorMapperNotFoundError maps a resource.NotFoundError to an API error response.
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
