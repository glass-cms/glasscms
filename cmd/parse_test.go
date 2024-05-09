package cmd_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/glass-cms/glasscms/cmd"
	"github.com/glass-cms/glasscms/item"
	"github.com/stretchr/testify/require"
)

func Test_ParseCommand(t *testing.T) {
	t.Parallel()

	tempDir, err := os.MkdirTemp("", "glasscms")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	command := cmd.NewParseCommand()
	command.SetArgs([]string{"../docs/commands", fmt.Sprintf("--%s", cmd.ArgOutput), tempDir})

	err = command.Command.Execute()
	require.NoError(t, err)

	// Assert that the output directory contains the expected JSON file.
	// And that the file is not empty.

	// Read the JSON file.
	var items []*item.Item
	file, err := os.Open(filepath.Join(tempDir, "output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(&items)
	if err != nil {
		t.Fatal(err)
	}

	require.NotEmpty(t, items)
}
