package db

import (
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database tools",
	}

	cmd.AddCommand(
		BackupCommands(),
		PgRestoreCommands(),
	)

	utils.AddFlagWorkingDir(cmd)

	cmd.PersistentFlags().String(
		constants.FLAG_HOST,
		"localhost",
		"database host",
	)

	cmd.PersistentFlags().Uint16(
		constants.FLAG_PORT,
		5432,
		"database port",
	)

	cmd.PersistentFlags().String(
		constants.FLAG_DB_NAME,
		"postgres",
		"database name",
	)

	cmd.PersistentFlags().String(
		constants.FLAG_USER_NAME,
		"postgres",
		"database username",
	)

	return cmd
}
