// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package query

import (
	"time"
)

type Token struct {
	ID         string    `db:"id"`
	Suffix     string    `db:"suffix"`
	Hash       string    `db:"hash"`
	CreateTime time.Time `db:"create_time"`
	ExpireTime time.Time `db:"expire_time"`
}
