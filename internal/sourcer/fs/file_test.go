package fs_test

import (
	"os"
	"testing"

	"github.com/djherbis/times"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempFile() (*os.File, error) {
	tempFile, err := os.CreateTemp("", "source")
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

func TestNewFileSource(t *testing.T) {
	t.Parallel()

	// Arrange
	tempFile, err := createTempFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Act
	fileSource, err := fs.NewFileSource(tempFile, "")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, fileSource)

	stats, err := times.StatFile(tempFile)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, stats.BirthTime(), fileSource.CreateTime())
	assert.Equal(t, stats.ModTime(), fileSource.UpdateTime())
}
