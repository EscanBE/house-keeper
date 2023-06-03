package types

import (
	"github.com/EscanBE/go-app-name/config"
	"github.com/EscanBE/go-lib/logging"
	"github.com/EscanBE/go-lib/telegram/bot"
	cmap "github.com/orcaman/concurrent-map/v2"
)

// WorkerWorkingContext hold the working context for each of the worker.
// In here, we can save identity, config, caches,...
type WorkerWorkingContext struct {
	WorkerID byte
	AppCtx   config.ApplicationExecutionContext
	Logger   logging.Logger
	RoCfg    WorkerReadonlyConfig
	RwCache  *WorkerWritableCache
	Queues   *WorkerQueues
}

// WorkerReadonlyConfig contains readonly configuration options
type WorkerReadonlyConfig struct {
	// TODO Sample of template
	// TODO remove this sample
	SampleToken string
}

// WorkerWritableCache contains caches, resources shared across workers, or local init & use, depends on implementation
type WorkerWritableCache struct {
	// TODO Sample of template
	// TODO remove this sample
	SampleConcurrentMap *cmap.ConcurrentMap[bool]
}

// WorkerTask to be used to put into a channel which shares across all workers.
// Each worker will take this WorkerTask from the channel and execute it based on input arguments within
type WorkerTask struct {
	TaskNo int
	// TODO Sample of template
	// TODO implement parameters for task here
}

// WorkerQueues contains queues shared across workers
type WorkerQueues struct {
	// Tasks is channel of worker's tasks, and all workers are sharing/using the same channel to get task from.
	// To assign task to a worker, put a WorkerTask with desired arguments into this channel
	Tasks chan WorkerTask
}

// GetTelegramBot returns bot.TelegramBot instance
func (wc WorkerWorkingContext) GetTelegramBot() *bot.TelegramBot {
	return wc.AppCtx.Bot
}
