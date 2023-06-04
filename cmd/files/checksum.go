package files

import (
	"bufio"
	"fmt"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// ChecksumCommands registers a sub-tree of commands
func ChecksumCommands() *cobra.Command {
	//goland:noinspection SpellCheckingInspection
	cmd := &cobra.Command{
		Use:     "checksum [file]",
		Short:   "Checksum file using shasum/sha1sum",
		Aliases: []string{"shasum", "sha1sum"},
		Args:    cobra.ExactArgs(1),
		Run:     checksumFile,
	}

	cmd.PersistentFlags().String(
		constants.FLAG_TOOL_FILE,
		"",
		"absolute file path of the checksum tool",
	)

	return cmd
}

func checksumFile(cmd *cobra.Command, args []string) {
	file := strings.TrimSpace(args[0])
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		panic(fmt.Errorf("file does not exists: %s", file))
	}

	var toolName string

	customToolName, _ := cmd.Flags().GetString(constants.FLAG_TOOL_FILE)
	customToolName = strings.TrimSpace(customToolName)
	if len(customToolName) > 0 {
		_, err := os.Stat(customToolName)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("custom tool file does not exists: %s", customToolName))
		}

		toolName = customToolName
	} else {
		cmdAppSha1Sum := exec.Command("sha1sum", "--version")
		if err := cmdAppSha1Sum.Run(); err == nil {
			toolName = "sha1sum"
		} else {
			cmdAppShaSum := exec.Command("shasum", "--version")
			if err := cmdAppShaSum.Run(); err == nil {
				toolName = "shasum"
			} else {
				panic("require at least either tools sha1sum or shasum")
			}
		}
	}

	rsyncCmd := exec.Command(toolName, file)

	rsyncCmd.Env = os.Environ()
	stdout, _ := rsyncCmd.StdoutPipe()
	stderr, _ := rsyncCmd.StderrPipe()
	rsyncStdOutScanner := bufio.NewScanner(stdout)
	rsyncStdErrScanner := bufio.NewScanner(stderr)
	err = rsyncCmd.Start()
	if err != nil {
		fmt.Println("problem when starting app", toolName, err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			oScan := rsyncStdOutScanner.Scan()
			eScan := rsyncStdErrScanner.Scan()
			if oScan {
				fmt.Println(rsyncStdOutScanner.Text())
			}
			if eScan {
				fmt.Println(rsyncStdErrScanner.Text())
			}
			if !oScan && !eScan {
				break
			}
		}
		err = rsyncCmd.Wait()
		if err != nil {
			fmt.Println("problem when waiting process", err)
		}
		defer wg.Done()
	}()

	wg.Wait()
}
