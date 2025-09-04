package service

import (
	"context"
	"file-flow-service/internal/shutdown"
	"file-flow-service/internal/service/executor"
	"file-flow-service/internal/taskmanager"
	"file-flow-service/internal/restart"
	"file-flow-service/internal/service/monitor"
	"file-flow-service/utils/logger"
	"file-flow-service/config"
	"file-flow-service/internal/threadpool"
	"file-flow-service/internal/service/interfaces"
	"mime/multipart"
	"strings"
)

type Service struct {
	*config.AppConfig
	TaskManager   taskmanager.TaskManager
	ShutdownManager *shutdown.ShutdownManager
	RestartManager *restart.RestartManager
	Monitor       *monitor.MonitorImpl
	Executor      *executor.BaseExecutor
	logger        logger.Logger
}

func NewService(config *config.AppConfig, logger logger.Logger) *Service {
	threadPool := threadpool.NewThreadPool()

	taskManager := taskmanager.NewTaskManager(config, threadPool, logger)
	shutdownManager := shutdown.NewShutdownManager(nil, logger, config)
	restartManager := restart.NewRestartManager(config, logger, nil)
	monitorImpl := monitor.NewMonitorImpl(logger, config)
	executor := executor.NewExecutor(config, logger)

	return &Service{
		AppConfig:           config,
		TaskManager:         taskManager,
		ShutdownManager:     shutdownManager,
		RestartManager:      restartManager,
		Monitor:             monitorImpl,
		Executor:            executor,
		logger:              logger,
	}
}

func (s *Service) Start() {
	s.TaskManager.Start()
	s.RestartManager.Start()
	s.Monitor.Start(context.Background())
}

func (s *Service) Stop() {
	s.ShutdownManager.GracefulShutdown()
	s.TaskManager.Stop()
	s.RestartManager.Stop()
	s.Monitor.Stop(context.Background())
}

func (s *Service) GetStatus() string {
	return "Service Status: Active"
}

func (s *Service) UploadFile(file *multipart.FileHeader) (string, error) {
	s.logger.Info("Uploading file: " + file.Filename)
	return file.Filename, nil
}

func (s *Service) ExecuteCommand(cmd string, args []string) error {
	s.logger.Info("Executing command: " + cmd + " " + strings.Join(args, " "))
	return nil
}

func (s *Service) GetCommandHelp() string {
	return "Available commands: help, status, restart"
}

func (s *Service) UpdateTask(taskID string, req interfaces.UpdateTaskRequest) error {
	return s.TaskManager.UpdateTask(taskID, req.Status)
}

func (s *Service) DeleteTask(taskID string) error {
	return nil
}

func (s *Service) DownloadFile(fileID string) (string, error) {
	return "", nil
}

func (s *Service) GetConfigList() []map[string]string {
	return []map[string]string{}
}

func (s *Service) GetExecutorStatus() string {
	return ""
}

func (s *Service) GetLogs(logType string, since string) ([]string, error) {
	return []string{}, nil
}

func (s *Service) GetProcessList() ([]*interfaces.ProcessInfo, error) {
	return []*interfaces.ProcessInfo{}, nil
}

func (s *Service) UpdateConfig(key string, value string) error {
	return nil
}

func (s *Service) GetHardwareStats() (*interfaces.HardwareStats, error) {
	return &interfaces.HardwareStats{}, nil
}

func (s *Service) GetSystemInfo() (*interfaces.SystemInfo, error) {
	return &interfaces.SystemInfo{}, nil
}

func (s *Service) GetTaskStats() (*interfaces.TaskStats, error) {
	return &interfaces.TaskStats{}, nil
}

func (s *Service) GetThreadPoolStats() (*interfaces.ThreadPoolStats, error) {
	return &interfaces.ThreadPoolStats{}, nil
}

func (s *Service) GracefulShutdown() {
	// Implement graceful shutdown logic here
}