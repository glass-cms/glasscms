package item

import "time"

const (
	PropertyTitle = "title"
)

type Item struct {
	UID string
	// Name is the full resource name of the item.
	// Format: collections/{collection}/items/{item}
	Name        string
	DisplayName string
	Path        string
	Content     string
	Hash        string
	CreateTime  time.Time
	UpdateTime  time.Time
	Properties  map[string]any
}

// Title returns the title property of the item if it exists.
func (i *Item) Title() *string {
	if i.Properties == nil {
		return nil
	}

	if v, ok := i.Properties[PropertyTitle]; ok {
		var s string
		if s, ok = v.(string); ok {
			return &s
		}
	}

	return nil
}
