package cmd

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
)

// budCmd represents the Butler install command
var budCmd = &cobra.Command{
	Use:   "bud",
	Args:  cobra.NoArgs,
	Short: fmt.Sprintf("Install Butler binary \"%s\"", constants.BUTLER_BINARY_NAME),
	Run: func(cmd *cobra.Command, args []string) {
		home, errGetUserHomeDir := os.UserHomeDir()
		if errGetUserHomeDir != nil {
			libutils.PrintlnStdErr("ERR: failed to get home directory:", errGetUserHomeDir.Error())
			os.Exit(1)
		}

		netrcPath := path.Join(home, ".netrc")
		isNetrcFileExists, err := utils.IsFileAndExists(netrcPath)
		if err != nil {
			libutils.PrintlnStdErr("ERR: failed to check if .netrc file exists:", err.Error())
			os.Exit(1)
		}

		var createNetrcFile bool
		if !isNetrcFileExists {
			createNetrcFile = true
		} else {
			bz, _ := os.ReadFile(netrcPath)
			createNetrcFile = len(bz) < 23
		}

		if createNetrcFile {
			gitUser := readStdIn("Git user:")
			gitToken := readStdIn("Git token:")
			netrcContent := fmt.Sprintf("\nmachine github.com\nlogin %s\npassword %s\n", gitUser, gitToken)

			outputFile, err := os.OpenFile(netrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				libutils.PrintlnStdErr("ERR: failed to open .netrc file:", err)
				os.Exit(1)
			}

			if _, err := outputFile.WriteString(netrcContent); err != nil {
				libutils.PrintlnStdErr("ERR: failed to write to .netrc file:", err)
				_ = outputFile.Close()
				os.Exit(1)
			}

			_ = outputFile.Close()
		}

		butlerRepoPath := path.Join(home, constants.BUTLER_REPO_DIR_NAME)

		runGit := func(workingDir string, args ...string) {
			ec := utils.LaunchAppWithDirectStd("git", append([]string{"-C", workingDir}, args...), nil)
			if ec != 0 {
				libutils.PrintlnStdErr("ERR: Exited with code", ec)
				os.Exit(ec)
			}
		}

		dirInfo, err := os.Stat(butlerRepoPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Clone Butler repo to", butlerRepoPath)

				runGit(home, "clone", constants.BUTLER_REPO)
			} else {
				libutils.PrintlnStdErr("Failed to check if Butler project directory exists:", err.Error())
				os.Exit(1)
			}
		} else {
			if !dirInfo.IsDir() {
				libutils.PrintlnStdErr("Butler project path is not a directory:", butlerRepoPath)
				os.Exit(1)
			}

			fmt.Println("Using existing Butler repo", butlerRepoPath)
		}

		fmt.Println("Updating Butler")
		runGit(butlerRepoPath, "reset", "--hard")
		runGit(butlerRepoPath, "clean", "-fd")
		runGit(butlerRepoPath, "checkout", "main")
		runGit(butlerRepoPath, "pull")

		fmt.Println("Installing Butler")
		ec := utils.LaunchAppWithSetup("make", []string{"install"}, func(launchCmd *exec.Cmd) {
			launchCmd.Dir = butlerRepoPath
			launchCmd.Stdin = os.Stdin
			launchCmd.Stdout = os.Stdout
			launchCmd.Stderr = os.Stderr
		})
		if ec != 0 {
			libutils.PrintlnStdErr("ERR: Exited with code", ec)
			os.Exit(ec)
		}

		fmt.Println("Butler version:")
		_ = utils.LaunchAppWithDirectStd(constants.BUTLER_BINARY_NAME, []string{"version"}, nil)
	},
}

func init() {
	rootCmd.AddCommand(budCmd)
}

func readStdIn(question string) string {
	fmt.Print(question)
	var input string
	_, _ = fmt.Scanln(&input)
	if len(input) < 1 {
		libutils.PrintlnStdErr("ERR: input is empty, please input again")
		return readStdIn(question)
	}
	return input
}
