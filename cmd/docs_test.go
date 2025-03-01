package cmd_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/glass-cms/glasscms/cmd"
	"github.com/stretchr/testify/assert"
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
	require.Positive(t, timestamp, int64(0))
}

func TestDocsCommand(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		slugify     bool
		prefix      string
		linkPattern string
	}{
		{
			name:        "default behavior",
			args:        []string{},
			slugify:     false,
			prefix:      "",
			linkPattern: `\[.*?\]\([^)]*\.md\)`,
		},
		{
			name:        "with link prefix",
			args:        []string{"--link-prefix", "../commands/"},
			slugify:     false,
			prefix:      "../commands/",
			linkPattern: `\[.*?\]\(\.\.\/commands\/[^)]*\.md\)`,
		},
		{
			name:        "with slugify enabled",
			args:        []string{"--slugify-links"},
			slugify:     true,
			prefix:      "",
			linkPattern: `\[.*?\]\([^)]*[^\.md]\)`,
		},
		{
			name:        "with slugify and link prefix",
			args:        []string{"--slugify-links", "--link-prefix", "../commands/"},
			slugify:     true,
			prefix:      "../commands/",
			linkPattern: `\[.*?\]\(\.\.\/commands\/[^)]*[^\.md]\)`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "docs-test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tempDir)

			command := cmd.NewDocsCommand().Command
			args := append([]string{"--output-dir", tempDir}, tc.args...)
			command.SetArgs(args)
			err = command.Execute()
			require.NoError(t, err)

			require.DirExists(t, tempDir)

			files, err := os.ReadDir(tempDir)
			require.NoError(t, err)

			var mdFiles []string
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".md") {
					mdFiles = append(mdFiles, file.Name())
				}
			}
			require.NotEmpty(t, mdFiles, "No markdown files found in output directory")

			foundCorrectLinkFormat := false
			for _, fileName := range mdFiles {
				filePath := filepath.Join(tempDir, fileName)
				fileContent, readErr := os.ReadFile(filePath)
				require.NoError(t, readErr)

				if checkLinkFormat(t, fileContent, tc.linkPattern) {
					foundCorrectLinkFormat = true
					break
				}
			}

			assert.True(t, foundCorrectLinkFormat, "Expected to find links matching pattern %s", tc.linkPattern)

			for _, fileName := range mdFiles {
				filePath := filepath.Join(tempDir, fileName)
				fileContent, readErr := os.ReadFile(filePath)
				require.NoError(t, readErr)

				contentStr := string(fileContent)

				validateLinks(t, contentStr, tc.slugify, tc.prefix)
			}
		})
	}
}

func checkLinkFormat(_ *testing.T, content []byte, pattern string) bool {
	linkRegex := regexp.MustCompile(pattern)
	return linkRegex.Match(content)
}

func validateLinks(t *testing.T, contentStr string, slugify bool, prefix string) {
	linkMatches := regexp.MustCompile(`\[.*?\]\((.*?)\)`).FindAllStringSubmatch(contentStr, -1)
	for _, match := range linkMatches {
		if len(match) <= 1 {
			continue
		}

		linkTarget := match[1]

		if strings.HasPrefix(linkTarget, "http") || strings.HasPrefix(linkTarget, "#") {
			continue
		}

		if slugify {
			assert.NotContains(t, linkTarget, ".md",
				"With slugify enabled, links should not contain .md extension: %s", linkTarget)
		} else if !strings.HasPrefix(linkTarget, "http") && !strings.HasPrefix(linkTarget, "#") {
			assert.Contains(t, linkTarget, ".md",
				"Without slugify, links should contain .md extension: %s", linkTarget)
		}

		if prefix != "" {
			assert.True(t, strings.HasPrefix(linkTarget, prefix),
				"With prefix %s, links should start with the prefix: %s", prefix, linkTarget)
		}
	}
}
