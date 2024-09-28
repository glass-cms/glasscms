package sourcer

import (
	"io"
	"time"
)

// Source is a data source that can be read from.
type Source interface {
	io.ReadCloser

	// Name returns the name of the source.
	Name() string

	// CreateTime returns the time when the source was created.
	CreateTime() time.Time

	// UpdatTime returns the time when the source was last modified.
	UpdateTime() time.Time
}
