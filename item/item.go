package item

import "time"

const ItemResource = "item"

// Item is the core data structure for the content management system.
// An item represent a single piece of content. It is the structured version of a markdown file.
type Item struct {
	Name        string         `json:"name"`
	DisplayName string         `json:"display_name"`
	Content     string         `json:"content"`
	Hash        string         `json:"hash"`
	CreateTime  time.Time      `json:"create_time"`
	UpdateTime  time.Time      `json:"update_time"`
	DeleteTime  *time.Time     `json:"delete_time"`
	Properties  map[string]any `json:"properties"`
	Metadata    map[string]any `json:"metadata"`
}
