package item

import "time"

const ItemResource = "item"

// Item is the core data structure for the content management system.
// An item represent a single piece of content. It is the structured version of a markdown file.
type Item struct {
	Name        string
	DisplayName string
	Content     string
	Hash        string
	CreateTime  time.Time
	UpdateTime  time.Time
	DeleteTime  *time.Time
	Properties  map[string]any
	Metadata    map[string]any
}
