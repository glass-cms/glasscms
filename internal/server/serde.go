package server

import (
	"encoding/json"
	"net/http"

	"github.com/glass-cms/glasscms/pkg/mediatype"
)

func writeJSONResponse[T any](w http.ResponseWriter, statusCode int, data T) {
	w.Header().Set("Content-Type", mediatype.ApplicationJSON)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
