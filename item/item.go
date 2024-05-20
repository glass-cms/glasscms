package item

import "time"

const (
	// PropertyTitle is the key for the title property.
	PropertyTitle = "title"
)

type Item struct {
	Name       string         `json:"name" yaml:"name"`
	Path       string         `json:"path" yaml:"path"`
	Content    string         `json:"content" yaml:"content"`
	CreateTime time.Time      `json:"createTime" yaml:"createTime"`
	UpdateTime time.Time      `json:"updateTime" yaml:"updateTime"`
	Properties map[string]any `json:"properties" yaml:"properties"`
}

// Title returns the title property of the item if it exists.
func (i *Item) Title() *string {
	if i.Properties == nil {
		return nil
	}

	if v, ok := i.Properties[PropertyTitle]; ok {
		if s, ok := v.(string); ok {
			return &s
		}
	}

	return nil
}
