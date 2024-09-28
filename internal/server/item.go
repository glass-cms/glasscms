package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/pkg/api"
)

func (s *Server) ItemsCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to read request body: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	var request *api.ItemsCreateJSONRequestBody
	if err = json.Unmarshal(reqBody, &request); err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to unmarshal request body: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	err = s.itemService.CreateItem(ctx, ToItem(request))
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to create item: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	// TODO: Write response.

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) ItemsGet(w http.ResponseWriter, r *http.Request, name api.ItemKey) {
	ctx := r.Context()

	item, err := s.itemService.GetItem(ctx, name)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to get item: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	SerializeResponse(w, r, http.StatusOK, item)
}

// ToItem converts an api.ItemCreate to an item.Item.
func ToItem(i *api.ItemCreate) *item.Item {
	if i == nil {
		return nil
	}

	return &item.Item{
		Name:        i.Name,
		DisplayName: i.DisplayName,
		Content:     i.Content,
		Hash:        parser.HashContent([]byte(i.Content)),
		CreateTime:  i.CreateTime,
		UpdateTime:  i.UpdateTime,
		Properties:  i.Properties,
		Metadata:    i.Metadata,
	}
}

func FromItem(item *item.Item) *api.Item {
	if item == nil {
		return nil
	}

	return &api.Item{
		Name:        item.Name,
		DisplayName: item.DisplayName,
		Content:     item.Content,
		CreateTime:  item.CreateTime,
		UpdateTime:  item.UpdateTime,
		Properties:  item.Properties,
		Metadata:    item.Metadata,
	}
}
