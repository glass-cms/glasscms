package v1

import (
	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/parser"
)

// MapToDomain converts an api.Item to an item.Item.
func (i *Item) MapToDomain() *item.Item {
	if i == nil {
		return nil
	}
	return &item.Item{
		UID:        i.Id,
		Name:       i.Name,
		Path:       i.Path,
		Content:    i.Content,
		Hash:       parser.HashContent([]byte(i.Content)),
		CreateTime: i.CreateTime,
		UpdateTime: i.UpdateTime,
		Properties: i.Properties,
	}
}
