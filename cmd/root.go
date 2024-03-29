package cmd

import (
	"github.com/EscanBE/house-keeper/cmd/config"
	"github.com/EscanBE/house-keeper/cmd/db"
	list "github.com/EscanBE/house-keeper/cmd/files"
	"github.com/EscanBE/house-keeper/cmd/gen"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   constants.BINARY_NAME,
	Short: constants.APP_DESC,
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
	rootCmd.AddCommand(list.Commands())
	rootCmd.AddCommand(db.Commands())
	rootCmd.AddCommand(config.Commands())
	rootCmd.AddCommand(gen.Commands())
}
