package server

import (
	"fmt"
	"net/http"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/pkg/api"
)

// ItemsCreate creates a new item.
func (s *Server) ItemsCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	createRequest, err := DeserializeJSONRequestBody[api.ItemsCreateJSONRequestBody](r)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to read request body: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	var createdItem item.Item
	createdItem, err = s.itemService.CreateItem(ctx, itemCreateToItem(createRequest))
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to create item: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	SerializeJSONResponse(w, http.StatusCreated, createdItem)
}

// ItemsGet retrieves an item by name.
func (s *Server) ItemsGet(w http.ResponseWriter, r *http.Request, name string) {
	ctx := r.Context()
	s.logger.DebugContext(ctx, fmt.Sprintf("getting item: %s", name))

	item, err := s.itemService.GetItem(ctx, name)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to get item: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	SerializeJSONResponse(w, http.StatusOK, item)
}

// ItemsUpdate updates an item by name.
func (s *Server) ItemsUpdate(w http.ResponseWriter, _ *http.Request, _ string) {
	SerializeJSONResponse[any](w, http.StatusNotImplemented, nil)
}

func (s *Server) ItemsList(w http.ResponseWriter, _ *http.Request, _ api.ItemsListParams) {
	SerializeJSONResponse[any](w, http.StatusNotImplemented, nil)
}

func itemCreateToItem(i *api.ItemCreate) item.Item {
	return item.Item{
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
