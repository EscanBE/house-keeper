package work

import (
	"context"
	"fmt"
	"github.com/EscanBE/go-app-name/database"
	workertypes "github.com/EscanBE/go-app-name/work/types"
	libapp "github.com/EscanBE/go-lib/app"
	"github.com/EscanBE/go-lib/logging"
	libutils "github.com/EscanBE/go-lib/utils"
)

// Worker represents for a worker, itself holds things needed for doing business logic, especially its own `WorkerWorkingContext`
type Worker struct {
	Id     byte
	Ctx    *workertypes.WorkerWorkingContext
	logger logging.Logger
}

// NewWorker creates new worker and inject needed information
func NewWorker(wCtx *workertypes.WorkerWorkingContext) Worker {
	return Worker{
		Id:     wCtx.WorkerID,
		Ctx:    wCtx,
		logger: wCtx.AppCtx.Logger,
	}
}

// BeginDatabaseTransaction uses the initialized DB obj to start new transaction
// By using tx methodology, it would help data remains consistency
func (w Worker) BeginDatabaseTransaction() (database.DbTransaction, error) {
	return w.Ctx.AppCtx.Database.BeginDatabaseTransaction(context.Background())
}

// Start performs business logic of worker
func (w Worker) Start() {
	defer libapp.TryRecoverAndExecuteExitFunctionIfRecovered(w.logger)

	// TODO Sample of template
	// TODO implement logic here

	// TODO Sample of template
	// TODO remove this sample
	for task := range w.Ctx.Queues.Tasks {
		// TODO implement task processing logic here
		fmt.Printf("Worker %d is processing task %d\n", w.Id, task.TaskNo)

		dbTx, err := w.BeginDatabaseTransaction()
		if err != nil {
			w.logger.Error("failed to process task", "task", task.TaskNo, "worker", w.Id, "error", err.Error())
			// w.Ctx.Queues.Tasks <- task // Re-Enqueue if needed
			continue
		}

		errTask := w.doTask(task)
		if errTask != nil {
			errRb := dbTx.RollbackTransaction()
			if errRb != nil {
				w.logger.Error(
					"failed to rollback tx after processing task failed",
					"task", task.TaskNo,
					"worker", w.Id,
					"error", errTask.Error(),
					"rollback-error", errRb.Error(),
				)
				// w.Ctx.Queues.Tasks <- task // Re-Enqueue if needed
			}
		} else {
			errCm := dbTx.CommitTransaction()
			if errCm != nil {
				w.logger.Error(
					"failed to commit tx after processing task successfully",
					"task", task.TaskNo,
					"worker", w.Id,
					"error", errTask.Error(),
					"commit-error", errCm.Error(),
				)
				// w.Ctx.Queues.Tasks <- task // Re-Enqueue if needed
			}
		}
	}
}

// TODO Sample of template
// TODO remove this sample
func (w Worker) doTask(workertypes.WorkerTask) error {
	nowS := libutils.NowS()
	var err error
	if nowS%2 == 0 { // simulate error
		err = fmt.Errorf("sample error message")
	}
	return libutils.NilOrWrapIfError(err, "failed to process task")
}
