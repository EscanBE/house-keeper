package files

import (
	"fmt"
	"github.com/EscanBE/go-ienumerable/goe"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	flagContains   = "contains"
	flagRegex      = "regex"
	flagSilent     = "silent"
	flagSkip       = "skip"
	flagDeleteFile = "delete"
	flagOrderBy    = "order-by"
	flagDescending = "desc"

	flagValueOrderByName = "name"
	flagValueOrderByDate = "date"
)

// ListingCommands registers a sub-tree of commands
func ListingCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Listing files and print them with absolute path",
		Args:  cobra.NoArgs,
		Run:   listFiles,
	}

	utils.AddFlagWorkingDir(cmd)

	cmd.PersistentFlags().Bool(
		flagSilent,
		false,
		"when error occurs, process will exit immediately with exit code 0, no result will be printed",
	)

	cmd.PersistentFlags().StringArray(
		flagContains,
		make([]string, 0),
		fmt.Sprintf("print only files contains specific string in file name, can be repeated multiple times, eg: --%s abc --%s def", flagContains, flagContains),
	)

	cmd.PersistentFlags().String(
		flagRegex,
		"",
		fmt.Sprintf("print only files which name satisfy regex pattern, eg: --%s '^backup_.+' (must quoted by single quotes)", flagRegex),
	)

	cmd.PersistentFlags().String(
		flagOrderBy,
		flagValueOrderByName,
		fmt.Sprintf("order files by %s or %s (creation date time)", flagValueOrderByName, flagValueOrderByDate),
	)

	cmd.PersistentFlags().Bool(
		flagDescending,
		false,
		"listing files by descending order",
	)

	cmd.PersistentFlags().Int(
		flagSkip,
		0,
		"skip first N results",
	)

	cmd.PersistentFlags().Bool(
		flagDeleteFile,
		false,
		"files in result will be deleted, make sure permission setup correctly",
	)

	return cmd
}

func listFiles(cmd *cobra.Command, _ []string) {
	ignoreError, _ := cmd.Flags().GetBool(flagSilent)

	defer func() {
		if ignoreError {
			_ = recover()
		}
	}()

	orderBy, _ := cmd.Flags().GetString(flagOrderBy)
	if len(orderBy) < 1 {
		panic(fmt.Errorf("missing value for mandatory flag --%s", flagOrderBy))
	}

	orderByDesc, _ := cmd.Flags().GetBool(flagDescending)

	deleteResultFiles, _ := cmd.Flags().GetBool(flagDeleteFile)

	skip, _ := cmd.Flags().GetInt(flagSkip)
	if skip < 0 {
		panic(fmt.Errorf("negative value for flag --%s", flagSkip))
	}

	containsString, _ := cmd.Flags().GetStringArray(flagContains)

	regexPattern, _ := cmd.Flags().GetString(flagRegex)
	var regex *regexp.Regexp
	if len(regexPattern) > 0 {
		var errRegex error
		regex, errRegex = regexp.Compile(regexPattern)
		if errRegex != nil {
			panic(errors.Wrap(errRegex, fmt.Sprintf("failed to parse regex pattern provided by flag --%s (do you forgot quoted it with single quotes?)", flagRegex)))
		}
	}

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
	if regex != nil {
		files = files.Where(func(file string) bool {
			_, fileName := path.Split(file)
			return regex.MatchString(fileName)
		})
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
		panic(fmt.Errorf("not supported value \"%s\" for flag --%s", orderBy, flagOrderBy))
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
