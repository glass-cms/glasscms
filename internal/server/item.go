package server

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/fieldmask"
)

// TODO: Add option to parse wikilinks in the content from the API.

// ItemsCreate creates a new item.
func (s *Server) ItemsCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	createRequest, err := DeserializeJSONRequestBody[api.ItemsCreateJSONRequestBody](r)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to read request body: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	itemToCreate, err := itemCreateToItem(createRequest)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to convert item create to item: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	var createdItem *item.Item
	createdItem, err = s.itemService.CreateItem(ctx, itemToCreate)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to create item: %w", err).Error())

		s.errorHandler.HandleError(w, r, err)
		return
	}

	SerializeJSONResponse(w, http.StatusCreated, FromItem(createdItem))
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

	SerializeJSONResponse(w, http.StatusOK, FromItem(item))
}

// ItemsUpdate updates an item by name.
func (s *Server) ItemsUpdate(w http.ResponseWriter, _ *http.Request, _ string) {
	// TODO: Implement item update.
	SerializeJSONResponse[any](w, http.StatusNotImplemented, nil)
}

func (s *Server) ItemsUpsert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	upsertRequest, err := DeserializeJSONRequestBody[api.ItemsUpsertJSONRequestBody](r)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to read request body: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	items := make([]item.Item, len(*upsertRequest))
	for i, itemUpsert := range *upsertRequest {
		var item item.Item
		item, err = itemUpsertToItem(&itemUpsert)
		if err != nil {
			s.logger.ErrorContext(ctx, fmt.Errorf("failed to convert item upsert to item: %w", err).Error())
			s.errorHandler.HandleError(w, r, err)
			return
		}

		items[i] = item
	}

	upsertedItems, err := s.itemService.UpsertItems(ctx, items)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to upsert items: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	apiItems := make([]*api.Item, len(upsertedItems))
	for i, item := range upsertedItems {
		apiItems[i] = FromItem(item)
	}

	SerializeJSONResponse(w, http.StatusOK, apiItems)
}

func (s *Server) ItemsList(w http.ResponseWriter, r *http.Request, params api.ItemsListParams) {
	ctx := r.Context()
	s.logger.DebugContext(ctx, "listing items")

	fm := getFieldmask(params.Fields)
	if err := api.ValidateItemFieldMask(fm); err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to validate field mask: %w", err).Error())
		s.errorHandler.HandleError(w, r, fieldmask.NewInvalidFieldMaskError(fm))
		return
	}

	items, err := s.itemService.ListItems(ctx, fm)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to list items: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	// Convert items to API items.
	var apiItems = make([]*api.Item, len(items))
	for i, item := range items {
		apiItems[i] = FromItem(item)
	}

	if len(fm) == 0 || fm == nil {
		SerializeJSONResponse(w, http.StatusOK, apiItems)
		return
	}

	SerializeJSONResponse(w, http.StatusOK, applyItemFieldMask(apiItems, fm))
}

func (s *Server) ItemsDeleteMany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s.logger.DebugContext(ctx, "deleting items")

	deleteRequest, err := DeserializeJSONRequestBody[api.ItemsDeleteManyJSONRequestBody](r)
	if err != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to read request body: %w", err).Error())
		s.errorHandler.HandleError(w, r, err)
		return
	}

	// TODO: Add valdidation for empty slice of request.

	if deleteErr := s.itemService.DeleteItems(ctx, deleteRequest.Names); deleteErr != nil {
		s.logger.ErrorContext(ctx, fmt.Errorf("failed to delete items: %w", deleteErr).Error())
		s.errorHandler.HandleError(w, r, deleteErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getFieldmask(r *[]string) []string {
	if r == nil {
		return make([]string, 0)
	}
	return *r
}

func applyItemFieldMask(items []*api.Item, fieldmask []string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(items))

	for i, item := range items {
		maskedItem := make(map[string]interface{})
		itemMap := itemToMap(item)
		for _, field := range fieldmask {
			if value, ok := itemMap[field]; ok {
				maskedItem[field] = value
			}
		}
		result[i] = maskedItem
	}

	return result
}

func itemToMap(item *api.Item) map[string]interface{} {
	itemMap := make(map[string]interface{})
	itemValue := reflect.ValueOf(item).Elem()
	itemType := itemValue.Type()

	for i := 0; i < itemType.NumField(); i++ {
		field := itemType.Field(i)
		fieldValue := itemValue.Field(i).Interface()

		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			tagParts := strings.Split(jsonTag, ",")
			fieldName := tagParts[0]
			itemMap[fieldName] = fieldValue
		}
	}

	return itemMap
}

func itemCreateToItem(i *api.ItemCreate) (item.Item, error) {
	hash, err := api.HashItem(i.Content, i.Properties, i.Metadata)
	if err != nil {
		return item.Item{}, err
	}

	return item.Item{
		Name:        i.Name,
		DisplayName: i.DisplayName,
		Content:     i.Content,
		Hash:        hash,
		CreateTime:  i.CreateTime,
		UpdateTime:  i.UpdateTime,
		Properties:  i.Properties,
		Metadata:    i.Metadata,
		DeleteTime:  nil,
	}, nil
}

func itemUpsertToItem(i *api.ItemUpsert) (item.Item, error) {
	hash, err := api.HashItem(i.Content, i.Properties, i.Metadata)
	if err != nil {
		return item.Item{}, err
	}

	return item.Item{
		Name:        i.Name,
		DisplayName: i.DisplayName,
		Content:     i.Content,
		Hash:        hash,
		CreateTime:  i.CreateTime,
		UpdateTime:  i.UpdateTime,
		Properties:  i.Properties,
		Metadata:    i.Metadata,
		DeleteTime:  i.DeleteTime,
	}, nil
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
		Hash:        &item.Hash,
		DeleteTime:  item.DeleteTime,
	}
}
