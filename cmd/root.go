package cmd

import (
	"github.com/EscanBE/go-app-name/cmd/utils"
	"github.com/EscanBE/go-app-name/constants"
	"github.com/spf13/cobra"
	"os"
)

// homeDir holds the home directory which was passed by flag `--home`, or default kinda `~/.binaryName`
var homeDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.BINARY_NAME,
	Short: constants.APP_DESC,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true    // hide the 'completion' subcommand
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true}) // hide the 'help' subcommand

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&homeDir, constants.FLAG_HOME, utils.GetDefaultHomeDirectory(), "Specify the home directory location instead of default")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
