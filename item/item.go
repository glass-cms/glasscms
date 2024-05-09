package item

import "time"

type Item struct {
	Name       string         `json:"name"`
	Content    string         `json:"content"`
	CreateTime time.Time      `json:"create_time"`
	UpdateTime time.Time      `json:"update_time"`
	Properties map[string]any `json:"properties"`
}
