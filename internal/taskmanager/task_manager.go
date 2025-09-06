package taskmanager

import (
	"time"
	"file-flow-service/config"
	"file-flow-service/utils/logger"
	"sync"
	"file-flow-service/internal/threadpool"
	"file-flow-service/internal/service/interfaces"
	"file-flow-service/database"
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

	// 获取任务数据
	task, err := database.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return nil
	}
	
	// 更新状态
	task.Status = status
	// 保存到数据库
	if err := database.UpdateTask(task); err != nil {
		return err
	}
	return nil
}

func (tm *taskManager) GetAllTasks() ([]*interfaces.TaskInterface, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 获取所有任务
	dbTasks, err := database.GetTasks()
	if err != nil {
		return nil, err
	}
	
	// 转换为TaskInterface
	var tasks []*interfaces.TaskInterface
	for _, dbTask := range dbTasks {
		// 正确转换为接口：使用指针
		task := interfaces.TaskInterface(&dbTask)
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
	// 转换为数据库任务
	dbTask := database.Task{
		ID:          task.GetID(),
		Name:        task.GetName(),
		Status:      task.GetStatus(),
		Creator:     task.GetCreator(),
		CreatedAt:   task.GetCreatedAt(),
		AssignedTo:  task.GetAssignedTo(),
		Description: task.GetDescription(),
		ResultPath:  task.GetResultPath(),
		Progress:    task.GetProgress(),
	}
	
	if err := database.CreateTask(&dbTask); err != nil {
		return err
	}
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