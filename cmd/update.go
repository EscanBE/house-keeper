package cmd

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/EscanBE/house-keeper/cmd/utils"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/spf13/cobra"
	"os"
	"regexp"
)

// updateCmd represents the update command, it updates the house-keeper binary to latest or specified version
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Args:    cobra.RangeArgs(0, 1),
	Short:   fmt.Sprintf("Update %s binary version", constants.BINARY_NAME),
	Run: func(cmd *cobra.Command, args []string) {
		var version string
		if len(args) == 0 {
			version = "latest"
		} else {
			version = args[0]

			if !regexp.MustCompile("^v?\\d+\\.\\d+\\.\\d+$").MatchString(version) {
				libutils.PrintfStdErr("version format of %s is malformed\n", version)
				os.Exit(1)
			}
		}

		ec := utils.LaunchApp("go", []string{"install", "-v", fmt.Sprintf("github.com/EscanBE/house-keeper/cmd/hkd@%s", version)}, nil, true)
		if ec != 0 {
			libutils.PrintfStdErr("Exited with status code: %d\n", ec)
			os.Exit(ec)
		}

		fmt.Println("Updated", constants.BINARY_NAME)
		fmt.Printf("%s => ", constants.VERSION)
		_ = utils.LaunchApp(constants.BINARY_NAME, []string{"version"}, nil, true)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
