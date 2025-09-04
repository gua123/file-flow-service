package executor

import (
	"file-flow-service/config"
	"file-flow-service/internal/service/interfaces"
	"file-flow-service/utils/logger"
	"time"
	"file-flow-service/internal/threadpool"
)

type BaseExecutor struct {
	config    *config.AppConfig
	logger    logger.Logger
	threadpool *threadpool.ThreadPool
}

func NewExecutor(config *config.AppConfig, logger logger.Logger) *BaseExecutor {
	return &BaseExecutor{
		config:    config,
		logger:    logger,
		threadpool: threadpool.NewThreadPool(),
	}
}

func (e *BaseExecutor) GetPool() *threadpool.ThreadPool {
	return e.threadpool
}

func (e *BaseExecutor) Execute(task interfaces.TaskInterface) {
	defer func(start time.Time) {
		duration := time.Since(start)
		e.logger.Info("任务执行完成, task_id=" + task.GetID() + ", duration=" + duration.String())
	}(time.Now())

	e.threadpool.Submit(func() {
		e.logger.Info("任务提交到线程池, task_id=" + task.GetID())
		err := task.Execute()
		if err != nil {
			e.logger.Error("任务执行失败, task_id=" + task.GetID() + ", error=" + err.Error())
		}
	})
}

func (e *BaseExecutor) Stop() {
	e.logger.Info("停止执行器")
	e.threadpool.Stop()
}