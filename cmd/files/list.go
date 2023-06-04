package files

import (
	"fmt"
	"github.com/EscanBE/go-ienumerable/goe"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	flagValueOrderByName = "name"
	flagValueOrderByDate = "date"
)

// ListingCommands registers a sub-tree of commands
func ListingCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Listing files and print them with absolute path",
		Args:  cobra.ExactArgs(0),
		Run:   listFiles,
	}

	cmd.PersistentFlags().Bool(
		constants.FLAG_SILENT,
		false,
		"when error occurs, process will exit immediately with exit code 0, no result will be printed",
	)

	cmd.PersistentFlags().StringArray(
		constants.FLAG_CONTAINS,
		make([]string, 0),
		fmt.Sprintf("print only files contains specific string, can be repeated multiple times, eg: --%s abc --%s def", constants.FLAG_CONTAINS, constants.FLAG_CONTAINS),
	)

	cmd.PersistentFlags().String(
		constants.FLAG_ORDER_BY,
		flagValueOrderByName,
		fmt.Sprintf("order files by %s or %s (creation date time)", flagValueOrderByName, flagValueOrderByDate),
	)

	cmd.PersistentFlags().Bool(
		constants.FLAG_DESC,
		false,
		"listing files by descending order",
	)

	cmd.PersistentFlags().Int(
		constants.FLAG_SKIP,
		0,
		"skip first N results",
	)

	cmd.PersistentFlags().Bool(
		constants.FLAG_DELETE,
		false,
		"files in result will be deleted, make sure permission setup correctly",
	)

	return cmd
}

func listFiles(cmd *cobra.Command, _ []string) {
	ignoreError, _ := cmd.Flags().GetBool(constants.FLAG_SILENT)

	defer func() {
		if ignoreError {
			_ = recover()
		}
	}()

	orderBy, _ := cmd.Flags().GetString(constants.FLAG_ORDER_BY)
	if len(orderBy) < 1 {
		panic(fmt.Errorf("missing value for mandatory flag --%s", constants.FLAG_ORDER_BY))
	}

	orderByDesc, _ := cmd.Flags().GetBool(constants.FLAG_DESC)

	deleteResultFiles, _ := cmd.Flags().GetBool(constants.FLAG_DELETE)

	skip, _ := cmd.Flags().GetInt(constants.FLAG_SKIP)
	if skip < 0 {
		panic(fmt.Errorf("negative value for flag --%s", constants.FLAG_SKIP))
	}

	containsString, _ := cmd.Flags().GetStringArray(constants.FLAG_CONTAINS)

	workingDir := utils.ReadFlagWorkingDir(cmd)

	files := goe.NewIEnumerable[string](listFilesWithinDir(workingDir)...)

	if len(containsString) > 0 {
		for _, part := range containsString {
			if len(part) > 0 {
				files = files.Where(func(file string) bool {
					_, fileName := path.Split(file)
					return strings.Contains(fileName, part)
				})
			}
		}
	}

	if !files.Any() {
		return
	}

	if orderBy == flagValueOrderByName {
		var orderedFiles goe.IOrderedEnumerable[string]
		if orderByDesc {
			orderedFiles = files.OrderDescending()
		} else {
			orderedFiles = files.Order()
		}
		files = orderedFiles.GetOrderedEnumerable()
	} else if orderBy == flagValueOrderByDate {
		var orderedFiles goe.IOrderedEnumerable[string]
		if orderByDesc {
			orderedFiles = files.OrderByDescending(func(file string) any {
				return statsFile(file).ModTime()
			}, nil)
		} else {
			orderedFiles = files.OrderBy(func(file string) any {
				return statsFile(file).ModTime()
			}, nil)
		}
		files = orderedFiles.GetOrderedEnumerable()
	} else {
		panic(fmt.Errorf("not supported value \"%s\" for flag --%s", orderBy, constants.FLAG_ORDER_BY))
	}

	files = files.Skip(skip)

	if !files.Any() {
		return
	}

	if deleteResultFiles {
		for _, file := range files.ToArray() {
			err := os.RemoveAll(file)
			if err != nil {
				panic(errors.Wrap(err, "failed to delete file"))
			}
		}
	}

	for _, file := range files.ToArray() {
		fmt.Println(file)
	}
}

func listFilesWithinDir(dir string) []string {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		panic(errors.Wrap(err, "failed to listing entries in the directory"))
	}

	result := make([]string, 0)
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		abs, err := filepath.Abs(path.Join(dir, dirEntry.Name()))
		if err != nil {
			panic(errors.Wrap(err, "failed to convert into absolute path"))
		}
		result = append(result, abs)
	}

	return libutils.GetUniqueElements(result...)
}

func statsFile(file string) os.FileInfo {
	fi, err := os.Stat(file)
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to stats file %s", file)))
	}
	return fi
}
