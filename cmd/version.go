package cmd

import (
	"fmt"
	"github.com/EscanBE/go-app-name/constants"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command, it prints the current version of the binary
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show binary version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constants.APP_NAME)
		fmt.Printf("Commit: %s\n", constants.COMMIT_HASH)
		fmt.Printf("Version: %s\n", constants.VERSION)
		fmt.Printf("Build Date: %s\n", constants.BUILD_DATE)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
