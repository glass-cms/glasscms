package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	v1 "github.com/glass-cms/glasscms/api/v1"
)

func (s *APIHandler) ItemsCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to read request body: %w", err).Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var request *v1.ItemsCreateJSONRequestBody
	if err = json.Unmarshal(reqBody, &request); err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to unmarshal request body: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.itemService.CreateItem(ctx, request.ToItem())
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to create item: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *APIHandler) ItemsGet(w http.ResponseWriter, _ *http.Request, _ v1.ItemKey) {
	w.WriteHeader(http.StatusTeapot)
}
