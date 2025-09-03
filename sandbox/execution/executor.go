// Package execution 沙盒执行模块
// 提供文件执行的沙盒环境支持
package execution

import (
	"fmt"
	"os"
	"path/filepath"
	"file-flow-service/config"
	"file-flow-service/utils/logger"
	"file-flow-service/sandbox/environments"

	"go.uber.org/zap"
)

// SandboxExecutor 沙盒执行器接口
type SandboxExecutor interface {
	// Init 初始化沙盒执行器
	Init(config *config.AppConfig, logger logger.Logger, envManager environments.EnvironmentManager) error
	
	// ExecuteTask 执行任务
	ExecuteTask(taskID, taskDir, cmd string, args []string, envType, envVersion string) error
	
	// CreateTaskDirectory 创建任务执行目录
	CreateTaskDirectory(taskID string) (string, error)
	
	// CleanupTaskDirectory 清理任务执行目录
	CleanupTaskDirectory(taskID string) error
}

// sandboxExecutor 沙盒执行器实现
type sandboxExecutor struct {
	config        *config.AppConfig
	logger        logger.Logger
	envManager    environments.EnvironmentManager
	taskDirectories map[string]string
}

// NewSandboxExecutor 创建沙盒执行器实例
func NewSandboxExecutor() SandboxExecutor {
	return &sandboxExecutor{
		taskDirectories: make(map[string]string),
	}
}

// Init 初始化沙盒执行器
// 参数: config 配置对象, logger 日志对象, envManager 环境管理器
// 返回: 错误信息
func (se *sandboxExecutor) Init(config *config.AppConfig, logger logger.Logger, envManager environments.EnvironmentManager) error {
	se.config = config
	se.logger = logger
	se.envManager = envManager
	
	// 创建执行目录
	executionPath := config.Sandbox.Execution.BasePath
	if err := os.MkdirAll(executionPath, 0755); err != nil {
		return fmt.Errorf("创建执行目录失败: %v", err)
	}
	
	// 创建任务目录
	tasksPath := config.Sandbox.Execution.TasksPath
	if err := os.MkdirAll(tasksPath, 0755); err != nil {
		return fmt.Errorf("创建任务目录失败: %v", err)
	}
	
	// 创建临时目录
	tempPath := config.Sandbox.Execution.TempPath
	if err := os.MkdirAll(tempPath, 0755); err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	
	// 创建锁目录
	locksPath := config.Sandbox.Execution.LocksPath
	if err := os.MkdirAll(locksPath, 0755); err != nil {
		return fmt.Errorf("创建锁目录失败: %v", err)
	}
	
	se.logger.Info("沙盒执行器初始化完成")
	return nil
}

// ExecuteTask 执行任务
// 参数: taskID 任务ID, taskDir 任务目录, cmd 命令, args 参数, envType 环境类型, envVersion 环境版本
// 返回: 错误信息
func (se *sandboxExecutor) ExecuteTask(taskID, taskDir, cmd string, args []string, envType, envVersion string) error {
	if se.config == nil {
		return fmt.Errorf("沙盒执行器未初始化")
	}
	
	se.logger.Info("开始执行任务", 
		zap.String("task_id", taskID),
		zap.String("command", cmd),
		zap.String("environment", envType+"-"+envVersion))
	
	// 这里应该实现实际的任务执行逻辑
	// 包括环境选择、沙盒隔离等
	
	// 示例实现：记录任务执行
	se.logger.Info("任务执行完成", zap.String("task_id", taskID))
	return nil
}

// CreateTaskDirectory 创建任务执行目录
// 参数: taskID 任务ID
// 返回: 目录路径，错误信息
func (se *sandboxExecutor) CreateTaskDirectory(taskID string) (string, error) {
	if se.config == nil {
		return "", fmt.Errorf("沙盒执行器未初始化")
	}
	
	// 创建任务目录
	taskDir := filepath.Join(se.config.Sandbox.Execution.TasksPath, taskID)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		return "", fmt.Errorf("创建任务目录失败: %v", err)
	}
	
	// 记录任务目录
	se.taskDirectories[taskID] = taskDir
	
	se.logger.Info("创建任务目录", 
		zap.String("task_id", taskID),
		zap.String("path", taskDir))
	
	return taskDir, nil
}

// CleanupTaskDirectory 清理任务执行目录
// 参数: taskID 任务ID
// 返回: 错误信息
func (se *sandboxExecutor) CleanupTaskDirectory(taskID string) error {
	if se.config == nil {
		return fmt.Errorf("沙盒执行器未初始化")
	}
	
	taskDir, exists := se.taskDirectories[taskID]
	if !exists {
		return fmt.Errorf("任务目录不存在: %s", taskID)
	}
	
	// 删除任务目录
	if err := os.RemoveAll(taskDir); err != nil {
		return fmt.Errorf("清理任务目录失败: %v", err)
	}
	
	// 从记录中移除
	delete(se.taskDirectories, taskID)
	
	se.logger.Info("清理任务目录", 
		zap.String("task_id", taskID),
		zap.String("path", taskDir))
	
	return nil
}
