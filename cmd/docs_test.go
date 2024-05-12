package cmd_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/glass-cms/glasscms/cmd"
	"github.com/stretchr/testify/require"
)

func TestFilePrepender(t *testing.T) {
	filename := "test_file.md"
	expectedTitle := "Test File"

	result := cmd.DocFilePrepender(filename)

	// Check if the result starts with the YAML front matter
	require.True(t, strings.HasPrefix(result, "---\n"))

	// Check if the title is correct
	require.Contains(t, result, "title: "+expectedTitle)

	// Check if the create timestamp is a valid Unix timestamp
	require.Contains(t, result, "createTime: ")
	timestampStr := strings.Split(strings.Split(result, "createTime: ")[1], "\n")[0]
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	require.NoError(t, err)
	require.Greater(t, timestamp, int64(0))
}
