package fs

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djherbis/times"
	"github.com/glass-cms/glasscms/internal/sourcer"
)

var _ sourcer.Source = &FileSource{}

type FileSource struct {
	*os.File

	birthtime time.Time
	modtime   time.Time

	rootPath string
}

func NewFileSource(file *os.File, rootPath string) (*FileSource, error) {
	stats, err := times.StatFile(file)
	if err != nil {
		return nil, err
	}

	return &FileSource{
		File:      file,
		birthtime: stats.BirthTime(),
		modtime:   stats.ModTime(),
		rootPath:  rootPath,
	}, nil
}

// Name returns the name of the file relative to the root path.
func (f *FileSource) Name() string {
	fn := f.File.Name()

	name, err := filepath.Rel(f.rootPath, fn)
	if err != nil {
		panic(err)
	}

	return strings.TrimSuffix(name, filepath.Ext(name))
}

func (f *FileSource) CreateTime() time.Time {
	return f.birthtime
}

func (f *FileSource) UpdateTime() time.Time {
	return f.modtime
}
