package db

import (
	"fmt"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

// PgRestoreCommands registers a sub-tree of commands
func PgRestoreCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pg_restore [file_name]",
		Short: "Restore DB using backup file (PostgreSQL)",
		Args:  cobra.ExactArgs(1),
		Run:   restorePgDatabase,
	}

	cmd.PersistentFlags().String(
		constants.FLAG_PASSWORD_FILE,
		"",
		"file path which store password of the user which will be used to backup the database",
	)

	cmd.PersistentFlags().String(
		constants.FLAG_TOOL_FILE,
		"",
		"custom file path for the pg_restore utility",
	)

	cmd.PersistentFlags().String(
		constants.FLAG_SUPER_USER,
		"",
		fmt.Sprintf("specify the superuser user name to disable triggers when restoring in data-only mode. This is relevant only if --%s (enabled by default) is used", constants.FLAG_DATA_ONLY),
	)

	cmd.PersistentFlags().Bool(
		constants.FLAG_DATA_ONLY,
		true,
		fmt.Sprintf("restore data only, triggers are disabled, requires --%s flag", constants.FLAG_SUPER_USER),
	)

	cmd.PersistentFlags().Bool(
		constants.FLAG_NO_PUB_SUB,
		true,
		"do not output commands to restore publications/subscriptions, even if the archive contains them.",
	)

	return cmd
}

func restorePgDatabase(cmd *cobra.Command, args []string) {
	inputFilePath := args[0]
	inputFilePath = strings.TrimSpace(inputFilePath)
	if len(inputFilePath) < 1 {
		panic("bad input backup file")
	}

	// workingDir := utils.ReadFlagWorkingDir(cmd)

	_, err := os.Stat(inputFilePath)
	if err == nil {
		// ok
	} else {
		if os.IsNotExist(err) {
			panic("bad input backup file, file does not exists")
		} else {
			panic(errors.Wrap(err, fmt.Sprintf("problem when checking input file %s", inputFilePath)))
		}
	}

	host, _ := cmd.Flags().GetString(constants.FLAG_HOST)
	host = strings.TrimSpace(host)

	port, _ := cmd.Flags().GetUint16(constants.FLAG_PORT)
	if port == 0 {
		panic(fmt.Errorf("missing value for mandatory flag --%s", constants.FLAG_PORT))
	}

	dbName, _ := cmd.Flags().GetString(constants.FLAG_DB_NAME)
	dbName = strings.TrimSpace(dbName)
	if len(dbName) < 1 {
		panic(fmt.Errorf("missing value for mandatory flag --%s", constants.FLAG_DB_NAME))
	}

	userName, _ := cmd.Flags().GetString(constants.FLAG_USER_NAME)
	userName = strings.TrimSpace(userName)

	dataOnly, _ := cmd.Flags().GetBool(constants.FLAG_DATA_ONLY)

	superUser, _ := cmd.Flags().GetString(constants.FLAG_SUPER_USER)
	superUser = strings.TrimSpace(superUser)
	if dataOnly {
		if len(superUser) == 0 {
			panic(fmt.Sprintf("flag --%s is mandatory when --%s=true (default)", constants.FLAG_SUPER_USER, constants.FLAG_DATA_ONLY))
		}
	} else {
		if len(superUser) > 0 {
			panic(fmt.Sprintf("flag --%s is not allowed when --%s=false", constants.FLAG_SUPER_USER, constants.FLAG_DATA_ONLY))
		}
	}

	schema, _ := cmd.Flags().GetString(constants.FLAG_SCHEMA)
	schema = strings.TrimSpace(schema)
	if len(schema) > 0 {
		panic(fmt.Sprintf("not yet supported flag --%s", constants.FLAG_SCHEMA))
	}

	toolName := "pg_restore"
	customToolName, _ := cmd.Flags().GetString(constants.FLAG_TOOL_FILE)
	customToolName = strings.TrimSpace(customToolName)
	if len(customToolName) > 0 {
		_, err = os.Stat(customToolName)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("custom pg_restore file path does not exists: %s", customToolName))
		}

		toolName = customToolName
	}

	var envVars []string

	passwordFile, _ := cmd.Flags().GetString(constants.FLAG_PASSWORD_FILE)
	if len(passwordFile) < 1 {
		if len(strings.TrimSpace(os.Getenv(constants.ENV_PG_PASSWORD))) < 1 {
			panic(fmt.Errorf("missing password for user %s, either environment variable %s or flag --%s is required", userName, constants.ENV_PG_PASSWORD, constants.FLAG_PASSWORD_FILE))
		}

		envVars = os.Environ()
	} else {
		fip, err := os.Stat(passwordFile)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("supplied password file does not exists %s", passwordFile))
		}

		bz, err := os.ReadFile(passwordFile)
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("failed to read password file %s", passwordFile)))
		}

		if fip.Mode().Perm() != constants.FILE_PERMISSION_400 &&
			fip.Mode().Perm() != constants.FILE_PERMISSION_600 &&
			fip.Mode().Perm() != constants.FILE_PERMISSION_700 {
			//goland:noinspection GoBoolExpressions
			panic(fmt.Errorf("incorrect permission of password file, must be %s or %s or %s", constants.FILE_PERMISSION_400_STR, constants.FILE_PERMISSION_600_STR, constants.FILE_PERMISSION_700_STR))
		}

		pgPassword := strings.TrimSpace(string(bz))
		if len(pgPassword) < 1 {
			panic(fmt.Errorf("password file is empty: %s", passwordFile))
		}

		envVars = append(envVars, fmt.Sprintf("%s=%s", constants.ENV_PG_PASSWORD, pgPassword))
	}

	restoreArgs := make([]string, 0)
	if len(host) > 0 {
		restoreArgs = append(restoreArgs, fmt.Sprintf("--host=%s", host))
	}
	if port > 0 {
		restoreArgs = append(restoreArgs, fmt.Sprintf("--port=%d", port))
	}
	if len(userName) > 0 {
		restoreArgs = append(restoreArgs, fmt.Sprintf("--username=%s", userName))
	}
	restoreArgs = append(restoreArgs, fmt.Sprintf("--dbname=%s", dbName))
	restoreArgs = append(restoreArgs, "--single-transaction")
	restoreArgs = append(restoreArgs, "--no-publications")
	restoreArgs = append(restoreArgs, "--no-subscriptions")
	restoreArgs = append(restoreArgs, "--no-owner")
	if dataOnly {
		restoreArgs = append(restoreArgs, "--data-only")
		restoreArgs = append(restoreArgs, "--disable-triggers")
		if len(superUser) > 0 {
			restoreArgs = append(restoreArgs, fmt.Sprintf("--superuser=%s", superUser))
		} else {
			panic("require superuser")
		}
	}
	restoreArgs = append(restoreArgs, inputFilePath)

	fmt.Println("Input file:", inputFilePath)
	fmt.Println("Restore arguments:\n", toolName, strings.Join(restoreArgs, " "))
	fmt.Println("Begin restore", inputFilePath, "at", time.Now().Format("2006-Jan-02 15:04:05"))

	exitCode := utils.LaunchApp(toolName, restoreArgs, envVars)
	if exitCode == 0 {
		fmt.Println("Finished restore", inputFilePath, "at", time.Now().Format("2006-Jan-02 15:04:05"))
	} else {
		fmt.Println("Failed to restore", inputFilePath, "at", time.Now().Format("2006-Jan-02 15:04:05"))
	}
	os.Exit(exitCode)
}
