package db

import (
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	flagOutputFile = "output-file"
	flagToolFile   = "tool-file"

	flagHost         = "host"
	flagPort         = "port"
	flagDbName       = "dbname"
	flagUsername     = "username"
	flagSchema       = "schema"
	flagPasswordFile = "password-file"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database tools",
	}

	cmd.AddCommand(
		PgDumpCommands(),
		PgRestoreCommands(),
	)

	utils.AddFlagWorkingDir(cmd)

	cmd.PersistentFlags().String(
		flagHost,
		"localhost",
		"database host",
	)

	cmd.PersistentFlags().Uint16(
		flagPort,
		5432,
		"database port",
	)

	cmd.PersistentFlags().String(
		flagDbName,
		"postgres",
		"database name",
	)

	cmd.PersistentFlags().String(
		flagUsername,
		"postgres",
		"username to be used to connect to database",
	)

	return cmd
}
