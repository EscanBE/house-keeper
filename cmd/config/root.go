package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure something",
	}

	cmd.AddCommand(
		ConfigureSshCommands(),
	)

	return cmd
}

func isFileExists(file string) bool {
	fi, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		panic(errors.Wrap(err, fmt.Sprintf("problem while checking target file %s", file)))
	}

	if fi.IsDir() {
		panic(fmt.Sprintf("require file but found directory: %s", file))
	}

	return true
}
