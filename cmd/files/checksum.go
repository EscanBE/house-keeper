package files

import (
	"bufio"
	"fmt"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
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

	cmd.PersistentFlags().String(
		constants.FLAG_OUTPUT_FILE,
		"",
		"append output to file",
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

	var outputFile *os.File
	outputFilePath, _ := cmd.Flags().GetString(constants.FLAG_OUTPUT_FILE)
	outputFilePath = strings.TrimSpace(outputFilePath)
	if len(outputFilePath) > 0 {
		outputFile, err = os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(errors.Wrap(err, "failed to open file for append: "+outputFilePath))
		}
	} else {
		outputFile = nil
	}

	defer func() {
		if outputFile != nil {
			_ = outputFile.Close()
		}
	}()

	appendOutput := func(text string) {
		if outputFile == nil {
			return
		}

		if _, err := outputFile.WriteString(text); err != nil {
			fmt.Println("failed to append to output file", err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			oScan := rsyncStdOutScanner.Scan()
			eScan := rsyncStdErrScanner.Scan()
			if oScan {
				msg := rsyncStdOutScanner.Text()
				fmt.Println(msg)
				appendOutput(msg + "\n")
			}
			if eScan {
				msg := rsyncStdErrScanner.Text()
				fmt.Println(msg)
				appendOutput(msg + "\n")
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
