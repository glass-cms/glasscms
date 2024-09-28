package sourcer

import (
	"io"
	"time"
)

// NilSource is a no-op sentinel Source.
var NilSource = nilSource{
	ReadCloser: io.NopCloser(nil),
}

type nilSource struct {
	io.ReadCloser
}

func (nilSource) Name() string {
	return ""
}

func (nilSource) CreateTime() time.Time {
	return time.Time{}
}

func (nilSource) UpdateTime() time.Time {
	return time.Time{}
}
