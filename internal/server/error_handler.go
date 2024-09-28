package server

import (
	"net/http"
	"reflect"

	"github.com/glass-cms/glasscms/pkg/api"
)

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
func (h *ErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
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

		SerializeResponse(w, r, statusCode, errResp)
		return
	}

	// Fallback on generic error response if we don't have a specific error mapper.
	errResp := &api.Error{
		Code:    api.ProcessingError,
		Message: "An error occurred while processing the request.",
		Type:    api.ApiError,
	}
	SerializeResponse(w, r, http.StatusInternalServerError, errResp)
}
