package fs_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
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
		{content: "# Daily Note\nToday's thoughts...", depth: 0, pattern: "*.md"},
		{content: "# Meeting Notes\nProject discussion...", depth: 0, pattern: "*.md"},
		{content: "# Article Draft\nContent for blog...", depth: 0, pattern: "*.md"},
	}

	tempDir, err := createTempDirWithFiles(fileData)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	src, err := fs.NewSourcer(tempDir, nil)
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
	assert.Equal(t, fs.ErrDone, err)
}

func TestFileSystemSourcer_Size(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input []testFile
		want  int
	}{
		"empty vault": {
			input: []testFile{},
			want:  0,
		},
		"single note": {
			input: []testFile{{content: "# Daily Note\nMy thoughts today...", depth: 0, pattern: "*.md"}},
			want:  1,
		},
		"multiple notes in vault": {
			input: []testFile{
				{content: "# Project Planning\nIdeas for the project...", depth: 0, pattern: "*.md"},
				{content: "# Meeting Notes\nDiscussion points...", depth: 0, pattern: "*.md"},
				{content: "# Research Article\nFindings and analysis...", depth: 0, pattern: "*.md"},
			},
			want: 3,
		},
		"nested vault structure": {
			input: []testFile{
				{content: "# Home\nWelcome to my vault", depth: 1, pattern: "*.md"},
				{content: "# Projects Overview\nCurrent projects...", depth: 2, pattern: "*.md"},
				{content: "# Archive Note\nOld content...", depth: 3, pattern: "*.md"},
			},
			want: 3,
		},
		"non-content files ignored": {
			input: []testFile{
				{content: "config data", depth: 0, pattern: "*.json"},
				{content: "image data", depth: 0, pattern: "*.png"},
				{content: "text file", depth: 0, pattern: "*.txt"},
			},
			want: 0,
		},
		"mixed content types": {
			input: []testFile{
				{content: "# Note\nContent here...", depth: 0, pattern: "*.md"},
				{content: "config", depth: 0, pattern: "*.json"},
				{content: "# Another Note\nMore content...", depth: 0, pattern: "*.md"},
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

			sourcer, err := fs.NewSourcer(tempDir, nil)
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

	fp, err := createTempDirWithFiles([]testFile{{content: "# Welcome\nMy knowledge vault", depth: 0, pattern: "*.md"}})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fp)

	require.NoError(t, fs.IsValidFileSystemSource(fp))
	require.Error(t, fs.IsValidFileSystemSource("non-existent"))
}

func TestFileSystemSourcer_IgnorePatterns(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		files          []testFile
		ignorePatterns []string
		expectedCount  int
	}{
		"no ignore patterns - sync all notes": {
			files: []testFile{
				{content: "# Home\nWelcome to my vault", depth: 0, pattern: "*.md"},
				{content: "# Project Note\nDetails about project", depth: 1, pattern: "*.md"},
				{content: "# Meeting Summary\nMeeting outcomes", depth: 1, pattern: "*.md"},
			},
			ignorePatterns: nil,
			expectedCount:  3,
		},
		"ignore templates folder": {
			files: []testFile{
				{content: "# Home\nMain vault note", depth: 0, pattern: "*.md"},
				{content: "# Project Note\nProject details", depth: 1, pattern: "*.md"}, // in dir0/
			},
			ignorePatterns: []string{"Templates"},
			expectedCount:  2, // should not affect existing dirs
		},
		"ignore obsidian system folders": {
			files: []testFile{
				{content: "# Daily Note\nToday's notes", depth: 0, pattern: "*.md"},
				{content: "# Archive Note\nOld content", depth: 1, pattern: "*.md"}, // in dir0/
			},
			ignorePatterns: []string{".*"},
			expectedCount:  2, // should not affect non-hidden dirs
		},
		"ignore draft folders": {
			files: []testFile{
				{content: "# Published Article\nFinished content", depth: 0, pattern: "*.md"},
				{content: "# Archive Note\nOld notes", depth: 1, pattern: "*.md"}, // in dir0/
			},
			ignorePatterns: []string{"Draft*"},
			expectedCount:  2, // should ignore any Draft folders
		},
		"multiple cms ignore patterns": {
			files: []testFile{
				{content: "# Home\nMain page", depth: 0, pattern: "*.md"},
				{content: "# Project Note\nProject info", depth: 1, pattern: "*.md"},
				{content: "# Deep Archive\nOld content", depth: 2, pattern: "*.md"},
			},
			ignorePatterns: []string{"dir1", ".*"},
			expectedCount:  2, // should ignore dir1 folders and hidden dirs
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create temp directory structure
			tempDir, err := createTempDirWithFiles(tt.files)
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			// Create sourcer with ignore patterns
			sourcer, err := fs.NewSourcer(tempDir, tt.ignorePatterns)
			if err != nil {
				t.Fatal(err)
			}

			// Verify expected file count
			assert.Equal(t, tt.expectedCount, sourcer.Size())
		})
	}
}

func TestFileSystemSourcer_IgnorePatternsWithObsidianVault(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		ignorePatterns []string
		expectedCount  int
	}{
		"ignore templates only": {
			ignorePatterns: []string{"Templates"},
			expectedCount:  4, // all except Templates/note-template.md
		},
		"ignore obsidian system folders": {
			ignorePatterns: []string{".*"},
			expectedCount:  4, // all except .obsidian/config.md
		},
		"ignore drafts and templates": {
			ignorePatterns: []string{"Templates", "Drafts"},
			expectedCount:  3, // home.md, .obsidian/config.md, Projects/project-a.md
		},
		"typical obsidian vault sync": {
			ignorePatterns: []string{"Templates", ".*", "Drafts"},
			expectedCount:  2, // home.md, Projects/project-a.md
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempDir := createObsidianVaultStructure(t)
			defer os.RemoveAll(tempDir)

			sourcer, err := fs.NewSourcer(tempDir, tt.ignorePatterns)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCount, sourcer.Size(), "Expected %d files but got %d", tt.expectedCount, sourcer.Size())
		})
	}
}

// createObsidianVaultStructure creates a realistic Obsidian vault structure for testing.
func createObsidianVaultStructure(t *testing.T) string {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "test-ignore")
	if err != nil {
		t.Fatal(err)
	}

	// Create Obsidian vault structure
	dirs := []string{"Templates", ".obsidian", "Drafts", "Projects"}
	files := map[string]string{
		"home.md":                    "# Welcome to My Vault\n\nThis is my knowledge base.",
		"Templates/note-template.md": "# Template\n\n## Date: {{date}}\n\n## Notes:\n\n",
		".obsidian/config.md":        "Obsidian configuration file",
		"Drafts/draft-article.md":    "# Draft Article\n\nWork in progress content...",
		"Projects/project-a.md":      "# Project A\n\n## Status: Active\n\n## Goals:\n- Complete phase 1",
	}

	// Create directories
	for _, dir := range dirs {
		if createErr := os.MkdirAll(filepath.Join(tempDir, dir), 0755); createErr != nil {
			t.Fatal(createErr)
		}
	}

	// Create files
	for filePath, content := range files {
		fullPath := filepath.Join(tempDir, filePath)
		if writeErr := os.WriteFile(fullPath, []byte(content), 0644); writeErr != nil {
			t.Fatal(writeErr)
		}
	}

	return tempDir
}
