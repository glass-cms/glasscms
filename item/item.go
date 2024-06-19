package item

import "time"

const (
	PropertyTitle = "title"
)

type Item struct {
	UID string `json:"uid" yaml:"uid"`
	// Name is the full resource name of the item.
	// Format: collections/{collection}/items/{item}
	Name        string         `json:"name" yaml:"name"`
	DisplayName string         `json:"display_name" yaml:"display_name"`
	Path        string         `json:"path" yaml:"path"`
	Content     string         `json:"content" yaml:"content"`
	Hash        string         `json:"hash" yaml:"hash"`
	CreateTime  time.Time      `json:"create_time" yaml:"create_time"`
	UpdateTime  time.Time      `json:"update_time" yaml:"update_time"`
	Properties  map[string]any `json:"properties" yaml:"properties"`
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
