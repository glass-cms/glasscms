package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/glass-cms/glasscms/cmd"
	"github.com/stretchr/testify/require"
)

func Test_ConvertCommandJSON(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "json*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	command := cmd.NewConvertCommand()
	command.SetArgs([]string{"../docs/commands", fmt.Sprintf("--%s", cmd.ArgOutput), tempDir})

	err = command.Command.Execute()
	require.NoError(t, err)

	// Check that there are multiple json files in the output directory
	files, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, file := range files {
		require.Contains(t, file.Name(), ".json")
	}
}

func Test_ConvertCommandYAML(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "yml*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	command := cmd.NewConvertCommand()
	command.SetArgs([]string{
		"../docs/commands",
		fmt.Sprintf("--%s", cmd.ArgOutput),
		tempDir,
		fmt.Sprintf("--%s", cmd.ArgFormat),
		cmd.FormatYAML})

	err = command.Command.Execute()
	require.NoError(t, err)

	// Check that there are multiple yaml files in the output directory
	files, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	require.NotEmpty(t, files)

	for _, file := range files {
		require.Contains(t, file.Name(), ".yaml")
	}
}
