package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"file-flow-service/internal/service/api"
	"file-flow-service/internal/test/testutils"
	"file-flow-service/web"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockWebAPI struct {
	taskID string
	err    error
}

func (m *mockWebAPI) CreateTask(task api.Task) (string, error) {
	return "test-task-1", m.err
}

func (m *mockWebAPI) GetTasks(page, pageSize int) ([]api.Task, error) {
	return []api.Task{}, nil
}

func (m *mockWebAPI) UpdateTask(taskID string, req api.UpdateTaskRequest) error {
	return nil
}

func (m *mockWebAPI) DeleteTask(taskID string) error {
	return nil
}

func (m *mockWebAPI) CancelTask(taskID string) error {
	return nil
}

func (m *mockWebAPI) GetExecutorStatus() string {
	return "active"
}

func (m *mockWebAPI) GetLogs(logType string, since string) ([]string, error) {
	return []string{}, nil
}

func (m *mockWebAPI) GetConfigList() []map[string]string {
	return []map[string]string{}
}

func (m *mockWebAPI) UpdateConfig(key string, value string) error {
	return nil
}

func (m *mockWebAPI) UploadFile(file *multipart.FileHeader) error {
	return nil
}

func (m *mockWebAPI) DownloadFile(fileID string) (string, error) {
	return "/tmp/test-result", nil
}

func (m *mockWebAPI) GetHardwareStats() (*api.HardwareStats, error) {
	return &api.HardwareStats{
		CPU:     []float64{50.0},
		Memory:  []int{1024 * 1024 * 1024},
		Disk:    []int{100 * 1024 * 1024 * 1024},
		Network: []int{1024, 1024},
	}, nil
}

func (m *mockWebAPI) GetSystemInfo() (*api.SystemInfo, error) {
	return &api.SystemInfo{
		Hostname:     "test-host",
		OS:           "linux",
		Platform:     "x86_64",
		Architecture: "amd64",
		Kernel:       "5.4.0",
		GoVersion:    "go1.24.5",
		StartTime:    1234567890,
	}, nil
}

func (m *mockWebAPI) GetProcessList() ([]*api.ProcessInfo, error) {
	return []*api.ProcessInfo{}, nil
}

func (m *mockWebAPI) GetTaskStats() (*api.TaskStats, error) {
	return &api.TaskStats{
		TotalTasks:     10,
		ActiveTasks:    2,
		CompletedTasks: 5,
		FailedTasks:    1,
		CPUUsage:       30.0,
		MemoryUsage:    1024 * 1024 * 1024,
		QueueLength:    3,
		ActiveWorkers:  2,
		Timestamp:      1234567890,
	}, nil
}

func (m *mockWebAPI) GetThreadPoolStats() (*api.ThreadPoolStats, error) {
	return &api.ThreadPoolStats{
		ActiveWorkers:  2,
		QueueLength:    3,
		TotalTasks:     100,
		FailedTasks:    5,
		CompletedTasks: 95,
	}, nil
}

func TestAPIIntegration(t *testing.T) {
	// 测试API接口的基本功能
	t.Run("TestTaskStruct", func(t *testing.T) {
		// 测试任务结构体
		task := api.Task{
			ID:          "test-task-1",
			Name:        "Test Task",
			Description: "A test task for integration testing",
			Status:      "pending",
			Progress:    0.0,
			Creator:     "test-user",
			AssignedTo:  "test-user",
			ResultPath:  "/tmp/test-result",
			Cmd:         "echo",
			Args:        []string{"hello", "world"},
			Dir:         "/tmp",
		}

		assert.Equal(t, "test-task-1", task.ID)
		assert.Equal(t, "Test Task", task.Name)
		assert.Equal(t, "A test task for integration testing", task.Description)
		assert.Equal(t, "pending", task.Status)
		assert.Equal(t, 0.0, task.Progress)
		assert.Equal(t, "test-user", task.Creator)
		assert.Equal(t, "test-user", task.AssignedTo)
		assert.Equal(t, "/tmp/test-result", task.ResultPath)
		assert.Equal(t, "echo", task.Cmd)
		assert.Equal(t, []string{"hello", "world"}, task.Args)
		assert.Equal(t, "/tmp", task.Dir)
	})

	t.Run("TestUpdateTaskRequest", func(t *testing.T) {
		// 测试更新任务请求结构体
		updateReq := api.UpdateTaskRequest{
			Name:   "Updated Task",
			Status: "running",
		}

		assert.Equal(t, "Updated Task", updateReq.Name)
		assert.Equal(t, "running", updateReq.Status)
	})

	t.Run("TestHardwareStats", func(t *testing.T) {
		// 测试硬件统计结构体
		stats, err := (&mockWebAPI{}).GetHardwareStats()
		assert.NoError(t, err)
		assert.Equal(t, float64(50.0), stats.CPU[0])
		assert.Equal(t, 1024*1024*1024, stats.Memory[0])
		assert.Equal(t, 100*1024*1024*1024, stats.Disk[0])
		assert.Equal(t, 1024, stats.Network[0])
		assert.Equal(t, 1024, stats.Network[1])
	})

	t.Run("TestSystemInfo", func(t *testing.T) {
		// 测试系统信息结构体
		info, err := (&mockWebAPI{}).GetSystemInfo()
		assert.NoError(t, err)
		assert.Equal(t, "test-host", info.Hostname)
		assert.Equal(t, "linux", info.OS)
		assert.Equal(t, "x86_64", info.Platform)
		assert.Equal(t, "amd64", info.Architecture)
		assert.Equal(t, "5.4.0", info.Kernel)
		assert.Equal(t, "go1.24.5", info.GoVersion)
		assert.Equal(t, int64(1234567890), info.StartTime)
	})

	t.Run("TestProcessInfo", func(t *testing.T) {
		// 测试进程信息结构体
		process := &api.ProcessInfo{
			PID:         1234,
			Name:        "test-process",
			CPUUsage:    25.0,
			Memory:      1024 * 1024,
			MemoryUsage: 10.0,
			Status:      "running",
			CmdLine:     "/usr/bin/test-process",
		}

		assert.Equal(t, int32(1234), process.PID)
		assert.Equal(t, "test-process", process.Name)
		assert.Equal(t, 25.0, process.CPUUsage)
		assert.Equal(t, 1024*1024, process.Memory)
		assert.Equal(t, 10.0, process.MemoryUsage)
		assert.Equal(t, "running", process.Status)
		assert.Equal(t, "/usr/bin/test-process", process.CmdLine)
	})

	t.Run("TestTaskStats", func(t *testing.T) {
		// 测试任务统计结构体
		stats, err := (&mockWebAPI{}).GetTaskStats()
		assert.NoError(t, err)
		assert.Equal(t, 10, stats.TotalTasks)
		assert.Equal(t, 2, stats.ActiveTasks)
		assert.Equal(t, 5, stats.CompletedTasks)
		assert.Equal(t, 1, stats.FailedTasks)
		assert.Equal(t, 30.0, stats.CPUUsage)
		assert.Equal(t, 1024*1024*1024, stats.MemoryUsage)
		assert.Equal(t, 3, stats.QueueLength)
		assert.Equal(t, 2, stats.ActiveWorkers)
		assert.Equal(t, int64(1234567890), stats.Timestamp)
	})

	t.Run("TestThreadPoolStats", func(t *testing.T) {
		// 测试线程池统计结构体
		stats, err := (&mockWebAPI{}).GetThreadPoolStats()
		assert.NoError(t, err)
		assert.Equal(t, 2, stats.ActiveWorkers)
		assert.Equal(t, 3, stats.QueueLength)
		assert.Equal(t, 100, stats.TotalTasks)
		assert.Equal(t, 5, stats.FailedTasks)
		assert.Equal(t, 95, stats.CompletedTasks)
	})
}

func TestHTTPAPI(t *testing.T) {
	// Setup test server
	mockService := &mockWebAPI{taskID: "test-task-id", err: nil}
	router := web.NewRouter(mockService)
	server := httptest.NewServer(router)
	defer server.Close()

	// Test CreateTask
	task := api.Task{
		Name:        "Test Task",
		Description: "Test description",
		Status:      "pending",
		Creator:     "test-user",
		Dir:         "/tmp",
		Cmd:         "echo",
		Args:        []string{"hello"},
	}
	jsonData, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestAPIFunctionality(t *testing.T) {
	// 测试API功能的基本验证
	t.Run("TestConfigCreation", func(t *testing.T) {
		cfg := testutils.CreateTestConfig()
		assert.NotNil(t, cfg)
		assert.Equal(t, 8080, cfg.App.Port)
		assert.Equal(t, "test-service", cfg.App.Name)
		assert.Equal(t, "test", cfg.App.Env)
	})

	t.Run("TestTaskCreation", func(t *testing.T) {
		task := testutils.CreateTestTask()
		assert.NotEmpty(t, task.ID)
		assert.Equal(t, "Test Task", task.Name)
		assert.Equal(t, "A test task for testing purposes", task.Description)
		assert.Equal(t, "pending", task.Status)
		assert.Equal(t, 0.0, task.Progress)
		assert.Equal(t, "test-user", task.Creator)
		assert.Equal(t, "test-user", task.AssignedTo)
		assert.Equal(t, "/tmp/test-result", task.ResultPath)
		assert.Equal(t, "echo", task.Cmd)
		assert.Equal(t, []string{"hello"}, task.Args)
		assert.Equal(t, "/tmp", task.Dir)
	})

	t.Run("TestTaskCreationWithError", func(t *testing.T) {
		mockService := &mockWebAPI{err: errors.New("mock error")}
		router := web.NewRouter(mockService)
		server := httptest.NewServer(router)
		defer server.Close()

		task := api.Task{
			Name:        "Test Task",
			Description: "Test description",
			Status:      "pending",
			Creator:     "test-user",
			Dir:         "/tmp",
			Cmd:         "echo",
			Args:        []string{"hello"},
		}
		jsonData, _ := json.Marshal(task)
		req, _ := http.NewRequest("POST", server.URL+"/tasks", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
