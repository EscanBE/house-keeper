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
		if len(constants.BUILD_FROM_SOURCE) > 0 {
			fmt.Println(constants.VERSION, "(build from source)")
		} else {
			fmt.Println(constants.VERSION)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
