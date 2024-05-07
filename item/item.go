package item

import "time"

type Item struct {
	Title      string
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
	Properties map[string]any
}
