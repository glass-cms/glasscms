// Package sourcer provides abstractions for accessing various data repositories
// in an iterative manner. It defines interfaces for data sources and data sourcers
// which can be implemented to read and iterate over data from different repositories.
package sourcer

import (
	"io"
	"time"
)

//go:generate moq -out mock_sourcer.go . Sourcer

// Sourcer is an iterator that provides data to be parsed.
type Sourcer interface {
	// Next returns the next piece of data to be parsed.
	Next() (Source, error)

	// Remaining returns the number of pieces of data remaining to be parsed.
	Remaining() int

	// Size returns the total number of pieces of data to be parsed.
	Size() int
}

// Source is a data source that can be read from.
type Source interface {
	io.ReadCloser

	// Name returns the name of the source.
	Name() string

	// CreateTime returns the time when the source was created.
	CreateTime() time.Time

	// UpdateTime returns the time when the source was last modified.
	UpdateTime() time.Time
}
