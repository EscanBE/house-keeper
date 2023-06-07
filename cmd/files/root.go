package files

import (
	"github.com/spf13/cobra"
)

const (
	flagToolFile = "tool-file"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "files",
		Aliases: []string{"f"},
		Short:   "Interacting with files",
	}

	cmd.AddCommand(
		ListingCommands(),
		RsyncCommands(),
		ChecksumCommands(),
	)

	return cmd
}
