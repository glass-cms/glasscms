package v1

import (
	"net/http"
	"reflect"

	v1 "github.com/glass-cms/glasscms/api/v1"
	"github.com/glass-cms/glasscms/server/handler"
)

// ErrorMapper is a function that maps an error to an API error response.
type ErrorMapper func(error) *v1.Error

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

		statusCode, ok := v1.ErrorCodeMapping[errResp.Code]
		if !ok {
			statusCode = http.StatusInternalServerError
		}

		handler.RespondWithJSON(w, statusCode, errResp)
	}

	// TODO: Default error handling.
}
