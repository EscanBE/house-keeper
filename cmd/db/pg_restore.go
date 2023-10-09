package db

import (
	"fmt"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

//goland:noinspection SpellCheckingInspection
const (
	flagNoPubSub  = "no-pubsub"
	flagDataOnly  = "data-only"
	flagSuperUser = "superuser"
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
		flagPasswordFile,
		"",
		"file path which store password of the user which will be used to backup the database",
	)

	cmd.PersistentFlags().String(
		flagToolFile,
		"",
		"custom file path for the pg_restore utility",
	)

	cmd.PersistentFlags().String(
		flagSuperUser,
		"",
		fmt.Sprintf("specify the superuser user name to disable triggers when restoring in data-only mode. This is relevant only if --%s (enabled by default) is used", flagDataOnly),
	)

	cmd.PersistentFlags().Bool(
		flagDataOnly,
		true,
		fmt.Sprintf("restore data only, triggers are disabled, requires --%s flag", flagSuperUser),
	)

	cmd.PersistentFlags().Bool(
		flagNoPubSub,
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

	host, _ := cmd.Flags().GetString(flagHost)
	host = strings.TrimSpace(host)

	port, _ := cmd.Flags().GetUint16(flagPort)
	if port == 0 {
		panic(fmt.Errorf("missing value for mandatory flag --%s", flagPort))
	}

	dbName, _ := cmd.Flags().GetString(flagDbName)
	dbName = strings.TrimSpace(dbName)
	if len(dbName) < 1 {
		panic(fmt.Errorf("missing value for mandatory flag --%s", flagDbName))
	}

	userName, _ := cmd.Flags().GetString(flagUsername)
	userName = strings.TrimSpace(userName)

	dataOnly, _ := cmd.Flags().GetBool(flagDataOnly)

	superUser, _ := cmd.Flags().GetString(flagSuperUser)
	superUser = strings.TrimSpace(superUser)
	if dataOnly {
		if len(superUser) == 0 {
			panic(fmt.Sprintf("flag --%s is mandatory when --%s=true (default)", flagSuperUser, flagDataOnly))
		}
	} else {
		if len(superUser) > 0 {
			panic(fmt.Sprintf("flag --%s is not allowed when --%s=false", flagSuperUser, flagDataOnly))
		}
	}

	schema, _ := cmd.Flags().GetString(flagSchema)
	schema = strings.TrimSpace(schema)
	if len(schema) > 0 {
		panic(fmt.Sprintf("not yet supported flag --%s", flagSchema))
	}

	toolName := "pg_restore"
	customToolName, _ := cmd.Flags().GetString(flagToolFile)
	customToolName = strings.TrimSpace(customToolName)
	if len(customToolName) > 0 {
		_, err = os.Stat(customToolName)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("custom pg_restore file path does not exists: %s", customToolName))
		}

		toolName = customToolName
	}

	var envVars []string

	passwordFile, _ := cmd.Flags().GetString(flagPasswordFile)
	if len(passwordFile) < 1 {
		if len(strings.TrimSpace(os.Getenv(constants.ENV_PG_PASSWORD))) < 1 {
			panic(fmt.Errorf("missing password for user %s, either environment variable %s or flag --%s is required", userName, constants.ENV_PG_PASSWORD, flagPasswordFile))
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

		fipPerm := fip.Mode().Perm()
		errPerm := utils.ValidatePasswordFileMode(fipPerm)
		if errPerm != nil {
			fmt.Printf("Incorrect permission '%o' of password file: %s\n", fipPerm, errPerm)
			fmt.Printf("Suggest setting permission to '%o'\n", constants.RECOMMENDED_FILE_PERMISSION)
			os.Exit(1)
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
	noPubSub, _ := cmd.Flags().GetBool(flagNoPubSub)
	if noPubSub {
		restoreArgs = append(restoreArgs, "--no-publications")
		restoreArgs = append(restoreArgs, "--no-subscriptions")
	}
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
	fmt.Println("Begin restore", inputFilePath, "at", utils.NowStr())

	ec := utils.LaunchApp(toolName, restoreArgs, envVars, false)
	if ec != 0 {
		fmt.Println("Failed to restore", inputFilePath, "at", utils.NowStr())
		os.Exit(ec)
	}

	fmt.Println("Finished restore", inputFilePath, "at", utils.NowStr())
}
