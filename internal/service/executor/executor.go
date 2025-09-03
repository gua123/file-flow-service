// Package executor 实现任务执行器的核心逻辑，包括任务提交、执行和状态管理。
// 该包负责管理线程池、任务调度和健康检查，确保任务高效执行。
// 与 threadpool 模块紧密协作，提供任务执行能力
package executor

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"file-flow-service/config"
	"file-flow-service/internal/service/api"
	"file-flow-service/internal/threadpool"
	"file-flow-service/utils/logger"

	"go.uber.org/zap"
)

// BaseExecutor 任务执行器基础实现
// 负责管理线程池、任务提交和执行状态控制
type BaseExecutor struct {
	config    *config.AppConfig
	pool      *threadpool.ThreadPool
	mu        sync.RWMutex
	Status    string
	Logger    logger.Logger
	startTime time.Time // 启动时间字段
}

// GetPool 获取线程池实例
// 参数：无
// 返回：线程池实例
func (e *BaseExecutor) GetPool() *threadpool.ThreadPool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.pool
}

// NewExecutor 创建任务执行器实例
// 参数：cfg 配置对象，logger 日志记录器
// 返回：任务执行器实例
func NewExecutor(cfg *config.AppConfig, log logger.Logger) *BaseExecutor {
	return &BaseExecutor{
		config: cfg,
		Logger: log,
		Status: "stopped",
	}
}

// Start 启动执行器，初始化线程池并设置状态为运行中。
// 参数：无
// 返回：错误信息，如果启动失败则返回错误
// 上下承接关系：初始化线程池，设置状态为运行中，记录启动时间
func (e *BaseExecutor) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.pool != nil && e.pool.IsRunning() {
		return fmt.Errorf("service already running")
	}

	timeout, err := time.ParseDuration(e.config.Threadpool.TaskTimeout)
	if err != nil {
		return fmt.Errorf("解析任务超时时间失败: %v", err)
	}
	
	pool, err := threadpool.NewThreadPool(
		e.config.Threadpool.MaxWorkers,
		timeout,
		e.Logger,
	)
	if err != nil {
		return fmt.Errorf("创建线程池失败: %v", err)
	}
	
	e.pool = pool
	e.pool.Start()
	e.Status = "running"
	e.startTime = time.Now() // 初始化启动时间

	e.Logger.Info(
		fmt.Sprintf("服务 %s 启动", e.config.App.Name),
		zap.Int("max_workers", e.config.Threadpool.MaxWorkers),
		zap.Duration("task_timeout", timeout),
	)
	return nil
}

// Stop 停止执行器，关闭线程池并重置状态。
// 参数：无
// 返回：错误信息，如果停止失败则返回错误
// 上下承接关系：关闭线程池，重置状态，记录运行时长
func (e *BaseExecutor) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.pool == nil || !e.pool.IsRunning() {
		return fmt.Errorf("服务未运行")
	}

	e.pool.Stop()
	e.pool = nil
	e.Status = "stopped"

	duration := time.Since(e.startTime)
	e.Logger.Info(
		fmt.Sprintf("服务 %s 停止", e.config.App.Name),
		zap.Duration("duration", duration),
	)
	return nil
}

// HealthCheck 检查执行器的健康状态，返回是否运行中。
// 参数：无
// 返回：布尔值，表示执行器是否健康运行
// 上下承接关系：检查线程池状态，判断执行器是否正常运行
func (e *BaseExecutor) HealthCheck() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.pool != nil && e.pool.IsRunning()
}

// SubmitTask 提交任务到线程池，返回错误。
// 参数：task 要提交的任务
// 返回：错误信息，如果提交失败则返回错误
// 上下承接关系：将任务提交到线程池队列，由工作协程执行
func (e *BaseExecutor) SubmitTask(task Task) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.pool == nil || !e.pool.IsRunning() {
		return fmt.Errorf("服务未运行，无法提交任务")
	}

	e.pool.Submit(task)
	e.Logger.Debug("任务已提交到线程池", zap.String("task_id", task.GetID()))
	return nil
}

// Task 任务接口
// 定义任务必须实现的方法
type Task interface {
	Execute(ctx context.Context) error
	GetID() string
	GetCmd() string
	GetArgs() []string
}

// RunTask 提交任务到执行器，使用API任务结构体。
// 参数：task API任务结构体
// 返回：错误信息，如果提交失败则返回错误
// 上下承接关系：将API任务转换为内部任务结构体并提交到线程池
func (e *BaseExecutor) RunTask(task api.Task) error {
	scriptTask := &ScriptTask{
		ID:   task.ID,
		Dir:  task.Dir,
		Cmd:  task.Cmd,
		Args: task.Args,
	}
	
	e.Logger.Info("提交任务到执行器", zap.String("task_id", task.ID))
	return e.SubmitTask(scriptTask)
}

// ScriptTask 脚本任务实现
// 实现Task接口，用于执行系统命令
type ScriptTask struct {
	ID   string
	Dir  string
	Cmd  string
	Args []string
}

// Execute 执行命令，返回错误。
// 参数：ctx 上下文，用于控制执行超时
// 返回：错误信息，如果执行失败则返回错误
// 上下承接关系：使用exec包执行系统命令，支持上下文超时控制
func (t *ScriptTask) Execute(ctx context.Context) error {
	// 创建命令
	cmd := exec.CommandContext(ctx, t.Cmd, t.Args...)
	cmd.Dir = t.Dir
	
	// 执行命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动命令失败: %v", err)
	}
	
	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("命令执行失败: %v", err)
	}
	
	return nil
}

// GetID 获取任务ID
// 参数：无
// 返回：任务ID字符串
func (t *ScriptTask) GetID() string { return t.ID }

// GetCmd 获取命令
// 参数：无
// 返回：命令字符串
func (t *ScriptTask) GetCmd() string { return t.Cmd }

// GetArgs 获取参数列表
// 参数：无
// 返回：参数字符串切片
func (t *ScriptTask) GetArgs() []string { return t.Args }
