package cmd_test

import (
	"testing"

	"github.com/glass-cms/glasscms/cmd"
	"github.com/stretchr/testify/require"
)

func Test_SyncCommandURLValidation(t *testing.T) {
	command := cmd.NewSyncCommand()
	command.SetArgs([]string{
		"filesystem",
		"../docs/commands",
		"--server", "localhost:8080",
	})

	err := command.Command.Execute()
	require.Error(t, err)
}

func Test_parseIgnorePatterns(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []string
	}{
		"empty patterns - sync all content": {
			input:    "",
			expected: nil,
		},
		"ignore template folder": {
			input:    "Templates",
			expected: []string{"Templates"},
		},
		"ignore common obsidian folders": {
			input:    ".obsidian,.trash,Templates,Daily Notes",
			expected: []string{".obsidian", ".trash", "Templates", "Daily Notes"},
		},
		"ignore folders with spaces in names": {
			input:    "Daily Notes, Meeting Notes , Project Archive",
			expected: []string{"Daily Notes", "Meeting Notes", "Project Archive"},
		},
		"ignore drafts and private folders": {
			input:    "Drafts,,Private,  ,Archive",
			expected: []string{"Drafts", "Private", "Archive"},
		},
		"whitespace only pattern": {
			input:    "   ",
			expected: nil,
		},
		"malformed comma pattern": {
			input:    ",,,",
			expected: nil,
		},
		"ignore hidden and temporary folders": {
			input:    ".*,*temp*,*draft*,*backup*",
			expected: []string{".*", "*temp*", "*draft*", "*backup*"},
		},
		"real obsidian vault structure": {
			input:    ".obsidian, .trash, Templates, Daily Notes, *.tmp, Archive",
			expected: []string{".obsidian", ".trash", "Templates", "Daily Notes", "*.tmp", "Archive"},
		},
		"content management workflow": {
			input:    "Inbox,Drafts,Private,Admin,*.conflict*",
			expected: []string{"Inbox", "Drafts", "Private", "Admin", "*.conflict*"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := cmd.ParseIgnorePatterns(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
