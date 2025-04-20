package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/glass-cms/glasscms/internal/database"
	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/item/repository/query"
)

var _ item.Repository = &ItemRepository{}

type ItemRepository struct {
	db           *sql.DB
	errorHandler database.ErrorHandler
	queries      query.Queries
}

func NewRepository(db *sql.DB, errorHandler database.ErrorHandler) *ItemRepository {
	return &ItemRepository{
		db:           db,
		errorHandler: errorHandler,
		queries:      *query.New(db),
	}
}

// CreateItem creates a new item in the database.
// If a transaction is provided, the item will be created within the transaction.
//
//nolint:dupl // Very similar to UpsertItem, but with different query parameters.
func (r *ItemRepository) CreateItem(ctx context.Context, tx *sql.Tx, item item.Item) (*item.Item, error) {
	q := r.queries.WithTx(tx)

	propertiesJSON, metadataJSON, err := marshalItemData(item)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	params := query.CreateItemParams{
		Name:        item.Name,
		DisplayName: item.DisplayName,
		CreateTime:  item.CreateTime,
		UpdateTime:  item.UpdateTime,
		DeleteTime: sql.NullTime{
			Time: func() time.Time {
				if item.DeleteTime != nil {
					return *item.DeleteTime
				}
				return time.Time{}
			}(),
			Valid: item.DeleteTime != nil,
		},
		Hash: sql.NullString{
			String: item.Hash,
			Valid:  true,
		},
		Content: sql.NullString{
			String: item.Content,
			Valid:  true,
		},
		Properties: propertiesJSON,
		Metadata:   metadataJSON,
	}

	i, err := q.CreateItem(ctx, params)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	newItem, err := ConvertQueryItem(i)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	return newItem, nil
}

func marshalItemData(item item.Item) ([]byte, []byte, error) {
	propertiesJSON, err := json.Marshal(item.Properties)
	if err != nil {
		return nil, nil, err
	}

	metadataJSON, err := json.Marshal(item.Metadata)
	if err != nil {
		return nil, nil, err
	}

	return propertiesJSON, metadataJSON, nil
}

// GetItem retrieves an item from the database by its resource name.
func (r *ItemRepository) GetItem(ctx context.Context, tx *sql.Tx, name string) (*item.Item, error) {
	q := r.queries.WithTx(tx)

	i, err := q.GetItem(ctx, name)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	foundItem, err := ConvertQueryItem(i)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	return foundItem, nil
}

// UpdateItem updates an existing item in the database.
func (r *ItemRepository) UpdateItem(ctx context.Context, tx *sql.Tx, item item.Item) (*item.Item, error) {
	q := r.queries.WithTx(tx)

	propertiesJSON, metadataJSON, err := marshalItemData(item)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	params := query.UpdateItemParams{
		Name:        item.Name,
		DisplayName: item.DisplayName,
		UpdateTime:  item.UpdateTime,
		Hash: sql.NullString{
			String: item.Hash,
			Valid:  true,
		},
		Content: sql.NullString{
			String: item.Content,
			Valid:  true,
		},
		Properties: propertiesJSON,
		Metadata:   metadataJSON,
		Name_2:     item.Name,
	}

	i, err := q.UpdateItem(ctx, params)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	updatedItem, err := ConvertQueryItem(i)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	return updatedItem, nil
}

// DeleteItem marks an item as deleted in the database.
func (r *ItemRepository) DeleteItem(ctx context.Context, tx *sql.Tx, name string) error {
	q := r.queries.WithTx(tx)

	params := query.DeleteItemParams{
		DeleteTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Name: name,
	}

	err := q.DeleteItem(ctx, params)
	if err != nil {
		return r.errorHandler.HandleError(ctx, err)
	}

	return nil
}

// ListItems retrieves a list of items from the database with optional fieldmask.
func (r *ItemRepository) ListItems(ctx context.Context, tx *sql.Tx, fieldmask []string) ([]*item.Item, error) {
	if len(fieldmask) > 0 {
		return r.listItemsWithFieldmask(ctx, tx, fieldmask)
	}

	items, err := r.queries.WithTx(tx).ListItems(ctx)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	itemList := make([]*item.Item, len(items))
	for index, item := range items {
		convertedItem, convertErr := ConvertQueryItem(item)
		if convertErr != nil {
			return nil, r.errorHandler.HandleError(ctx, convertErr)
		}

		itemList[index] = convertedItem
	}
	return itemList, nil
}

// listItemsWithFieldmask retrieves a list of items from the database with the specified field mask.
// The field mask determines which columns are selected in the query.
func (r *ItemRepository) listItemsWithFieldmask(
	ctx context.Context,
	tx *sql.Tx,
	fieldmask []string,
) ([]*item.Item, error) {
	qry := "SELECT ? FROM items WHERE delete_time IS NULL"
	qry = strings.Replace(qry, "?", strings.Join(fieldmask, ","), 1)

	rows, err := tx.QueryContext(ctx, qry)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	defer rows.Close()

	// FIXME: Scan error on column index 2, name "metadata": unsupported Scan,
	// storing driver.Value type []uint8 into type *map[string]interface {}
	// When selecting properties or metadata in the fieldmask.

	var items []*item.Item
	if err = sqlscan.ScanAll(&items, rows); err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	if err = rows.Err(); err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	return items, nil
}

// UpsertItem creates a new item if it does not exist, otherwise it updates the existing item.
//
//nolint:dupl // Very similar to CreateItem, but with different query parameters.
func (r *ItemRepository) UpsertItem(ctx context.Context, tx *sql.Tx, item item.Item) (*item.Item, error) {
	q := r.queries.WithTx(tx)

	propertiesJSON, metadataJSON, err := marshalItemData(item)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	params := query.UpsertItemParams{
		Name:        item.Name,
		DisplayName: item.DisplayName,
		CreateTime:  item.CreateTime,
		UpdateTime:  item.UpdateTime,
		DeleteTime: sql.NullTime{
			Time: func() time.Time {
				if item.DeleteTime != nil {
					return *item.DeleteTime
				}
				return time.Time{}
			}(),
			Valid: item.DeleteTime != nil,
		},
		Hash: sql.NullString{
			String: item.Hash,
			Valid:  true,
		},
		Content: sql.NullString{
			String: item.Content,
			Valid:  true,
		},
		Properties: propertiesJSON,
		Metadata:   metadataJSON,
	}

	i, err := q.UpsertItem(ctx, params)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	newItem, err := ConvertQueryItem(i)
	if err != nil {
		return nil, r.errorHandler.HandleError(ctx, err)
	}

	return newItem, nil
}

func (r *ItemRepository) DeleteItems(ctx context.Context, tx *sql.Tx, names []string) error {
	return r.queries.WithTx(tx).DeleteItems(ctx, names)
}
