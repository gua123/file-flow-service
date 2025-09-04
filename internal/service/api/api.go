package api

import (
	"file-flow-service/internal/service/interfaces"
	"mime/multipart"
)

type Task struct {
	ID     string
	Dir    string
	Cmd    string
	Args   []string
}

type API struct {
	service interfaces.Service
}

func NewAPI(service interfaces.Service) *API {
	return &API{service: service}
}

// API implementation methods

func (a *API) UpdateTask(taskID string, req interfaces.UpdateTaskRequest) error {
	return a.service.UpdateTask(taskID, req)
}

func (a *API) DeleteTask(taskID string) error {
	return a.service.DeleteTask(taskID)
}

func (a *API) GetLogs(logType string, since string) ([]string, error) {
	return a.service.GetLogs(logType, since)
}

func (a *API) GetProcessList() ([]*interfaces.ProcessInfo, error) {
	return a.service.GetProcessList()
}

func (a *API) UpdateConfig(key string, value string) error {
	return a.service.UpdateConfig(key, value)
}

func (a *API) DownloadFile(fileID string) (string, error) {
	return a.service.DownloadFile(fileID)
}

func (a *API) GetHardwareStats() (*interfaces.HardwareStats, error) {
	return a.service.GetHardwareStats()
}

func (a *API) GetSystemInfo() (*interfaces.SystemInfo, error) {
	return a.service.GetSystemInfo()
}

func (a *API) GetTaskStats() (*interfaces.TaskStats, error) {
	return a.service.GetTaskStats()
}

func (a *API) GetThreadPoolStats() (*interfaces.ThreadPoolStats, error) {
	return a.service.GetThreadPoolStats()
}

func (a *API) UploadFile(file *multipart.FileHeader) (string, error) {
	return a.service.UploadFile(file)
}

func (a *API) ExecuteCommand(cmd string, args []string) error {
	return a.service.ExecuteCommand(cmd, args)
}

func (a *API) GetCommandHelp() string {
	return a.service.GetCommandHelp()
}

func (a *API) GetStatus() string {
	return a.service.GetStatus()
}