package item

import "time"

type Item struct {
	Name       string         `json:"name"`
	Content    string         `json:"content"`
	CreateTime time.Time      `json:"createTime"`
	UpdateTime time.Time      `json:"updateTime"`
	Properties map[string]any `json:"properties"`
}
