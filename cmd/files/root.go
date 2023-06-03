package files

import (
	"github.com/EscanBE/house-keeper/constants"
	"github.com/spf13/cobra"
	"os"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Interacting with files",
	}

	cmd.AddCommand(
		ListingCommands(),
	)

	curDir, err := os.Getwd()
	if err != nil {
		panic("failed to get current directory")
	}

	cmd.PersistentFlags().String(
		constants.FLAG_WORKING_DIR,
		curDir,
		"the working directory",
	)

	return cmd
}
