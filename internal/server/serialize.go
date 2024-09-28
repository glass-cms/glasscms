package server

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/glass-cms/glasscms/pkg/mediatype"
)

// SerializeResponse writes the provided data to the response writer in the
// appropriate media type based on the request's Accept header.
func SerializeResponse[T any](w http.ResponseWriter, r *http.Request, statusCode int, data T) {
	w.WriteHeader(statusCode)

	switch r.Header.Get("Accept") {
	case mediatype.ApplicationJSON:
		w.Header().Set("Content-Type", mediatype.ApplicationJSON)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	case mediatype.ApplicationXML:
		w.Header().Set("Content-Type", mediatype.ApplicationXML)

		if err := xml.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode XML", http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Unsupported media type", http.StatusNotAcceptable)
	}
}
