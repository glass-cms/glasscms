package item

import "time"

const ItemResource = "item"

// Item is the core data structure for the content management system.
// An item represent a single piece of content. It is the structured version of a markdown file.
type Item struct {
	Name        string         `mapstructure:"name"`
	DisplayName string         `mapstructure:"display_name"`
	Content     string         `mapstructure:"content"`
	Hash        string         `mapstructure:"hash"`
	CreateTime  time.Time      `mapstructure:"create_time"`
	UpdateTime  time.Time      `mapstructure:"update_time"`
	DeleteTime  *time.Time     `mapstructure:"delete_time"`
	Properties  map[string]any `mapstructure:"properties"`
	Metadata    map[string]any `mapstructure:"metadata"`
}

// New creates a new Item instance with the given parameters.
func New(
	name string,
	displayName string,
	content string,
	createTime time.Time,
	updateTime time.Time,
	properties map[string]any,
	metadata map[string]any,
) *Item {
	return &Item{
		Name:        name,
		DisplayName: displayName,
		Content:     content,
		CreateTime:  createTime,
		UpdateTime:  updateTime,
		Properties:  properties,
		Metadata:    metadata,
	}
}
