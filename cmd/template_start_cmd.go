package cmd

/*
import (
	"fmt"
	libapp "github.com/EscanBE/go-lib/app"
	"github.com/EscanBE/house-keeper/constants"
	"github.com/spf13/cobra"
	"sync"
)

var (
	waitGroup sync.WaitGroup
)

// startCmd represents the start command, it launches the main business logic of this app
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start job",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := createApplicationContext()
		logger := ctx.Logger

		logger.Debug("Application starts")

		_, _ = ctx.SendTelegramLogMessage(fmt.Sprintf("[%s] Application Start", constants.APP_NAME))

		// Increase the waitGroup by one and decrease within trapExitSignal
		waitGroup.Add(1)

		// Register the function which should be executed upon exit.
		// After register, when you want to clean-up things before exit,
		// call libapp.ExecuteExitFunction(ctx) the same was as trapExitSignal method did
		libapp.RegisterExitFunction(func(params ...any) {
			// finalize
			defer waitGroup.Done()
			ctx.Logger.Debug("defer waitGroup::Done")

			// Legacy TODO Implements close connection, resources,... here to prevent resource leak
			if ctx.Bot != nil {
				ctx.Bot.StopReceivingUpdates()
			}

		})

		// Listen for and trap any OS signal to gracefully shutdown and exit
		trapExitSignal(ctx)

		// TODO implement business logic

		// end
		waitGroup.Wait()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
*/
