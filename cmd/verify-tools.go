package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

// verifyToolsCmd represents the verify-tools command, it checks required tools exists
var verifyToolsCmd = &cobra.Command{
	Use:   "verify-tools",
	Short: "Checks required tools exists",
	Run: func(cmd *cobra.Command, args []string) {
		var anyError bool

		cmdApp := exec.Command("pg_dump", "--help")
		if err := cmdApp.Run(); err != nil {
			fmt.Println("pg_dump might not exists", err)
			anyError = true
		}

		cmdApp = exec.Command("rsync", "--help")
		if err := cmdApp.Run(); err != nil {
			fmt.Println("rsync might not exists", err)
			anyError = true
		}

		cmdApp = exec.Command("sshpass", "-V")
		if err := cmdApp.Run(); err != nil {
			fmt.Println("sshpass might not exists", err)
			anyError = true
		}

		if !anyError {
			fmt.Println("Successfully checking, all tools were installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyToolsCmd)
}
