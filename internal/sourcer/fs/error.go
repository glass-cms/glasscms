package fs

import "errors"

var (
	ErrInvalidFileSystemSource = errors.New("invalid file system source")

	// ErrDone is returned when there are no items left in the data source.
	ErrDone = errors.New("no items left in the data source")
)
