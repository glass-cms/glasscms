package item

import "time"

const (
	PropertyTitle = "title"
)

type Item struct {
	UID string `json:"uid" yaml:"uid"`
	// Name is the full resource name of the item.
	Name       string         `json:"name" yaml:"name"`
	Path       string         `json:"path" yaml:"path"`
	Content    string         `json:"content" yaml:"content"`
	Hash       string         `json:"hash" yaml:"hash"`
	CreateTime time.Time      `json:"create_time" yaml:"create_time"`
	UpdateTime time.Time      `json:"update_time" yaml:"update_time"`
	Properties map[string]any `json:"properties" yaml:"properties"`
}
