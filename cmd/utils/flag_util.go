package utils

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	flagWorkingDir = "working-directory"
)

func AddFlagWorkingDir(cmd *cobra.Command) {
	curDir, err := os.Getwd()
	if err != nil {
		panic("failed to get current directory")
	}

	cmd.PersistentFlags().String(
		flagWorkingDir,
		curDir,
		"the working directory",
	)
}

func ReadFlagWorkingDir(cmd *cobra.Command) string {
	workingDir, _ := cmd.Flags().GetString(flagWorkingDir)
	workingDir = strings.TrimSpace(workingDir)
	if len(workingDir) < 1 {
		panic(fmt.Errorf("empty working directory"))
	}
	workingDirInfo, err := os.Stat(workingDir)
	if os.IsNotExist(err) {
		panic(fmt.Errorf("specified working directory does not exists: %s", workingDir))
	}
	if !workingDirInfo.IsDir() {
		panic(fmt.Errorf("specified working directory is not a directory"))
	}
	return workingDir
}
