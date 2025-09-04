package service

import (
	"fmt"
	"time"
	"file-flow-service/internal/shutdown"
	"file-flow-service/internal/service/executor"
	"file-flow-service/internal/service/interfaces"
	"file-flow-service/internal/taskmanager"
	"file-flow-service/internal/restart"
	"file-flow-service/internal/service/monitor"
	"file-flow-service/utils/logger"
	"file-flow-service/config"
	"mime/multipart"
	"go.uber.org/zap"
)

type Service struct {
	logger     logger.Logger
	config     *config.AppConfig
	shutdown   *shutdown.ShutdownManager
	monitor    *monitor.MonitorImpl
	taskManager taskmanager.TaskManager
	restartManager *restart.RestartManager
	executor   *executor.Executor
}

// NewService creates a new Service instance
func NewService(config *config.AppConfig, logger logger.Logger, shutdown *shutdown.ShutdownManager, monitor *monitor.MonitorImpl, taskManager taskmanager.TaskManager, restartManager *restart.RestartManager, executor *executor.Executor) *Service {
	return &Service{
		logger:         logger,
		config:         config,
		shutdown:       shutdown,
		monitor:        monitor,
		taskManager:    taskManager,
		restartManager: restartManager,
		executor:       executor,
	}
}

// UpdateTask updates task information
func (s *Service) UpdateTask(taskID string, req interfaces.UpdateTaskRequest) error {
	if s.taskManager == nil {
		return fmt.Errorf("task manager not initialized")
	}
	
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] UpdateTask: params=%s, duration=%dms", zap.String("task_id", taskID), zap.Int64("duration", duration))
	}()
	
	// Call task manager to update task
	return s.taskManager.UpdateTask(taskID, &taskmanager.TaskInfo{
		Name:    req.Name,
		Status:  req.Status,
	})
}

// DeleteTask deletes a task
func (s *Service) DeleteTask(taskID string) error {
	if s.taskManager == nil {
		return fmt.Errorf("task manager not initialized")
	}

	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] DeleteTask: params=%s, duration=%dms", zap.String("task_id", taskID), zap.Int64("duration", duration))
	}()

	return s.taskManager.DeleteTask(taskID)
}

// GetLogs returns logs for a specific log type
func (s *Service) GetLogs(logType string, since string) ([]string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] GetLogs: params=%s, duration=%dms", zap.String("log_type", logType), zap.Int64("duration", duration))
	}()

	// Implementation to retrieve logs
	if logType == "error" {
		return []string{"Error log example", "Another error log"}, nil
	}
	return []string{"Normal log example"}, nil
}

// GetProcessList returns list of running processes
func (s *Service) GetProcessList() ([]*interfaces.ProcessInfo, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] GetProcessList: duration=%dms", zap.Int64("duration", duration))
	}()

	// Implementation to retrieve process list
	processes := []*interfaces.ProcessInfo{
		{ID: "proc1", Name: "example", CPU: 10.5, Memory: 1024, StartTime: "2023-01-01T00:00:00"},
	}
	return processes, nil
}

// UpdateConfig updates configuration values
func (s *Service) UpdateConfig(key string, value string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] UpdateConfig: params=%s=%s, duration=%dms", zap.String("key", key), zap.String("value", value), zap.Int64("duration", duration))
	}()
	
	// Implementation to update configuration
	return nil
}

// DownloadFile downloads a file
func (s *Service) DownloadFile(fileID string) (string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] DownloadFile: params=%s, duration=%dms", zap.String("file_id", fileID), zap.Int64("duration", duration))
	}()
	
	// Implementation to download file
	return "downloaded_file_path", nil
}

// GetHardwareStats returns hardware statistics
func (s *Service) GetHardwareStats() (*interfaces.HardwareStats, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] GetHardwareStats: duration=%dms", zap.Int64("duration", duration))
	}()

	// Implementation to get hardware stats
	return &interfaces.HardwareStats{
		CPUUsage:   50.2,
		MemoryUsed: 2048,
		DiskUsed:   1536,
	}, nil
}

// GetSystemInfo returns system information
func (s *Service) GetSystemInfo() (*interfaces.SystemInfo, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] GetSystemInfo: duration=%dms", zap.Int64("duration", duration))
	}()
	
	// Implementation to get system info
	return &interfaces.SystemInfo{
		OS:        "Linux",
		Arch:      "amd64",
		Hostname:  "server1",
		Users:     5,
		BootTime:  1686000000,
	}, nil
}

// GetTaskStats returns task statistics
func (s *Service) GetTaskStats() (*interfaces.TaskStats, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] GetTaskStats: duration=%dms", zap.Int64("duration", duration))
	}()

	if s.taskManager == nil {
		return nil, fmt.Errorf("task manager not initialized")
	}

	stats, err := s.taskManager.GetTaskStats()
	if err != nil {
		return nil, err
	}

	return &interfaces.TaskStats{
		TotalTasks:     stats.TotalTasks,
		ActiveTasks:    stats.ActiveTasks,
		CompletedTasks: stats.CompletedTasks,
		FailedTasks:    stats.FailedTasks,
		CPUUsage:       stats.CPUUsage,
		MemoryUsage:    stats.MemoryUsage,
		QueueLength:    stats.QueueLength,
		ActiveWorkers:  stats.ActiveWorkers,
		Timestamp:      stats.Timestamp,
	}, nil
}

// GetThreadPoolStats returns thread pool statistics
func (s *Service) GetThreadPoolStats() (*interfaces.ThreadPoolStats, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] GetThreadPoolStats: duration=%dms", zap.Int64("duration", duration))
	}()

	if s.executor == nil {
		return nil, fmt.Errorf("executor not initialized")
	}

	poolStats, err := s.executor.GetThreadPoolStats()
	if err != nil {
		return nil, err
	}
	
	return &interfaces.ThreadPoolStats{
		TotalPoolSize:  poolStats.TotalPoolSize,
		ActiveThreads:  poolStats.ActiveThreads,
		QueueSize:      poolStats.QueueSize,
		CompletedTasks: poolStats.CompletedTasks,
	}, nil
}

// UploadFile uploads a file
func (s *Service) UploadFile(file *multipart.FileHeader) (string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] UploadFile: duration=%dms", zap.Int64("duration", duration))
	}()
	
	// Implementation to upload file
	return "file_id_123", nil
}

// ExecuteCommand executes a command
func (s *Service) ExecuteCommand(cmd string, args []string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] ExecuteCommand: params=%s, duration=%dms", zap.String("cmd", cmd), zap.Int64("duration", duration))
	}()
	
	// Implementation to execute command
	return nil
}

// GetCommandHelp returns help for a command
func (s *Service) GetCommandHelp() string {
	s.logger.Info("[INFO] GetCommandHelp")
	return "Available commands: help, status, restart"
}

// GetStatus returns status information
func (s *Service) GetStatus() string {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] GetStatus: duration=%dms", zap.Int64("duration", duration))
	}()
	
	return s.monitor.Status()
}

// GetExecutorStatus returns executor status
func (s *Service) GetExecutorStatus() string {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] GetExecutorStatus: duration=%dms", zap.Int64("duration", duration))
	}()
	
	if s.executor != nil {
		return s.executor.Status()
	}
	return "executor not initialized"
}

// GetConfigList returns list of configuration values
func (s *Service) GetConfigList() []map[string]string {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.logger.Info("[INFO] GetConfigList: duration=%dms", zap.Int64("duration", duration))
	}()
	
	// Return configuration values as map
	configMap := []map[string]string{
		{"key": "database", "value": "pg"},
		{"key": "port", "value": "8080"},
	}
	return configMap
}

// Close cleans up resources
func (s *Service) Close() error {
	if s.shutdown != nil {
		return s.shutdown.Stop()
	}
	return nil
}