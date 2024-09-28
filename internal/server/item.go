package server

import (
	"fmt"
	"net/http"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/pkg/api"
)

func (s *Server) ItemsCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	createRequest, err := DeserializeJSONRequestBody[api.ItemsCreateJSONRequestBody](r)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to read request body: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	item := itemCreateToItem(createRequest)
	err = s.itemService.CreateItem(ctx, item)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to create item: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	SerializeJSONResponse(w, http.StatusCreated, item)
}

func (s *Server) ItemsGet(w http.ResponseWriter, r *http.Request, name api.ItemKey) {
	ctx := r.Context()

	item, err := s.itemService.GetItem(ctx, name)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to get item: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	SerializeJSONResponse(w, http.StatusOK, item)
}

func itemCreateToItem(i *api.ItemCreate) *item.Item {
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
