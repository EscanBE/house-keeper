package files

import (
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/spf13/cobra"
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

	utils.AddFlagWorkingDir(cmd)

	return cmd
}
