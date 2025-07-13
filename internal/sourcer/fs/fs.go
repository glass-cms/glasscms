package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/glass-cms/glasscms/internal/sourcer"
)

var _ sourcer.Sourcer = &FileSystemSourcer{}

// FileSystemSourcer is a DataSourcer that reads files from the file system.
type FileSystemSourcer struct {
	files          []string
	cursor         int // cursor is the index of the next file to be read
	rootPath       string
	ignorePatterns []string
}

// NewSourcer creates a new FileSystemSourcer.
func NewSourcer(rootPath string, ignorePatterns []string) (*FileSystemSourcer, error) {
	var files []string

	absRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(absRootPath, func(path string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories that match ignore patterns
		if dirEntry.IsDir() && shouldIgnoreDirectory(dirEntry.Name(), ignorePatterns) {
			return filepath.SkipDir
		}

		if !dirEntry.IsDir() && strings.HasSuffix(dirEntry.Name(), ".md") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &FileSystemSourcer{
		files:          files,
		rootPath:       absRootPath,
		ignorePatterns: ignorePatterns,
	}, nil
}

func (s *FileSystemSourcer) Next() (sourcer.Source, error) {
	if s.cursor >= len(s.files) {
		return nil, ErrDone
	}

	file, err := os.Open(s.files[s.cursor])
	s.cursor++

	if err != nil {
		return nil, err
	}

	return NewFileSource(file, s.rootPath)
}

func (s *FileSystemSourcer) Remaining() int {
	return s.Size() - s.cursor
}

func (s *FileSystemSourcer) Size() int {
	return len(s.files)
}

// IsValidFileSystemSource checks if the given path exist and is a directory.
func IsValidFileSystemSource(fp string) error {
	fileInfo, err := os.Stat(fp)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidFileSystemSource, err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%w: %s is not a directory", ErrInvalidFileSystemSource, fp)
	}

	return nil
}

// shouldIgnoreDirectory checks if a directory name matches any of the ignore patterns.
func shouldIgnoreDirectory(dirName string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, dirName); matched {
			return true
		}
	}
	return false
}
