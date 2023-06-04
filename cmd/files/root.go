package files

import (
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
		RsyncCommands(),
	)

	return cmd
}
