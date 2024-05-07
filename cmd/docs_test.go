package cmd

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFilePrepender(t *testing.T) {
	filename := "test_file.md"
	expectedTitle := "Test File"

	result := filePrepender(filename)

	// Check if the result starts with the YAML front matter
	require.True(t, strings.HasPrefix(result, "---\n"))

	// Check if the title is correct
	require.Contains(t, result, "title: "+expectedTitle)

	// Check if the create timestamp is a valid Unix timestamp
	require.Contains(t, result, "create_timestamp: ")
	timestampStr := strings.Split(strings.Split(result, "create_timestamp: ")[1], "\n")[0]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	require.NoError(t, err)
	require.True(t, timestamp > 0)
	require.True(t, time.Now().Unix() >= timestamp)
}
