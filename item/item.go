package item

import "time"

type Item struct {
	Name       string         `json:"name" yaml:"name"`
	Content    string         `json:"content" yaml:"content"`
	CreateTime time.Time      `json:"createTime" yaml:"createTime"`
	UpdateTime time.Time      `json:"updateTime" yaml:"updateTime"`
	Properties map[string]any `json:"properties" yaml:"properties"`
}
