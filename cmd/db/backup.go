package db

import (
	"fmt"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// BackupCommands registers a sub-tree of commands
func BackupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: fmt.Sprintf("** Deprecated command, use %s db pg_dump instead **", constants.BINARY_NAME),
		Args:  cobra.NoArgs,
		Run:   backupDatabase,
	}

	cmd.PersistentFlags().String(
		flagOutputFile,
		fmt.Sprintf("db-%s.dump", time.Now().Format("2006-01-02")),
		"specify name of the output backup file, file name only, default has layout: db-yyyy-MM-dd.dump",
	)

	cmd.PersistentFlags().String(
		flagPasswordFile,
		"",
		"file path which store password of the user which will be used to backup the database",
	)

	cmd.PersistentFlags().String(
		flagToolFile,
		"",
		"custom file path (absolute) for the backup utility (eg pg_dump of PostgreSQL)",
	)

	cmd.PersistentFlags().String(
		flagSchema,
		"public",
		"specify schema to backup",
	)

	return cmd
}

func backupDatabase(cmd *cobra.Command, _ []string) {
	fmt.Printf("** WARNING ** this is a deprecated function, use [%s db pg_dump ...] instead!\n", constants.BINARY_NAME)

	outputFileName, _ := cmd.Flags().GetString(flagOutputFile)
	outputFileName = strings.TrimSpace(outputFileName)
	if len(outputFileName) < 1 {
		panic(fmt.Errorf("missing value for mandatory flag --%s", flagOutputFile))
	}

	dir, outputFileName := path.Split(outputFileName)
	if len(dir) > 0 {
		panic("output file name must be file name alone, can not contains directory part")
	}

	workingDir := utils.ReadFlagWorkingDir(cmd)

	outputFilePath, err := filepath.Abs(path.Join(workingDir, outputFileName))
	if err != nil {
		panic(errors.Wrap(err, "failed to convert into absolute path"))
	}

	_, err = os.Stat(outputFilePath)
	if err == nil {
		panic(fmt.Errorf("output file already exists: %s", outputFilePath))
	} else {
		if os.IsNotExist(err) {
			// ok
		} else {
			panic(errors.Wrap(err, fmt.Sprintf("problem when checking file %s", outputFilePath)))
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

	schema, _ := cmd.Flags().GetString(flagSchema)
	schema = strings.TrimSpace(schema)

	toolName := "pg_dump"
	customToolName, _ := cmd.Flags().GetString(flagToolFile)
	customToolName = strings.TrimSpace(customToolName)
	if len(customToolName) > 0 {
		_, err = os.Stat(customToolName)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("custom tool file does not exists: %s", customToolName))
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

	dumpArgs := make([]string, 0)
	if len(host) > 0 {
		dumpArgs = append(dumpArgs, fmt.Sprintf("--host=%s", host))
	}
	if port > 0 {
		dumpArgs = append(dumpArgs, fmt.Sprintf("--port=%d", port))
	}
	if len(schema) > 0 {
		dumpArgs = append(dumpArgs, fmt.Sprintf("--schema=%s", schema))
	}
	dumpArgs = append(dumpArgs, "-Fc")
	if len(userName) > 0 {
		dumpArgs = append(dumpArgs, fmt.Sprintf("--username=%s", userName))
	}
	dumpArgs = append(dumpArgs, fmt.Sprintf("--file=%s", outputFilePath))
	dumpArgs = append(dumpArgs, dbName)

	fmt.Println("Output file:", outputFilePath)
	fmt.Println("Dump arguments:", strings.Join(dumpArgs, " "))
	fmt.Println("Begin backup", outputFileName, "at", utils.NowStr())

	exitCode := utils.LaunchApp(toolName, dumpArgs, envVars)
	if exitCode != 0 {
		fmt.Println("Failed to backup", outputFileName, "at", utils.NowStr())
		os.Exit(exitCode)
	}

	fmt.Println("Finished backup", outputFileName, "at", utils.NowStr())
}
