package cmd

import (
	"fmt"
	"github.com/EscanBE/go-app-name/config"
	"github.com/EscanBE/go-app-name/constants"
	"github.com/EscanBE/go-app-name/work"
	workertypes "github.com/EscanBE/go-app-name/work/types"
	libapp "github.com/EscanBE/go-lib/app"
	libcons "github.com/EscanBE/go-lib/constants"
	logtypes "github.com/EscanBE/go-lib/logging/types"
	libbot "github.com/EscanBE/go-lib/telegram/bot"
	libutils "github.com/EscanBE/go-lib/utils"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	waitGroup sync.WaitGroup
)

// startCmd represents the start command, it launches the main business logic of this app
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start job",
	Run: func(cmd *cobra.Command, args []string) {
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
		ctx := config.NewContext(conf, bot)
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

			// TODO Sample of template
			// Implements close connection, resources,... here to prevent resource leak

		})

		// Listen for and trap any OS signal to gracefully shutdown and exit
		trapExitSignal(ctx)

		// Create workers
		// Worker defines a job consumer that is responsible for getting assigned tasks and process business logic
		// to assign task to workers, use a channel

		// tasksChannel is a channel of WorkerTask which shared across all workers
		tasksChannel := make(chan workertypes.WorkerTask)

		// TODO Sample of template
		// TODO remove this sample
		sampleConcurrentMap := cmap.New[bool]()

		workers := make([]work.Worker, 0, conf.WorkerConfig.Count)
		for i := 0; i < conf.WorkerConfig.Count; i++ {
			workerWorkingCtx := &workertypes.WorkerWorkingContext{
				WorkerID: byte(i + 1),
				AppCtx:   *ctx,
				Logger:   ctx.Logger,
				RoCfg: workertypes.WorkerReadonlyConfig{
					SampleToken: conf.SecretConfig.SampleToken1,
				},
				RwCache: &workertypes.WorkerWritableCache{
					SampleConcurrentMap: &sampleConcurrentMap,
				},
				Queues: &workertypes.WorkerQueues{
					Tasks: tasksChannel,
				},
			}

			workers = append(workers, work.NewWorker(workerWorkingCtx))
		}

		// Start workers
		for _, worker := range workers {
			logger.Debug("Starting worker", "worker-id", worker.Id)
			go worker.Start() // parallel workers using go-routines
		}

		// TODO Sample of template
		// TODO remove this sample
		for i := 1; i <= 1_000; i++ {
			time.Sleep(time.Second)
			// enqueue a sample task to let worker process it
			tasksChannel <- workertypes.WorkerTask{
				TaskNo: i,
			}
		}

		// end
		waitGroup.Wait()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
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
