package taskmanager

import (
	"time"
	"file-flow-service/config"
	"file-flow-service/utils/logger"
	"sync"
	"file-flow-service/internal/threadpool"
	"file-flow-service/internal/service/interfaces"
)

type TaskManager interface {
	Start()
	Stop()
	GetRunningTaskCount() int
	GetTaskByID(taskID string) (*interfaces.TaskInterface, error)
	UpdateTask(taskID string, status string) error
	GetAllTasks() ([]*interfaces.TaskInterface, error)
	SubmitTask(task interfaces.TaskInterface) error
	CancelTask(taskID string) error
	GetThreadPoolStats() (*threadpool.ThreadPoolStats, error)
}

type taskManager struct {
	config          *config.AppConfig
	threadpool      *threadpool.ThreadPool
	logger          logger.Logger
	tasks           map[string]interfaces.TaskInterface
	mu              sync.RWMutex
	runningTasks    int
	totalTasks      int
	activeTaskCount int
}

func NewTaskManager(config *config.AppConfig, threadpool *threadpool.ThreadPool, logger logger.Logger) TaskManager {
	return &taskManager{
		config:   config,
		threadpool: threadpool,
		logger:   logger,
		tasks:    make(map[string]interfaces.TaskInterface),
	}
}

func (tm *taskManager) Start() {
	tm.logger.Info("Task manager started")
}

func (tm *taskManager) Stop() {
	tm.logger.Info("Task manager stopped")
	tm.threadpool.Stop()
}

func (tm *taskManager) GetRunningTaskCount() int {
	return tm.activeTaskCount
}

func (tm *taskManager) GetTaskByID(taskID string) (*interfaces.TaskInterface, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, nil
	}
	return &task, nil
}

func (tm *taskManager) UpdateTask(taskID string, status string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return nil
	}
	task.SetStatus(status)
	return nil
}

func (tm *taskManager) GetAllTasks() ([]*interfaces.TaskInterface, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tasks := make([]*interfaces.TaskInterface, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (tm *taskManager) SubmitTask(task interfaces.TaskInterface) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.tasks[task.GetID()]; exists {
		return nil
	}

	tm.totalTasks++
	tm.tasks[task.GetID()] = task
	tm.activeTaskCount++

	tm.logger.Info("任务已提交到执行器: " + task.GetID())

	tm.threadpool.Submit(func() {
		task.SetStatus("running")
		startTime := time.Now().Unix()
		task.SetStartTime(startTime)
		task.Execute()
		task.SetStatus("completed")
		finishTime := time.Now().Unix()
		duration := finishTime - startTime
		task.SetDuration(duration)
		task.SetFinishedAt(finishTime)
		tm.activeTaskCount--
	})

	return nil
}

func (tm *taskManager) CancelTask(taskID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return nil
	}

	task.SetStatus("cancelled")
	tm.UpdateTask(taskID, "cancelled")
	tm.activeTaskCount--
	return nil
}

func (tm *taskManager) GetThreadPoolStats() (*threadpool.ThreadPoolStats, error) {
	poolStats := tm.threadpool.GetStats()
	return &threadpool.ThreadPoolStats{
		TotalTasks:     poolStats.TotalTasks,
		ActiveTasks:    poolStats.ActiveTasks,
		CompletedTasks: poolStats.CompletedTasks,
	}, nil
}