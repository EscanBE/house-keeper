package files

import (
	"fmt"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
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
		"custom checksum tool's file path",
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

	outputFilePath, _ := cmd.Flags().GetString(constants.FLAG_OUTPUT_FILE)
	outputFilePath = strings.TrimSpace(outputFilePath)

	writeToOutputFile(outputFilePath, "") // test write

	outputCb := func(msg string) {
		writeToOutputFile(outputFilePath, msg+"\n")
	}
	utils.LaunchAppWithOutputCallback(toolName, []string{file}, os.Environ(), outputCb, outputCb)
}

func writeToOutputFile(outputFilePath string, content string) {
	if len(outputFilePath) < 1 {
		return
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(errors.Wrap(err, "failed to open output file for append: "+outputFilePath))
	}

	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

	if len(content) < 1 {
		return
	}

	if _, err := outputFile.WriteString(content); err != nil {
		fmt.Printf("failed to append content [%s] to output file [%s]", content, outputFilePath)
		fmt.Println(err)
	}
}
