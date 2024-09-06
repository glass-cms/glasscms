package handler

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/glass-cms/glasscms/lib/mediatype"
)

// SerializeResponse writes a JSON response to the response.
func SerializeResponse[T any](w http.ResponseWriter, r *http.Request, statusCode int, data T) {
	acceptHeader := r.Header.Get("Accept")

	switch acceptHeader {
	case mediatype.ApplicationJSON:
		w.Header().Set("Content-Type", mediatype.ApplicationJSON)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	case mediatype.ApplicationXML:
		w.Header().Set("Content-Type", mediatype.ApplicationXML)

		if err := xml.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode XML", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsupported media type", http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(statusCode)
}
