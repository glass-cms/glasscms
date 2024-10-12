package repository

import (
	"database/sql"
	"time"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/internal/item/repository/query"
)

func ConvertQueryItem(dbItem query.Item) item.Item {
	return item.Item{
		Name:        dbItem.Name,
		DisplayName: dbItem.DisplayName,
		Content:     convertNullString(dbItem.Content),
		Hash:        convertNullString(dbItem.Hash),
		CreateTime:  dbItem.CreateTime,
		UpdateTime:  dbItem.UpdateTime,
		DeleteTime:  convertNullTime(dbItem.DeleteTime),
		Properties:  convertInterfaceToMap(dbItem.Properties),
		Metadata:    convertInterfaceToMap(dbItem.Metadata),
	}
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

func convertInterfaceToMap(i interface{}) map[string]any {
	if m, ok := i.(map[string]any); ok {
		return m
	}
	return nil
}
