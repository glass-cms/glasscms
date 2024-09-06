package v1

import "net/http"

var ErrorCodeMapping = map[ErrorCode]int{
	ParameterInvalid:      http.StatusBadRequest,
	ParameterMissing:      http.StatusBadRequest,
	ProcessingError:       http.StatusInternalServerError,
	ResourceAlreadyExists: http.StatusConflict,
	ResourceMissing:       http.StatusNotFound,
}
