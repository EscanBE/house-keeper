package cmd

import (
	"fmt"
	libapp "github.com/EscanBE/go-lib/app"
	libcons "github.com/EscanBE/go-lib/constants"
	logtypes "github.com/EscanBE/go-lib/logging/types"
	libbot "github.com/EscanBE/go-lib/telegram/bot"
	libutils "github.com/EscanBE/go-lib/utils"
	list "github.com/EscanBE/house-keeper/cmd/files"
	"github.com/EscanBE/house-keeper/config"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
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

	rootCmd.AddCommand(list.Commands())

	// rootCmd.PersistentFlags().StringVar(&homeDir, constants.FLAG_HOME, utils.GetDefaultHomeDirectory(), "Specify the home directory location instead of default")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func createApplicationContext() *config.ApplicationExecutionContext {
	fmt.Printf("Process id: %d\n", os.Getpid())

	conf, err := config.LoadConfig(homeDir)
	libutils.ExitIfErr(err, "unable to load configuration")

	// Output some options to console
	conf.PrintOptions()

	// Perform validation
	err = conf.Validate()
	libutils.ExitIfErr(err, "failed to validate configuration")

	// Initialize bot
	var bot *libbot.TelegramBot
	if len(conf.SecretConfig.TelegramToken) > 0 {
		bot, err = libbot.NewBot(conf.SecretConfig.TelegramToken)
		if err != nil {
			panic(errors.Wrap(err, "Failed to initialize Telegram bot"))
		}
		bot.EnableDebug(conf.Logging.Level == logtypes.LOG_LEVEL_DEBUG)
	}

	// Init execution context
	return config.NewContext(conf, bot)
}

// trapExitSignal traps the signal which being emitted when interrupting the application. Implement connection/resource close to prevent resource leaks
func trapExitSignal(ctx *config.ApplicationExecutionContext) {
	var sigCh = make(chan os.Signal)

	signal.Notify(sigCh, libcons.TrapExitSignals...)

	go func() {
		sig := <-sigCh
		ctx.Logger.Info(
			"caught signal; shutting down...",
			"os.signal", sig.String(),
		)

		libapp.ExecuteExitFunction()
	}()
}
