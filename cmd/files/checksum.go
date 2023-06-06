package files

import (
	"fmt"
	"github.com/EscanBE/go-ienumerable/goe"
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
		Use:     "checksum [file1] [file2...]",
		Short:   "Checksum file using shasum/sha1sum",
		Aliases: []string{"shasum", "sha1sum"},
		Args:    cobra.MinimumNArgs(1),
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
	var toolName string

	customToolName, _ := cmd.Flags().GetString(constants.FLAG_TOOL_FILE)
	customToolName = strings.TrimSpace(customToolName)
	if len(customToolName) > 0 {
		_, err := os.Stat(customToolName)
		if err != nil {
			if os.IsNotExist(err) {
				panic(fmt.Errorf("custom tool file does not exists: %s", customToolName))
			}

			panic(errors.Wrap(err, "problem while checking custom tool file path"))
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

	files := goe.NewIEnumerable(args...).SelectNewValue(func(file string) string {
		return strings.TrimSpace(file)
	}).Where(func(file string) bool {
		return len(file) > 0
	}).Distinct(nil).ToArray()

	if len(files) < 1 {
		panic("no file was provided")
	}

	checkInputFile := func(file string) {
		_, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				panic(fmt.Errorf("file does not exists: %s", file))
			}

			panic(errors.Wrap(err, fmt.Sprintf("problem while checking target file %s", file)))
		}
	}

	// check files before start checksum files
	for _, file := range files {
		checkInputFile(file)
	}

	// start checksum files one by one
	for _, file := range files {
		checkInputFile(file)

		fmt.Println("start checksum file", file)
		exitCode := utils.LaunchAppWithOutputCallback(toolName, []string{file}, os.Environ(), outputCb, outputCb)
		if exitCode != 0 {
			fmt.Println("failed to checksum file", file)
			os.Exit(exitCode)
		}
	}
}

func writeToOutputFile(outputFilePath string, content string) {
	if len(outputFilePath) < 1 {
		return
	}

	outputFile, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to open file [%s] to write content [%s]", outputFilePath, content)
	}

	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

	if len(content) < 1 {
		return
	}

	if _, err := outputFile.WriteString(content); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to append content [%s] to output file [%s]", content, outputFilePath)
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
