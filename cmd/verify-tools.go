package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

// verifyToolsCmd represents the verify-tools command, it checks required tools exists
var verifyToolsCmd = &cobra.Command{
	Use:     "verify-tools",
	Short:   "Checks required tools exists",
	Aliases: []string{"verify"},
	Run: func(cmd *cobra.Command, args []string) {
		var anyError bool

		cmdAppPgDump := exec.Command("pg_dump", "--help")
		if err := cmdAppPgDump.Run(); err != nil {
			fmt.Println("pg_dump might not exists", err)
			anyError = true
		}

		cmdAppRsync := exec.Command("rsync", "--help")
		if err := cmdAppRsync.Run(); err != nil {
			fmt.Println("rsync might not exists", err)
			anyError = true
		}

		cmdAppSshPass := exec.Command("sshpass", "-V")
		if err := cmdAppSshPass.Run(); err != nil {
			fmt.Println("sshpass might not exists", err)
			anyError = true
		}

		cmdAppSha1Sum := exec.Command("sha1sum", "--version")
		if err := cmdAppSha1Sum.Run(); err != nil {
			cmdAppShaSum := exec.Command("shasum", "--version")
			if err := cmdAppShaSum.Run(); err != nil {
				fmt.Println("both application shasum and sha1sum might not exists", err)
				anyError = true
			}
		}

		if !anyError {
			fmt.Println("Successfully checking, all tools were installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyToolsCmd)
}
