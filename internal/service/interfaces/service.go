package interfaces

import "mime/multipart"

type UpdateTaskRequest struct {
	Name   string
	Status string
}

type Service interface {
	GracefulShutdown()
	UploadFile(file *multipart.FileHeader) (string, error)
	ExecuteCommand(cmd string, args []string) error
	GetCommandHelp() string
	GetStatus() string
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
	GetExecutorStatus() string
	GetConfigList() []map[string]string
}

type TaskInterface interface {
	GetID() string
	GetName() string
	GetStatus() string
	SetStatus(status string)
	Execute() error
	GetStartTime() int64
	SetStartTime(startTime int64)
	GetDuration() int64
	SetDuration(duration int64)
	GetFinishedAt() int64
	SetFinishedAt(finishTime int64)
	GetCreatedAt() string
	GetCreator() string
	GetAssignedTo() string
	GetDescription() string
	GetResultPath() string
	GetProgress() int64
}

type ProcessInfo struct {
	ID       string
	Name     string
	CPU      float64
	Memory   uint64
	StartTime string
}

type HardwareStats struct {
	CPUUsage   float64
	MemoryUsed uint64
	DiskUsed   uint64
}

type SystemInfo struct {
	OS        string
	Arch      string
	Hostname  string
	Users     int
	BootTime  int64
}

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

type ThreadPoolStats struct {
	TotalTasks     int
	ActiveTasks    int
	CompletedTasks int
}

type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	StartedAt int64     `json:"started_at"`
	FinishedAt int64    `json:"finished_at"`
	Duration  int64     `json:"duration"`
	Logs      []string  `json:"logs"`
}

func (t *Task) GetID() string {
	return t.ID
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) GetStatus() string {
	return t.Status
}

func (t *Task) SetStatus(status string) {
	t.Status = status
}

func (t *Task) Execute() error {
	return nil
}

func (t *Task) GetStartTime() int64 {
	return t.StartedAt
}

func (t *Task) SetStartTime(startTime int64) {
	t.StartedAt = startTime
}

func (t *Task) GetDuration() int64 {
	return t.Duration
}

func (t *Task) SetDuration(duration int64) {
	t.Duration = duration
}

func (t *Task) GetFinishedAt() int64 {
	return t.FinishedAt
}

func (t *Task) SetFinishedAt(finishTime int64) {
	t.FinishedAt = finishTime
}

func (t *Task) GetCreatedAt() string {
	return t.CreatedAt
}