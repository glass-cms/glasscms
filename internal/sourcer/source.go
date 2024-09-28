package sourcer

import (
	"io"
	"time"
)

// Source is a data source that can be read from.
type Source interface {
	io.ReadCloser
	Name() string

	CreatedAt() time.Time
	ModifiedAt() time.Time

	// Add a metadata command.
}
