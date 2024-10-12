// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package query

import (
	"database/sql"
	"time"
)

type Item struct {
	Name        string
	DisplayName string
	CreateTime  time.Time
	UpdateTime  time.Time
	DeleteTime  sql.NullTime
	Hash        sql.NullString
	Content     sql.NullString
	Properties  interface{}
	Metadata    interface{}
}
