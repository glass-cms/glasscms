package v1

import (
	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/parser"
)

// ToItem converts an api.ItemCreate to an item.Item.
func (i *ItemCreate) ToItem() *item.Item {
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
