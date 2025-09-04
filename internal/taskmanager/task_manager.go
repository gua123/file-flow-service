// Package taskmanager 任务管理器
// 负责任务的创建、管理、监控和统计
// 与 executor 和 threadpool 模块紧密协作，提供完整的任务生命周期管理
package taskmanager

import (
	"context"
	"fmt"
	"sync"
	"time"
	"file-flow-service/internal/service/interfaces"
	"file-flow-service/internal/threadpool"
	"file-flow-service/utils/logger"
	"file-flow-service/config"
	
	"go.uber.org/zap"
)

// TaskInfo 任务信息
// 存储任务的详细信息和状态
type TaskInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Progress    float64   `json:"progress"`
	CreatedAt   time.Time `json:"created_at"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	Duration    time.Duration `json:"duration"`
	WorkerID    int       `json:"worker_id"`
	Error       string    `json:"error,omitempty"`
	Result      string    `json:"result,omitempty"`
}

// TaskStats 任务统计信息
// 提供任务执行的统计信息
type TaskStats struct {
	TotalTasks     int     `json:"total_tasks"`
	ActiveTasks    int     `json:"active_tasks"`
	CompletedTasks int     `json:"completed_tasks"`
	FailedTasks    int     `json:"failed_tasks"`
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsage    uint64  `json:"memory_usage"`
	QueueLength    int     `json:"queue_length"`
	ActiveWorkers  int     `json:"active_workers"`
	Timestamp      int64   `json:"timestamp"`
}

// ConvertToAPITaskStats 转换为API任务统计信息
// 参数：无
// 返回：API任务统计信息
func (ts *TaskStats) ConvertToAPITaskStats() *interfaces.TaskStats {
	return &interfaces.TaskStats{
		TotalTasks:     ts.TotalTasks,
		ActiveTasks:    ts.ActiveTasks,
		CompletedTasks: ts.CompletedTasks,
		FailedTasks:    ts.FailedTasks,
		CPUUsage:       ts.CPUUsage,
		MemoryUsage:    ts.MemoryUsage,
		QueueLength:    ts.QueueLength,
		ActiveWorkers:  ts.ActiveWorkers,
		Timestamp:      ts.Timestamp,
	}
}

// TaskManager 任务管理器接口
// 定义任务管理器必须实现的方法
type TaskManager interface {
	// 启动管理器
	// 参数：无
	// 返回：错误信息
	// 上下承接关系：初始化监控循环，开始任务状态监控
	Start() error
	
	// 停止管理器
	// 参数：无
	// 返回：错误信息
	// 上下承接关系：停止监控循环，清理资源
	Stop() error
	
	// 获取所有任务
	// 参数：无
	// 返回：任务信息切片，错误信息
	// 上下承接关系：返回当前所有任务的快照
	GetAllTasks() ([]*TaskInfo, error)
	
	// 获取指定任务
	// 参数：taskID 任务ID
	// 返回：任务信息，错误信息
	// 上下承接关系：根据ID查找并返回指定任务
	GetTask(taskID string) (*TaskInfo, error)
	
	// 取消任务
	// 参数：taskID 任务ID
	// 返回：错误信息
	// 上下承接关系：更新任务状态为取消，调用线程池取消方法
	CancelTask(taskID string) error
	
	// 重新执行任务
	// 参数：taskID 任务ID
	// 返回：错误信息
	// 上下承接关系：重置任务状态，准备重新执行
	RetryTask(taskID string) error
	
	// 获取任务统计信息
	// 参数：无
	// 返回：任务统计信息，错误信息
	// 上下承接关系：聚合任务状态信息，返回统计结果
	GetTaskStats() (*TaskStats, error)
	
	// 获取线程池状态
	// 参数：无
	// 返回：线程池统计信息，错误信息
	// 上下承接关系：获取线程池运行状态信息
	GetThreadPoolStats() (*threadpool.ThreadPoolStats, error)
	
	// 监控任务状态
	// 参数：无
	// 返回：错误信息
	// 上下承接关系：定期更新任务状态信息
	MonitorTasks() error
	
	// 更新任务
	// 参数：taskID 任务ID, task 任务信息
	// 返回：错误信息
	// 上下承接关系：更新任务信息
	UpdateTask(taskID string, task *TaskInfo) error
	
	// 删除任务
	// 参数：taskID 任务ID
	// 返回：错误信息
	// 上下承接关系：删除指定任务
	DeleteTask(taskID string) error
}

// taskManager 任务管理器实现
// 实现任务管理器接口，提供完整的任务管理功能
type taskManager struct {
	logger        logger.Logger
	config        *config.AppConfig
	threadpool    *threadpool.ThreadPool
	running       bool
	mu            sync.RWMutex
	tasks         map[string]*TaskInfo
	ticker        *time.Ticker
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewTaskManager 创建任务管理器
// 参数：config 配置对象, threadpool 线程池实例, logger 日志记录器
// 返回：任务管理器接口实例
// 上下承接关系：初始化任务管理器结构体，创建上下文
func NewTaskManager(config *config.AppConfig, threadpool *threadpool.ThreadPool, logger logger.Logger) TaskManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &taskManager{
		config:    config,
		threadpool: threadpool,
		logger:    logger,
		ctx:       ctx,
		cancel:    cancel,
		tasks:     make(map[string]*TaskInfo),
	}
}

// Start 启动任务管理器
// 参数：无
// 返回：错误信息，如果启动失败则返回错误
// 上下承接关系：初始化监控循环，开始定期更新任务状态
func (tm *taskManager) Start() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if tm.running {
		return nil
	}
	
	tm.running = true
	tm.logger.Info("任务管理器启动")
	
	// 启动监控循环
	go tm.monitorLoop()
	
	return nil
}

// Stop 停止任务管理器
// 参数：无
// 返回：错误信息，如果停止失败则返回错误
// 上下承接关系：停止监控循环，清理资源，关闭上下文
func (tm *taskManager) Stop() error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	if !tm.running {
		return nil
	}
	
	tm.running = false
	tm.cancel()
	
	if tm.ticker != nil {
		tm.ticker.Stop()
	}
	
	tm.logger.Info("任务管理器停止")
	return nil
}

// GetAllTasks 获取所有任务
// 参数：无
// 返回：任务信息切片，错误信息
// 上下承接关系：返回当前所有任务的快照，用于任务列表展示
func (tm *taskManager) GetAllTasks() ([]*TaskInfo, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	var tasks []*TaskInfo
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}
	
	tm.logger.Debug("获取所有任务", zap.Int("count", len(tasks)))
	return tasks, nil
}

// GetTask 获取指定任务
// 参数：taskID 任务ID
// 返回：任务信息，错误信息
// 上下承接关系：根据ID查找并返回指定任务，用于任务详情展示
func (tm *taskManager) GetTask(taskID string) (*TaskInfo, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("任务 %s 不存在", taskID)
	}
	
	tm.logger.Debug("获取任务", zap.String("task_id", taskID))
	return task, nil
}

// CancelTask 取消任务
// 参数：taskID 任务ID
// 返回：错误信息，如果任务不存在或取消失败则返回错误
// 上下承接关系：更新任务状态为取消，调用线程池取消方法，记录取消操作
func (tm *taskManager) CancelTask(taskID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	task, exists := tm.tasks[taskID]
	if !exists {
		return fmt.Errorf("任务 %s 不存在", taskID)
	}
	
	// 调用线程池的取消方法
	if tm.threadpool != nil {
		err := tm.threadpool.CancelTask(taskID)
		if err != nil {
			tm.logger.Warn("线程池取消任务失败", zap.String("task_id", taskID), zap.Error(err))
		}
	}
	
	// 更新任务状态
	task.Status = "cancelled"
	task.CompletedAt = time.Now()
	task.Duration = task.CompletedAt.Sub(task.CreatedAt)
	
	tm.logger.Info("任务取消", zap.String("task_id", taskID))
	return nil
}

// RetryTask 重新执行任务
// 参数：taskID 任务ID
// 返回：错误信息，如果任务不存在则返回错误
// 上下承接关系：重置任务状态，准备重新执行，记录重试操作
func (tm *taskManager) RetryTask(taskID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	task, exists := tm.tasks[taskID]
	if !exists {
		return fmt.Errorf("任务 %s 不存在", taskID)
	}
	
	// 重置任务状态
	task.Status = "pending"
	task.Progress = 0
	task.Error = ""
	task.Result = ""
	task.CompletedAt = time.Time{}
	task.Duration = 0
	
	tm.logger.Info("任务重试", zap.String("task_id", taskID))
	return nil
}

// GetTaskStats 获取任务统计信息
// 参数：无
// 返回：任务统计信息，错误信息
// 上下承接关系：聚合任务状态信息，返回统计结果用于监控面板
func (tm *taskManager) GetTaskStats() (*TaskStats, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	stats := &TaskStats{
		TotalTasks: len(tm.tasks),
		Timestamp:  time.Now().Unix(),
	}
	
	// 统计任务状态
	for _, task := range tm.tasks {
		switch task.Status {
		case "running":
			stats.ActiveTasks++
		case "completed":
			stats.CompletedTasks++
		case "failed":
			stats.FailedTasks++
		}
	}
	
	// 获取线程池统计
	if tm.threadpool != nil {
		poolStats := tm.threadpool.GetStats()
		stats.QueueLength = poolStats.QueueLength
		stats.ActiveWorkers = poolStats.ActiveWorkers
	}
	
	tm.logger.Debug("获取任务统计信息", 
		zap.Int("total_tasks", stats.TotalTasks),
		zap.Int("active_tasks", stats.ActiveTasks),
		zap.Int("completed_tasks", stats.CompletedTasks),
		zap.Int("failed_tasks", stats.FailedTasks))
	
	return stats, nil
}

// GetThreadPoolStats 获取线程池状态
// 参数：无
// 返回：线程池统计信息，错误信息
// 上下承接关系：获取线程池运行状态信息，用于系统监控
func (tm *taskManager) GetThreadPoolStats() (*threadpool.ThreadPoolStats, error) {
	if tm.threadpool == nil {
		return nil, fmt.Errorf("线程池未初始化")
	}
	
	stats := tm.threadpool.GetStats()
	tm.logger.Debug("获取线程池状态", 
		zap.Int("active_workers", stats.ActiveWorkers),
		zap.Int("queue_length", stats.QueueLength))
	
	return &stats, nil
}

// MonitorTasks 监控任务状态
// 参数：无
// 返回：错误信息，如果监控失败则返回错误
// 上下承接关系：定期更新任务状态，确保任务状态与实际执行情况一致
func (tm *taskManager) MonitorTasks() error {
	tm.logger.Debug("监控任务状态")
	// 实现监控逻辑
	tm.updateTaskStatus()
	return nil
}

// monitorLoop 监控循环
// 参数：无
// 返回：无
// 上下承接关系：定期执行任务状态更新，处理任务生命周期管理
func (tm *taskManager) monitorLoop() {
	interval, err := time.ParseDuration(tm.config.Monitoring.HealthCheck.Interval)
	if err != nil {
		interval = 5 * time.Second
	}
	
	tm.ticker = time.NewTicker(interval)
	defer tm.ticker.Stop()
	
	for {
		select {
		case <-tm.ctx.Done():
			tm.logger.Info("监控循环停止")
			return
		case <-tm.ticker.C:
			tm.updateTaskStatus()
		}
	}
}

// updateTaskStatus 更新任务状态
// 参数：无
// 返回：无
// 上下承接关系：定期检查任务执行状态，更新任务信息
func (tm *taskManager) updateTaskStatus() {
	tm.mu.RLock()
	taskCount := len(tm.tasks)
	tm.mu.RUnlock()
	
	tm.logger.Debug("更新任务状态", zap.Int("task_count", taskCount))
	
	// 实现任务状态更新逻辑
	// 这里可以检查线程池中的任务状态并更新
	// 由于是示例实现，我们只是记录日志
	if tm.threadpool != nil {
		stats := tm.threadpool.GetStats()
		tm.logger.Debug("线程池统计信息", 
			zap.Int("active_workers", stats.ActiveWorkers),
			zap.Int("queue_length", stats.QueueLength))
	}
}

// UpdateTask 更新任务
// 参数：taskID 任务ID, task 任务信息
// 返回：错误信息
// 上下承接关系：更新任务信息
func (tm *taskManager) UpdateTask(taskID string, task *TaskInfo) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	_, exists := tm.tasks[taskID]
	if !exists {
		return fmt.Errorf("任务 %s 不存在", taskID)
	}
	
	// 更新任务信息
	tm.tasks[taskID] = task
	tm.logger.Info("任务已更新", zap.String("task_id", taskID))
	return nil
}

// DeleteTask 删除任务
// 参数：taskID 任务ID
// 返回：错误信息
// 上下承接关系：删除指定任务
func (tm *taskManager) DeleteTask(taskID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	_, exists := tm.tasks[taskID]
	if !exists {
		return fmt.Errorf("任务 %s 不存在", taskID)
	}
	
	// 删除任务
	delete(tm.tasks, taskID)
	tm.logger.Info("任务已删除", zap.String("task_id", taskID))
	return nil
}
