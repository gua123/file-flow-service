package interfaces

import "mime/multipart"

type UpdateTaskRequest struct {
	Name   string
	Status string
}

type Service interface {
	UpdateTask(taskID string, req UpdateTaskRequest) error
	DeleteTask(taskID string) error
	GetLogs(logType string, since string) ([]string, error)
	GetProcessList() ([]*ProcessInfo, error)
	UpdateConfig(key string, value string) error
	DownloadFile(fileID string) (string, error)
	GetHardwareStats() (*HardwareStats, error)
	GetSystemInfo() (*SystemInfo, error)
	GetTaskStats() (*TaskStats, error)
	GetThreadPoolStats() (*ThreadPoolStats, error)
	UploadFile(file *multipart.FileHeader) (string, error)
	ExecuteCommand(cmd string, args []string) error
	GetCommandHelp() string
	GetStatus() string
	GetExecutorStatus() string
	GetConfigList() []map[string]string
}

// ProcessInfo is used by GetProcessList
type ProcessInfo struct {
	ID       string
	Name     string
	CPU      float64
	Memory   uint64
	StartTime string
}

// HardwareStats contains system hardware metrics
type HardwareStats struct {
	CPUUsage   float64
	MemoryUsed uint64
	DiskUsed   uint64
}

// SystemInfo contains system information
type SystemInfo struct {
	OS        string
	Arch      string
	Hostname  string
	Users     int
	BootTime  int64
}

// TaskStats contains task execution metrics
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

// ThreadPoolStats contains thread pool metrics
type ThreadPoolStats struct {
	TotalPoolSize  int
	ActiveThreads  int
	QueueSize      int
	CompletedTasks int
}