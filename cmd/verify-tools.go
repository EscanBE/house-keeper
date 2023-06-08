package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// verifyToolsCmd represents the verify-tools command, it checks required tools exists
var verifyToolsCmd = &cobra.Command{
	Use:     "verify-tools",
	Short:   "Checks required tools exists",
	Aliases: []string{"verify"},
	Run: func(cmd *cobra.Command, args []string) {
		var anyMandatoryToolsError bool

		defer func() {
			if anyMandatoryToolsError {
				os.Exit(1)
			}
		}()

		cmdAppPgDump := exec.Command("pg_dump", "--help")
		if err := cmdAppPgDump.Run(); err != nil {
			fmt.Println("pg_dump might not exists", err)
			anyMandatoryToolsError = true
		}

		cmdAppRsync := exec.Command("rsync", "--help")
		if err := cmdAppRsync.Run(); err != nil {
			fmt.Println("rsync might not exists", err)
			anyMandatoryToolsError = true
		}

		cmdAppSshPass := exec.Command("sshpass", "-V")
		if err := cmdAppSshPass.Run(); err != nil {
			fmt.Println("sshpass might not exists", err)
			anyMandatoryToolsError = true
		}

		cmdAppSha1Sum := exec.Command("sha1sum", "--version")
		if err := cmdAppSha1Sum.Run(); err != nil {
			cmdAppShaSum := exec.Command("shasum", "--version")
			if err := cmdAppShaSum.Run(); err != nil {
				fmt.Println("both applications shasum and sha1sum might not exists", err)
				anyMandatoryToolsError = true
			}
		}

		if !anyMandatoryToolsError {
			fmt.Println("Successfully checking, all mandatory tools were installed")
		}

		possiblyMissingOptionalTools := make(map[string]string)
		defer func() {
			if len(possiblyMissingOptionalTools) > 0 {
				fmt.Println("(Additionally) The following optional tools might not be installed yet:")
				for toolName, installCommand := range possiblyMissingOptionalTools {
					if len(installCommand) > 0 {
						fmt.Printf(" - %s (%s)\n", toolName, installCommand)
					} else {
						fmt.Println(" -", toolName)
					}
				}
			}
		}()

		checkToolPossiblyExists("telnet", "sudo apt install telnet -y", possiblyMissingOptionalTools)
		checkToolPossiblyExists("htop", "sudo apt install htop -y", possiblyMissingOptionalTools)
		checkToolPossiblyExists("screen", "sudo apt install screen -y", possiblyMissingOptionalTools)
		checkToolPossiblyExists("wget", "sudo apt install wget -y", possiblyMissingOptionalTools)
		checkToolPossiblyExists("jq", "sudo apt install jq -y", possiblyMissingOptionalTools)
		checkToolPossiblyExists("lz4", "sudo apt install snapd -y && sudo snap install lz4", possiblyMissingOptionalTools)
		checkToolPossiblyExists("psql", "sudo apt install postgresql-client", possiblyMissingOptionalTools)
	},
}

func checkToolPossiblyExists(tool string, installCommand string, tracker map[string]string) {
	_, err := exec.LookPath(tool)
	if err != nil {
		tracker[tool] = installCommand
	}
}

func init() {
	rootCmd.AddCommand(verifyToolsCmd)
}
