package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/glass-cms/glasscms/cmd"
	"github.com/stretchr/testify/require"
)

func Test_ParseCommand(t *testing.T) {
	t.Parallel()

	tempDir, err := os.MkdirTemp("", "glasscms")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("json files", func(t *testing.T) {
		command := cmd.NewParseCommand()
		command.SetArgs([]string{"../docs/commands", fmt.Sprintf("--%s", cmd.ArgOutput), tempDir})

		err = command.Command.Execute()
		require.NoError(t, err)

		// Check that there are multiple json files in the output directory
		files, err := os.ReadDir(tempDir)
		require.NoError(t, err)
		require.Greater(t, len(files), 0)

		for _, file := range files {
			require.Contains(t, file.Name(), ".json")
		}

		os.RemoveAll(tempDir)
	})

	t.Run("yaml files", func(t *testing.T) {
		command := cmd.NewParseCommand()
		command.SetArgs([]string{"../docs/commands", fmt.Sprintf("--%s", cmd.ArgOutput), tempDir, fmt.Sprintf("--%s", cmd.ArgFormat), cmd.FormatYAML})

		err = command.Command.Execute()
		require.NoError(t, err)

		// Check that there are multiple yaml files in the output directory
		files, err := os.ReadDir(tempDir)
		require.NoError(t, err)
		require.Greater(t, len(files), 0)

		for _, file := range files {
			require.Contains(t, file.Name(), ".yaml")
		}

		os.RemoveAll(tempDir)
	})
}
