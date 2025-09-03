package testutils

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"time"

	"file-flow-service/config"
	"file-flow-service/internal/service/api"
)

// CreateTestConfig 创建测试配置
func CreateTestConfig() *config.AppConfig {
	return &config.AppConfig{
		App: config.App{
			Port: 8080,
			Name: "test-service",
			Env:  "test",
		},
		LoggerConf: config.LoggerConf{
			BasePath: "./test-log",
			Levels: map[string]bool{
				"debug": true,
				"info":  true,
				"warn":  true,
				"error": true,
			},
		},
		Threadpool: config.Threadpool{
			MaxWorkers:  5,
			MaxQueue:    10,
			TaskTimeout: "5m",
		},
		File: config.File{
			StoragePath:   "./test-files",
			MaxUploadSize: 10485760,
		},
		Monitoring: config.Monitoring{
			HealthCheck: config.HealthCheck{
				Interval: "1m",
			},
		},
		Permissions: config.Permissions{
			DefaultFileMode: 0644,
			DefaultDirMode:  0755,
		},
	}
}

// CreateTestTask 创建测试任务
func CreateTestTask() api.Task {
	return api.Task{
		ID:          "test-task-1",
		Name:        "Test Task",
		Description: "A test task for testing purposes",
		Status:      "pending",
		Progress:    0.0,
		CreatedAt:   time.Now(),
		Creator:     "test-user",
		AssignedTo:  "test-user",
		ResultPath:  "/tmp/test-result",
		Cmd:         "echo",
		Args:        []string{"hello"},
		Dir:         "/tmp",
	}
}

// CreateTestMultipartForm 创建测试multipart表单
func CreateTestMultipartForm() (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// 添加文件字段
	fileWriter, err := writer.CreateFormFile("file", "test.txt")
	if err != nil {
		return nil, "", err
	}
	
	// 写入文件内容
	fileWriter.Write([]byte("test content"))
	
	// 关闭writer
	writer.Close()
	
	return &buf, writer.FormDataContentType(), nil
}

// CreateJSONRequest 创建JSON请求
func CreateJSONRequest(method, url string, data interface{}) (*http.Request, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	req := httptest.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// CreateTestResponse 创建测试响应
func CreateTestResponse(code int, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	}
}
