// Package api API layer for service
package api

import (
	"errors"
	"time"
	"mime/multipart"
	"file-flow-service/internal/service"
)

// HardwareStats provides hardware statistics.
type HardwareStats struct {
	CPU       []float64
	Memory    []int
	Disk      []int
	Network   []int
}

// SystemInfo provides system information.
type SystemInfo struct {
	Hostname     string
	OS           string
	Architecture string
	Kernel       string
}

// TaskStats provides task statistics.
type TaskStats struct {
	Active      int
	Completed   int
	Failed      int
	Running     int
}

// ThreadPoolStats provides thread pool statistics.
type ThreadPoolStats struct {
	ActiveWorkers  int
	QueueLength    int
	TotalTasks     int
	FailedTasks    int
	CompletedTasks int
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

type API struct {
	svc Service
}

func NewAPI() *API {
	return &API{svc: service.NewService()}
}

// ProcessInfo describes the details of a running process.
type ProcessInfo struct {
	PID         int
	Name        string
	CPUUsage    float64
	Memory      int
	MemoryUsage float64
	Status      string
	CmdLine     string
}

// Task represents a task in the API.
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

// UpdateTaskRequest represents the request payload for updating a task.
type UpdateTaskRequest struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// CreateTaskRequest represents the request payload for creating a task.
type CreateTaskRequest struct {
	Name string `json:"name"`
}

// GetPoolStatus gets thread pool status
func (a *API) GetPoolStatus() (*service.ThreadPoolStats, error) {
	return a.svc.GetThreadPoolStats()
}

// UpdateTask 更新任务
func (a *API) UpdateTask(taskID string, req UpdateTaskRequest) error {
	// 验证任务ID
	if taskID == "" {
		return errors.New("任务ID不能为空")
	}

	// 验证任务名称
	if req.Name == "" {
		return errors.New("任务名称不能为空")
	}

	// 验证任务状态
	if req.Status == "" {
		return errors.New("任务状态不能为空")
	}

	// 调用服务方法
	return a.svc.UpdateTask(taskID, req)
}

// DeleteTask 删除任务
func (a *API) DeleteTask(taskID string) error {
	// 验证任务ID
	if taskID == "" {
		return errors.New("任务ID不能为空")
	}

	// 调用服务方法
	return a.svc.DeleteTask(taskID)
}

// GetLogs 获取日志
func (a *API) GetLogs(logType string, since string) ([]string, error) {
	// 验证日志类型
	if logType == "" {
		return nil, errors.New("日志类型不能为空")
	}

	// 验证时间范围
	if since == "" {
		return nil, errors.New("时间范围不能为空")
	}

	// 调用服务方法
	return a.svc.GetLogs(logType, since)
}

// GetProcessList 获取进程列表
func (a *API) GetProcessList() ([]*ProcessInfo, error) {
	// 调用服务方法
	return a.svc.GetProcessList()
}

// UpdateConfig 更新配置
func (a *API) UpdateConfig(key string, value string) error {
	// 验证配置键
	if key == "" {
		return errors.New("配置键不能为空")
	}

	// 验证配置值
	if value == "" {
		return errors.New("配置值不能为空")
	}

	// 调用服务方法
	return a.svc.UpdateConfig(key, value)
}

// DownloadFile 下载文件
func (a *API) DownloadFile(fileID string) (string, error) {
	// 验证文件ID
	if fileID == "" {
		return "", errors.New("文件ID不能为空")
	}

	// 调用服务方法
	return a.svc.DownloadFile(fileID)
}