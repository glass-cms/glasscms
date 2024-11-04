package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/glass-cms/glasscms/pkg/mediatype"
)

// SerializeJSONResponse serializes the given data to JSON and writes it to the HTTP response.
// It sets the Content-Type header to "application/json" and the response status code to the provided statusCode.
func SerializeJSONResponse[T any](w http.ResponseWriter, statusCode int, data T) {
	w.Header().Set("Content-Type", mediatype.ApplicationJSON)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// DeserializeJSONRequestBody reads the JSON-encoded request body from an HTTP request
// and deserializes it into a value of type T.
func DeserializeJSONRequestBody[T any](r *http.Request) (*T, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var requestBody T
	if err = json.Unmarshal(body, &requestBody); err != nil {
		return nil, NewDeserializeError(err)
	}

	return &requestBody, nil
}

// DeserializeError is an error type that wraps an error that occurred during deserialization.
type DeserializeError struct {
	Err error
}

func (e *DeserializeError) Error() string {
	return e.Err.Error()
}

func NewDeserializeError(err error) *DeserializeError {
	return &DeserializeError{Err: err}
}
