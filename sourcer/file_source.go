package sourcer

import (
	"os"
	"time"

	"github.com/djherbis/times"
)

var _ Source = (*FileSource)(nil)

type FileSource struct {
	*os.File
	birthtime time.Time
	modtime   time.Time
}

func NewFileSource(file *os.File) (*FileSource, error) {
	stats, err := times.StatFile(file)
	if err != nil {
		return nil, err
	}

	return &FileSource{
		File:      file,
		birthtime: stats.BirthTime(),
		modtime:   stats.ModTime(),
	}, nil
}

func (f FileSource) CreatedAt() time.Time {
	return f.birthtime
}

func (f FileSource) ModifiedAt() time.Time {
	return f.modtime
}
