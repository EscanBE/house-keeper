package cmd

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/EscanBE/house-keeper/cmd/utils"
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

		fmt.Println("Mandatory tools checking...")

		if !utils.HasBinaryName("pg_dump") {
			libutils.PrintlnStdErr("- \"pg_dump\" might not exists")
			anyMandatoryToolsError = true
		}

		if !utils.HasBinaryName("rsync") {
			libutils.PrintlnStdErr("- \"rsync\" might not exists")
			anyMandatoryToolsError = true
		}

		if !utils.HasBinaryName("sshpass") {
			libutils.PrintlnStdErr("- \"sshpass\" might not exists")
			anyMandatoryToolsError = true
		}

		if !utils.HasBinaryName("sha1sum") && !utils.HasBinaryName("shasum") {
			libutils.PrintlnStdErr("- Both applications \"shasum\" and \"sha1sum\" might not exists")
			anyMandatoryToolsError = true
		}

		if !utils.HasBinaryName("aria2c") {
			libutils.PrintlnStdErr("- \"aria2c\" might not exists")
			anyMandatoryToolsError = true
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
		//goland:noinspection SpellCheckingInspection
		checkToolPossiblyExists("lz4", "sudo apt install snapd -y && sudo snap install lz4", possiblyMissingOptionalTools)
		checkToolPossiblyExists("psql", "sudo apt install postgresql-client", possiblyMissingOptionalTools)
		//goland:noinspection SpellCheckingInspection
		checkToolPossiblyExists("rclone", "sudo -v ; curl https://rclone.org/install.sh | sudo bash", possiblyMissingOptionalTools)
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
