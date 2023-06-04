package db

import (
	"fmt"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// BackupCommands registers a sub-tree of commands
func BackupCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Dump a database backup, dump file will be saved in working directory",
		Args:  cobra.ExactArgs(0),
		Run:   backupDatabase,
	}

	cmd.PersistentFlags().String(
		constants.FLAG_OUTPUT_FILE,
		fmt.Sprintf("db-%s.dump", time.Now().Format("2006-01-02")),
		"specify name of the output backup file, file name only, default has layout: db-yyyy-MM-dd.dump",
	)

	cmd.PersistentFlags().String(
		constants.FLAG_PASSWORD_FILE,
		"",
		"file path which store password of the user which will be used to backup the database",
	)

	cmd.PersistentFlags().String(
		constants.FLAG_TOOL_FILE,
		"",
		"absolute file path of the tool will be used to dump database, eg: /usr/bin/pg_dump for PostgreSQL",
	)

	return cmd
}

func backupDatabase(cmd *cobra.Command, _ []string) {
	outputFileName, _ := cmd.Flags().GetString(constants.FLAG_OUTPUT_FILE)
	outputFileName = strings.TrimSpace(outputFileName)
	if len(outputFileName) < 1 {
		panic(fmt.Errorf("missing value for mandatory flag --%s", constants.FLAG_OUTPUT_FILE))
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

	schema, _ := cmd.Flags().GetString(constants.FLAG_SCHEMA)
	schema = strings.TrimSpace(schema)

	toolName := "pg_dump"
	customToolName, _ := cmd.Flags().GetString(constants.FLAG_TOOL_FILE)
	customToolName = strings.TrimSpace(customToolName)
	if len(customToolName) > 0 {
		_, err = os.Stat(customToolName)
		if os.IsNotExist(err) {
			panic(fmt.Errorf("custom tool file does not exists: %s", customToolName))
		}

		toolName = customToolName
	}

	dbType, _ := cmd.Flags().GetString(constants.FLAG_TYPE)
	if len(dbType) < 1 {
		panic(fmt.Errorf("missing value for mandatory flag --%s", constants.FLAG_TYPE))
	}

	if dbType != dbTypePostgres {
		panic(fmt.Errorf("at this moment, only PostgreSQL db is supported"))
	}

	var envVar []string

	passwordFile, _ := cmd.Flags().GetString(constants.FLAG_PASSWORD_FILE)
	if len(passwordFile) < 1 {
		if len(strings.TrimSpace(os.Getenv(constants.ENV_PG_PASSWORD))) < 1 {
			panic(fmt.Errorf("missing password for user %s, either environment variable %s or flag --%s is required", userName, constants.ENV_PG_PASSWORD, constants.FLAG_PASSWORD_FILE))
		}

		envVar = os.Environ()
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
			panic(fmt.Errorf("password file is empty: %s", userName))
		}

		envVar = append(envVar, fmt.Sprintf("%s=%s", constants.ENV_PG_PASSWORD, pgPassword))
	}

	args := make([]string, 0)
	if len(host) > 0 {
		args = append(args, fmt.Sprintf("--host=%s", host))
	}
	if port > 0 {
		args = append(args, fmt.Sprintf("--port=%d", port))
	}
	if len(schema) > 0 {
		args = append(args, fmt.Sprintf("--schema=%s", schema))
	}
	args = append(args, "-Fc")
	if len(userName) > 0 {
		args = append(args, fmt.Sprintf("--username=%s", userName))
	}
	args = append(args, fmt.Sprintf("--file=%s", outputFilePath))
	args = append(args, dbName)

	fmt.Println("Output file:", outputFilePath)
	fmt.Println("Begin dump", outputFileName, "at", time.Now().Format("2006-Jan-02 15:04:05"))

	excCmd := exec.Command(toolName, args...)
	excCmd.Env = envVar
	stdout, err := excCmd.Output()
	if err != nil {
		fmt.Println("Failed to dump", outputFileName, "at", time.Now().Format("2006-Jan-02 15:04:05"))
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Finished dump", outputFileName, "at", time.Now().Format("2006-Jan-02 15:04:05"))
	if len(stdout) > 0 {
		fmt.Println(stdout)
	}
}
