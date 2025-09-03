package interfaces

import (
	"time"
	"mime/multipart"
)

var GlobalService Service

func GetService() Service {
	return GlobalService
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

type UpdateTaskRequest struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type CreateTaskRequest struct {
	Name string `json:"name"`
}

type ProcessInfo struct {
	PID        int
	Name       string
	CPUUsage   float64
	Memory     int
	MemoryUsage float64
	Status     string
	CmdLine    string
}

type Task struct {
	ID          string
	Name        string
	Description string
	Status      string
	Progress    int
	CreatedAt   time.Time
	Creator     string
	AssignedTo  string
	ResultPath  string
	Cmd         string
	Args        []string
	Dir         string
}

type HardwareStats struct {
	CPU       []float64
	Memory    []int
	Disk      []int
	Network   []int
	CPUUsage  float64 // Required by monitor
	MemoryUsed int    // Required by monitor
	MemoryTotal int   // Required by monitor
	MemoryFree int    // Required by monitor
	MemoryUsage float64 // Required by monitor
}

type SystemInfo struct {
	Hostname     string
	OS           string
	Architecture string
	Kernel       string
}

type TaskStats struct {
	Active      int
	Completed   int
	Failed      int
	Running     int
	TotalTasks  int   // Required by taskmanager
	ActiveTasks int   // Required by taskmanager
	CompletedTasks int  // Required by taskmanager
	FailedTasks int    // Required by taskmanager
	QueueLength int    // Required by monitor
	ActiveWorkers int  // Required by monitor
}

type ThreadPoolStats struct {
	ActiveWorkers  int
	QueueLength    int
	TotalTasks     int
	FailedTasks    int
	CompletedTasks int
}