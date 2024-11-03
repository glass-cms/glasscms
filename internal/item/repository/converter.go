package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/item/repository/query"
)

func ConvertQueryItem(dbItem query.Item) (*item.Item, error) {
	var domainItem item.Item
	var err error

	domainItem = item.Item{
		Name:        dbItem.Name,
		DisplayName: dbItem.DisplayName,
		Content:     convertNullString(dbItem.Content),
		Hash:        convertNullString(dbItem.Hash),
		CreateTime:  dbItem.CreateTime,
		UpdateTime:  dbItem.UpdateTime,
		DeleteTime:  convertNullTime(dbItem.DeleteTime),
	}

	domainItem.Properties, err = unmarshalJSONToMap(dbItem.Properties)
	if err != nil {
		return nil, errors.New("failed to unmarshal Properties: " + err.Error())
	}

	domainItem.Metadata, err = unmarshalJSONToMap(dbItem.Metadata)
	if err != nil {
		return nil, errors.New("failed to unmarshal Metadata: " + err.Error())
	}

	return &domainItem, nil
}

func convertNullString(nullStr sql.NullString) string {
	if nullStr.Valid {
		return nullStr.String
	}
	return ""
}

func convertNullTime(nullTime sql.NullTime) *time.Time {
	if nullTime.Valid {
		return &nullTime.Time
	}
	return nil
}

func unmarshalJSONToMap(data interface{}) (map[string]any, error) {
	if data == nil {
		return make(map[string]any), nil
	}

	var result map[string]any

	switch v := data.(type) {
	case []byte:
		if err := json.Unmarshal(v, &result); err != nil {
			return nil, err
		}
	case string:
		if err := json.Unmarshal([]byte(v), &result); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown data type for JSON unmarshal")
	}

	return result, nil
}
