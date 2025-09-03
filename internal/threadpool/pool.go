// Package threadpool 线程池实现，提供并发执行能力
// 用于文件执行任务的并发处理
// 与 executor 模块紧密协作，提供任务执行的并发支持
package threadpool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"file-flow-service/utils/logger"
)

// Task 任务接口
// 定义任务必须实现的方法
type Task interface {
	Execute(ctx context.Context) error
	GetID() string
}

// ThreadPool 线程池实现
// 负责管理工作线程、任务队列和并发执行控制
type ThreadPool struct {
	maxWorkers  int
	taskQueue   chan Task
	workerWG    sync.WaitGroup
	taskTimeout time.Duration
	logger      logger.Logger
	running     bool
	mu          sync.RWMutex
	// 统计字段
	activeWorkers  int
	totalTasks     int64
	failedTasks    int64
	completedTasks int64
}

// NewThreadPool 创建线程池
// 参数：maxWorkers 最大工作线程数, taskTimeout 任务超时时间, logger 日志记录器
// 返回：线程池实例，错误信息
// 上下承接关系：验证参数合法性，初始化线程池结构体
func NewThreadPool(maxWorkers int, taskTimeout time.Duration, logger logger.Logger) (*ThreadPool, error) {
	if maxWorkers <= 0 {
		return nil, fmt.Errorf("最大工作线程数必须大于0")
	}
	return &ThreadPool{
		maxWorkers:  maxWorkers,
		taskQueue:   make(chan Task),
		taskTimeout: taskTimeout,
		logger:      logger,
		running:     false,
	}, nil
}

// Start 启动线程池
// 参数：无
// 返回：无
// 上下承接关系：启动指定数量的工作协程，开始处理任务队列
func (tp *ThreadPool) Start() {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	
	if tp.running {
		tp.logger.Warn("线程池已运行，无需重复启动")
		return
	}
	
	tp.logger.Info("线程池启动", zap.Int("max_workers", tp.maxWorkers))
	tp.running = true
	
	// 启动工作协程
	for i := 0; i < tp.maxWorkers; i++ {
		tp.workerWG.Add(1)
		go tp.worker()
	}
}

// worker 工作协程函数
// 参数：无
// 返回：无
// 上下承接关系：从任务队列中获取任务并执行，处理任务完成或超时情况
func (tp *ThreadPool) worker() {
	defer tp.workerWG.Done()
	
	for task := range tp.taskQueue {
		// 记录任务开始执行
		tp.mu.Lock()
		tp.activeWorkers++
		tp.totalTasks++
		tp.mu.Unlock()
		
		// 创建带超时的上下文
		ctx, cancel := context.WithTimeout(context.Background(), tp.taskTimeout)
		defer cancel()
		
		// 执行任务
		err := task.Execute(ctx)
		
		// 更新统计信息
		tp.mu.Lock()
		tp.activeWorkers--
		if err != nil {
			tp.failedTasks++
			tp.logger.Error("任务执行失败", zap.Error(err), zap.String("task_id", task.GetID()))
		} else {
			tp.completedTasks++
			tp.logger.Debug("任务执行完成", zap.String("task_id", task.GetID()))
		}
		tp.mu.Unlock()
		
		// 检查是否超时
		select {
		case <-ctx.Done():
			tp.logger.Warn("任务超时", zap.String("task_id", task.GetID()))
		default:
		}
	}
}

// Stop 停止线程池
// 参数：无
// 返回：无
// 上下承接关系：关闭任务队列，等待所有工作协程结束，清理资源
func (tp *ThreadPool) Stop() {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	
	if !tp.running {
		tp.logger.Warn("线程池未运行，无需停止")
		return
	}
	
	tp.logger.Info("线程池停止")
	
	// 关闭任务队列
	close(tp.taskQueue)
	
	// 等待所有工作协程结束
	tp.workerWG.Wait()
	
	tp.running = false
	tp.logger.Info("线程池已停止", 
		zap.Int64("total_tasks", tp.totalTasks),
		zap.Int64("failed_tasks", tp.failedTasks),
		zap.Int64("completed_tasks", tp.completedTasks))
}

// IsRunning 检查线程池是否正在运行
// 参数：无
// 返回：布尔值，表示线程池是否运行中
// 上下承接关系：检查线程池的运行状态
func (tp *ThreadPool) IsRunning() bool {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	return tp.running && tp.taskQueue != nil
}

// CancelTask 取消指定任务
// 参数：taskID 任务ID
// 返回：错误信息，如果任务不存在或取消失败则返回错误
// 上下承接关系：从任务队列中移除指定任务（注意：此实现为简化版本）
func (tp *ThreadPool) CancelTask(taskID string) error {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	
	// 简化实现：实际应用中需要更复杂的任务取消机制
	// 这里只是记录日志，实际取消需要更复杂的实现
	tp.logger.Debug("任务取消请求", zap.String("task_id", taskID))
	return nil
}

// Submit 提交任务到线程池
// 参数：task 要提交的任务
// 返回：无
// 上下承接关系：将任务放入任务队列，等待工作协程处理
func (tp *ThreadPool) Submit(task Task) {
	select {
	case tp.taskQueue <- task:
		tp.logger.Debug("任务已提交", zap.String("task_id", task.GetID()))
	default:
		tp.logger.Warn("任务队列已满，任务提交失败", zap.String("task_id", task.GetID()))
	}
}

// SetTaskTimeout 设置任务超时时间
// 参数：d 新的超时时间
// 返回：无
// 上下承接关系：更新线程池的任务超时配置
func (tp *ThreadPool) SetTaskTimeout(d time.Duration) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.taskTimeout = d
	tp.logger.Info("任务超时时间已更新", zap.Duration("timeout", d))
}

// ThreadPoolStats 线程池统计信息
type ThreadPoolStats struct {
	ActiveWorkers  int   `json:"active_workers"`
	QueueLength    int   `json:"queue_length"`
	TotalTasks     int64 `json:"total_tasks"`
	FailedTasks    int64 `json:"failed_tasks"`
	CompletedTasks int64 `json:"completed_tasks"`
}

// GetStats 获取线程池统计信息
// 参数：无
// 返回：线程池统计信息
// 上下承接关系：返回当前线程池的运行统计信息
func (tp *ThreadPool) GetStats() ThreadPoolStats {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	
	stats := ThreadPoolStats{
		ActiveWorkers:  tp.activeWorkers,
		QueueLength:    len(tp.taskQueue),
		TotalTasks:     tp.totalTasks,
		FailedTasks:    tp.failedTasks,
		CompletedTasks: tp.completedTasks,
	}
	
	return stats
}
