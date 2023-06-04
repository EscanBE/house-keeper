package cmd

import (
	"fmt"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command, it prints the current version of the binary
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show binary version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(constants.VERSION)
		if len(constants.BUILD_FROM_SOURCE) > 0 {
			fmt.Println(" (build from source)")
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
