package sourcer_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/glass-cms/glasscms/sourcer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testFile struct {
	content string
	depth   int
	pattern string
}

// Helper function to create a temporary directory with files for testing.
func createTempDirWithFiles(files []testFile) (string, error) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		return "", err
	}

	for i, tf := range files {
		dir := tempDir

		for j := range tf.depth {
			dir = filepath.Join(dir, fmt.Sprintf("dir%d", j))
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return "", err
			}
		}

		var file *os.File
		file, err = os.CreateTemp(dir, fmt.Sprintf("file%d", i)+tf.pattern)
		if err != nil {
			return "", err
		}
		defer file.Close()

		_, err = file.WriteString(tf.content)
		if err != nil {
			return "", err
		}
	}

	return tempDir, nil
}

func TestFileSystemSourcer(t *testing.T) {
	t.Parallel()

	// Arrange.
	fileData := []testFile{
		{content: "file 1", depth: 0, pattern: "*.md"},
		{content: "file 2", depth: 0, pattern: "*.md"},
		{content: "file 3", depth: 0, pattern: "*.md"},
	}

	tempDir, err := createTempDirWithFiles(fileData)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	src, err := sourcer.NewFileSystemSourcer(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Walk until the end of the files.
	for i := range src.Size() {
		var data sourcer.Source
		data, err = src.Next()
		if err != nil {
			t.Fatal(err)
		}

		// Read the data from the file.
		var fileContent []byte
		fileContent, err = io.ReadAll(data)
		if err != nil {
			t.Fatal(err)
		}

		// Assert that data is the same as the fileData.
		assert.Equal(t, fileData[i].content, string(fileContent))

		// Assert that the remaining files is correct.
		assert.Equal(t, src.Size()-i-1, src.Remaining())
	}

	// Assert that the sourcer is done.
	_, err = src.Next()
	assert.Equal(t, sourcer.ErrDone, err)
}

func TestFileSystemSourcer_Size(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input []testFile
		want  int
	}{
		"no files": {
			input: []testFile{},
			want:  0,
		},
		"one markdown file": {
			input: []testFile{{content: "file 1", depth: 0, pattern: "*.md"}},
			want:  1,
		},
		"multiple markdown files": {
			input: []testFile{
				{content: "file 1", depth: 0, pattern: "*.md"},
				{content: "file 2", depth: 0, pattern: "*.md"},
				{content: "file 3", depth: 0, pattern: "*.md"},
			},
			want: 3,
		},
		"multiple markdown files with depth": {
			input: []testFile{
				{content: "file 1", depth: 1, pattern: "*.md"},
				{content: "file 2", depth: 2, pattern: "*.md"},
				{content: "file 3", depth: 3, pattern: "*.md"},
			},
			want: 3,
		},
		"non-markdown files": {
			input: []testFile{
				{content: "file 1", depth: 0, pattern: "*.txt"},
				{content: "file 2", depth: 0, pattern: "*.txt"},
				{content: "file 3", depth: 0, pattern: "*.txt"},
			},
			want: 0,
		},
		"mixed files": {
			input: []testFile{
				{content: "file 1", depth: 0, pattern: "*.md"},
				{content: "file 2", depth: 0, pattern: "*.txt"},
				{content: "file 3", depth: 0, pattern: "*.md"},
			},
			want: 2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange.
			tempDir, err := createTempDirWithFiles(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			sourcer, err := sourcer.NewFileSystemSourcer(tempDir)
			if err != nil {
				t.Fatal(err)
			}

			// Asserts.
			assert.Equal(t, tt.want, sourcer.Size())
		})
	}
}

func TestIsValidFileSystemSource(t *testing.T) {
	t.Parallel()

	fp, err := createTempDirWithFiles([]testFile{{content: "file 1", depth: 0, pattern: "*.md"}})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fp)

	require.NoError(t, sourcer.IsValidFileSystemSource(fp))
	require.Error(t, sourcer.IsValidFileSystemSource("non-existent"))
}
