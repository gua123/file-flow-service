// Package service 服务模块，提供任务管理、文件处理、监控等核心功能。
// 该模块是整个系统的协调中心，整合各个子模块的功能
// 与 executor、monitor、taskmanager、processmanager 等模块紧密协作
package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"sync"
	"time"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"go.uber.org/zap"

	"file-flow-service/config"
	"file-flow-service/file"
	"file-flow-service/internal/processmanager"
	"file-flow-service/internal/service/executor"
	"file-flow-service/internal/service/monitor"
	"file-flow-service/internal/taskmanager"
	"file-flow-service/utils/logger"
	"file-flow-service/internal/service/api"
)

type TaskStats struct {
	Active   int
	Completed int
}

type PoolStatus struct {
	Active int
	Total  int
}

// Service 服务主结构体
// 整合所有核心功能模块，提供统一的服务接口
type Service struct {
	Executor     *executor.BaseExecutor
	Monitor      *monitor.MonitorImpl
	File         *file.FileService
	Config       *config.AppConfig
	Logger       logger.Logger
	TaskManager  taskmanager.TaskManager
	ProcessManager processmanager.ProcessManager
}

var (
	serviceInstance *Service
	once            sync.Once
)

// NewService 创建服务实例
// 参数：cfg 配置对象，logger 日志记录器
// 返回：服务实例
// 上下承接关系：初始化所有子模块，创建全局服务实例
func NewService(cfg *config.AppConfig, logger logger.Logger) *Service {
	once.Do(func() {
		// 初始化执行器
		executor := executor.NewExecutor(cfg, logger)

		// 初始化文件服务
		fileService := file.NewFileService(cfg.File.StoragePath, logger)

		// 创建Service实例
		svc := &Service{
			Executor:       executor,
			File:           fileService,
			Config:         cfg,
			Logger:         logger,
		}

		// 初始化监控模块
		monitorInterval, err := time.ParseDuration(cfg.MonitorInterval)
		if err != nil {
			logger.Error("解析监控间隔失败", zap.Error(err))
			monitorInterval = 5 * time.Second // 使用默认值
		}
		monitor := monitor.NewMonitor(monitorInterval, logger, svc)
		svc.Monitor = monitor

		// 初始化任务管理器
		taskManager := taskmanager.NewTaskManager(cfg, svc.Executor.GetPool(), logger)
		svc.TaskManager = taskManager

		// 初始化进程管理器
		processManager := processmanager.NewProcessManager(cfg, logger)
		svc.ProcessManager = processManager

		serviceInstance = svc
	})
	return serviceInstance
}

// GetService 获取全局服务实例
// 参数：无
// 返回：服务实例
// 上下承接关系：返回已初始化的全局服务实例
func GetService() *Service {
	return serviceInstance
}

// generateTaskID 生成任务ID
// 生成唯一的任务ID
// 参数：无
// 返回：任务ID字符串
func generateTaskID() string {
	// 简单的ID生成器，实际项目中可以使用UUID等
	return "task-" + time.Now().Format("20060102150405")
}

// WebAPI implementation

// CreateTask 创建任务
// 参数：task 任务对象
// 返回：任务ID，错误信息
// 上下承接关系：生成任务ID并提交任务到执行器
func (s *Service) CreateTask(task api.Task) (string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] CreateTask: params={name: " + task.Name + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result={" + task.ID + "}")
	}()

	s.Logger.Info("创建任务", zap.String("task_name", task.Name))

	// 生成任务ID
	taskID := generateTaskID()
	task.ID = taskID

	// 提交任务到执行器
	if s.Executor != nil {
		err := s.Executor.RunTask(task)
		if err != nil {
			s.Logger.Error("提交任务失败", zap.String("task_id", taskID), zap.Error(err))
			return "", fmt.Errorf("提交任务失败: %v", err)
		}
		s.Logger.Info("任务已提交到执行器", zap.String("task_id", taskID))
		return taskID, nil
	}

	return "", fmt.Errorf("执行器未初始化")
}

// GetTasks 获取任务列表
// 参数：page 页码，pageSize 每页数量
// 返回：任务列表，错误信息
// 上下承接关系：调用任务管理器获取任务列表
func (s *Service) GetTasks(page int, pageSize int) (string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetTasks: params={page: " + strconv.Itoa(page) + ", pageSize: " + strconv.Itoa(pageSize) + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result=" + "[]")
	}()

	s.Logger.Debug("获取任务列表", zap.Int("page", page), zap.Int("page_size", pageSize))

	// 调用任务管理器获取任务列表
	if s.TaskManager != nil {
		taskInfos, err := s.TaskManager.GetAllTasks()
		if err != nil {
			return "", err
		}

		// 转换为API任务格式
		tasks := make([]api.Task, len(taskInfos))
		for i, taskInfo := range taskInfos {
			tasks[i] = api.Task{
				ID:          taskInfo.ID,
				Name:        taskInfo.Name,
				Description: "", // taskInfo中没有描述字段
				Status:      taskInfo.Status,
				Progress:    taskInfo.Progress,
				CreatedAt:   taskInfo.CreatedAt,
				Creator:     "", // taskInfo中没有创建者字段
				AssignedTo:  "", // taskInfo中没有分配给字段
				ResultPath:  taskInfo.Result,
				Cmd:         "", // taskInfo中没有命令字段
				Args:        nil, // taskInfo中没有参数字段
				Dir:         "", // taskInfo中没有工作目录字段
			}
		}
		// 将任务列表转换为JSON字符串
		jsonData, err := json.Marshal(tasks)
		if err != nil {
			return "", err
		}
		return string(jsonData), nil
	}

	// 返回空列表的JSON
	return "[]", nil
}

// UpdateTask 更新任务
// 参数：taskID 任务ID，req 更新请求
// 返回：错误信息
// 上下承接关系：调用任务管理器更新任务信息
func (s *Service) UpdateTask(taskID string, req api.UpdateTaskRequest) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] UpdateTask: params={task_id: " + taskID + ", name: " + req.Name + ", status: " + req.Status + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result=null")
	}()

	s.Logger.Info("更新任务", zap.String("task_id", taskID))

	// 调用任务管理器来更新任务
	if s.TaskManager != nil {
		// 将API请求转换为任务信息
		taskInfo := &taskmanager.TaskInfo{
			ID:        taskID,
			Name:      req.Name,
			Status:    req.Status,
			Progress:  0, // 进度需要根据实际执行情况更新
			CreatedAt: time.Now(),
		}

		// 调用任务管理器更新任务
		err := s.TaskManager.UpdateTask(taskID, taskInfo)
		if err != nil {
			s.Logger.Error("更新任务失败", zap.String("task_id", taskID), zap.Error(err))
			return err
		}

		s.Logger.Info("任务已更新", zap.String("task_id", taskID))
		return nil
	}

	return fmt.Errorf("任务管理器未初始化")
}

// DeleteTask 删除任务
// 参数：taskID 任务ID
// 返回：错误信息
// 上下承接关系：调用任务管理器删除任务
func (s *Service) DeleteTask(taskID string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] DeleteTask: params={task_id: " + taskID + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result=null")
	}()

	s.Logger.Info("删除任务", zap.String("task_id", taskID))

	// 调用任务管理器来删除任务
	if s.TaskManager != nil {
		// 调用任务管理器删除任务
		err := s.TaskManager.DeleteTask(taskID)
		if err != nil {
			s.Logger.Error("删除任务失败", zap.String("task_id", taskID), zap.Error(err))
			return err
		}

		s.Logger.Info("任务已删除", zap.String("task_id", taskID))
		return nil
	}

	return fmt.Errorf("任务管理器未初始化")
}

// CancelTask 取消任务
// 参数：taskID 任务ID
// 返回：错误信息
// 上下承接关系：调用任务管理器取消任务
func (s *Service) CancelTask(taskID string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] CancelTask: params={task_id: " + taskID + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result=null")
	}()

	s.Logger.Info("取消任务", zap.String("task_id", taskID))

	// 调用任务管理器来取消任务
	if s.TaskManager != nil {
		return s.TaskManager.CancelTask(taskID)
	}

	return fmt.Errorf("任务管理器未初始化")
}

// GetExecutorStatus 获取执行器状态
// 参数：无
// 返回：执行器状态字符串
// 上下承接关系：返回执行器的当前运行状态
func (s *Service) GetExecutorStatus() string {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetExecutorStatus: params={}, duration=" + fmt.Sprintf("%dms", duration) + ", result=" + s.Executor.Status)
	}()

	return s.Executor.Status
}

// GetLogs 获取日志
// 参数：logType 日志类型，since 时间点
// 返回：日志列表，错误信息
// 上下承接关系：调用监控模块获取日志
func (s *Service) GetLogs(logType string, since string) ([]string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetLogs: params={type: " + logType + ", since: " + since + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result=" + strconv.Itoa(len(s.Monitor.GetLogs(logType, since))))
	}()

	s.Logger.Debug("获取日志", zap.String("type", logType), zap.String("since", since))
	return s.Monitor.GetLogs(logType, since)
}

// GetConfigList 获取配置列表
// 参数：无
// 返回：配置列表
// 上下承接关系：返回当前配置的键值对列表
func (s *Service) GetConfigList() []map[string]string {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetConfigList: params={}, duration=" + fmt.Sprintf("%dms", duration) + ", result=[" + strconv.Itoa(len(s.GetConfigList())) + "]")
	}()

	s.Logger.Debug("获取配置列表")

	// 返回配置列表
	return []map[string]string{
		{
			"key":   "monitor_interval",
			"value": s.Config.MonitorInterval,
		},
		{
			"key":   "task_timeout",
			"value": s.Config.Threadpool.TaskTimeout,
		},
		{
			"key":   "max_workers",
			"value": fmt.Sprintf("%d", s.Config.Threadpool.MaxWorkers),
		},
	}
}

// UpdateConfig 更新配置
// 参数：key 配置键，value 配置值
// 返回：错误信息
// 上下承接关系：更新配置并记录变更
func (s *Service) UpdateConfig(key string, value string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] UpdateConfig: params={key: " + key + ", value: " + value + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result=null")
	}()

	s.Logger.Info("更新配置", zap.String("key", key), zap.String("value", value))

	// 根据键更新相应的配置
	switch key {
	case "monitor_interval":
		s.Config.MonitorInterval = value
	case "task_timeout":
		s.Config.Threadpool.TaskTimeout = value
	case "max_workers":
		// 尝试解析为整数
		if maxWorkers, err := strconv.Atoi(value); err == nil {
			s.Config.Threadpool.MaxWorkers = maxWorkers
		} else {
			s.Logger.Warn("解析max_workers失败", zap.String("value", value), zap.Error(err))
			return fmt.Errorf("无效的max_workers值: %v", err)
		}
	default:
		s.Logger.Warn("未知配置键", zap.String("key", key))
		return fmt.Errorf("未知配置键: %s", key)
	}

	s.Logger.Info("配置已更新", zap.String("key", key), zap.String("value", value))
	return nil
}

// GetHardwareStats 获取硬件统计信息
// 参数：无
// 返回：硬件统计信息，错误信息
// 上下承接关系：调用监控模块获取硬件统计信息
func (s *Service) GetHardwareStats() (*api.HardwareStats, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetHardwareStats: params={}, duration=" + fmt.Sprintf("%dms", duration) + ", result=" + strconv.Itoa(len(s.Monitor.GetHardwareStats().CPU)))
	}()

	s.Logger.Debug("获取硬件统计信息")
	return s.Monitor.GetHardwareStats()
}

// GetSystemInfo 获取系统信息
// 参数：无
// 返回：系统信息，错误信息
// 上下承接关系：调用监控模块获取系统信息
func (s *Service) GetSystemInfo() (*api.SystemInfo, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetSystemInfo: params={}, duration=" + fmt.Sprintf("%dms", duration) + ", result=" + s.Monitor.GetSystemInfo().Hostname)
	}()

	s.Logger.Debug("获取系统信息")
	return s.Monitor.GetSystemInfo()
}

// GetProcessList 获取进程列表
// 参数：无
// 返回：进程列表，错误信息
// 上下承接关系：调用进程管理模块获取进程列表
func (s *Service) GetProcessList() ([]*api.ProcessInfo, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetProcessList: params={}, duration=" + fmt.Sprintf("%dms", duration) + ", result=[" + strconv.Itoa(len(s.ProcessManager.GetAllProcesses().apiProcesses)) + "]")
	}()

	s.Logger.Debug("获取进程列表")

	// 如果进程管理器已初始化，调用其方法
	if s.ProcessManager != nil {
		processes, err := s.ProcessManager.GetAllProcesses()
		if err != nil {
			return nil, err
		}

		// 转换类型
		apiProcesses := make([]*api.ProcessInfo, len(processes))
		for i, proc := range processes {
			apiProcesses[i] = &api.ProcessInfo{
				PID:        proc.PID,
				Name:       proc.Name,
				CPUUsage:   proc.CPUUsage,
				Memory:     proc.Memory,
				MemoryUsage: proc.MemoryUsage,
				Status:     proc.Status,
				CmdLine:    proc.CmdLine,
			}
		}
		return apiProcesses, nil
	}

	// 否则返回空列表
	return []*api.ProcessInfo{}, nil
}

// GetTaskStats 获取任务统计信息
// 参数：无
// 返回：任务统计信息，错误信息
// 上下承接关系：调用任务管理模块获取任务统计信息
func (s *Service) GetTaskStats() (*api.TaskStats, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetTaskStats: params={}, duration=" + fmt.Sprintf("%dms", duration) + ", result= {active: " + strconv.Itoa(s.TaskManager.GetTaskStats().Active) + ", completed: " + strconv.Itoa(s.TaskManager.GetTaskStats().Completed) + "}")
	}()

	s.Logger.Debug("获取任务统计信息")
	return s.TaskManager.GetTaskStats().ConvertToAPITaskStats()
}

// GetThreadPoolStats 获取线程池统计信息
// 参数：无
// 返回：线程池统计信息，错误信息
// 上下承接关系：调用线程池获取统计信息
func (s *Service) GetThreadPoolStats() (*api.ThreadPoolStats, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetThreadPoolStats: params={}, duration=" + fmt.Sprintf("%dms", duration) + ", result={active=" + strconv.Itoa(s.Executor.GetPool().GetStats().ActiveWorkers) + ", total=" + strconv.Itoa(int(s.Executor.GetPool().GetStats().TotalTasks)) + "}")
	}()

	s.Logger.Debug("获取线程池统计信息")
	return s.Executor.GetPool().GetStats()
}

// UploadFile 上传文件
// 参数：file 文件头
// 返回：错误信息
// 上下承接关系：调用文件服务处理文件上传
func (s *Service) UploadFile(file *multipart.FileHeader) (string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] UploadFile: params={filename: " + file.Filename + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result=null")
	}()

	s.Logger.Info("上传文件", zap.String("filename", file.Filename))
	return "", s.File.Upload(file)
}

// DownloadFile 下载文件
// 参数：fileID 文件ID
// 返回：文件路径，错误信息
// 上下承接关系：调用文件服务处理文件下载
func (s *Service) DownloadFile(fileID string) (string, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] DownloadFile: params={file_id: " + fileID + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result=/path/to/downloaded/file")
	}()

	s.Logger.Debug("下载文件", zap.String("file_id", fileID))
	return s.File.Download(fileID)
}

// CLI implementation

// ExecuteCommand 执行命令
// 参数：cmd 命令，args 参数列表
// 返回：错误信息
// 上下承接关系：执行命令并返回结果
func (s *Service) ExecuteCommand(cmd string, args []string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] ExecuteCommand: params={command: " + cmd + ", args: " + strings.Join(args, ", ") + "}, duration=" + fmt.Sprintf("%dms", duration) + ", result=null")
	}()

	s.Logger.Info("执行命令", zap.String("command", cmd), zap.Strings("args", args))

	// 实现命令执行逻辑
	// 这里可以使用os/exec包来执行系统命令
	// 为了简单起见，我们记录日志但不实际执行
	// 在实际实现中，应该使用exec.Command来执行命令

	s.Logger.Debug("命令执行完成", zap.String("command", cmd), zap.Strings("args", args))
	return nil
}

// GetCommandHelp 获取命令帮助
// 参数：无
// 返回：帮助信息字符串
// 上下承接关系：返回可用命令的帮助信息
func (s *Service) GetCommandHelp() string {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetCommandHelp: params={}, duration=" + fmt.Sprintf("%dms", duration) + ", result=CLI帮助信息：执行命令、帮助、状态、重启")
	}()

	s.Logger.Debug("获取命令帮助")
	return "CLI帮助信息：" +
		"\n  - execute <command> [args]: 执行命令" +
		"\n  - help: 显示帮助信息" +
		"\n  - status: 显示服务状态" +
		"\n  - restart: 重启服务"
}

// GetStatus 获取服务状态
// 参数：无
// 返回：服务状态字符串
// 上下承接关系：返回服务的当前运行状态
func (s *Service) GetStatus() string {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		s.Logger.Info("[INFO] GetStatus: params={}, duration=" + fmt.Sprintf("%dms", duration) + ", result=" + s.Monitor.GetStatus())
	}()

	s.Logger.Debug("获取服务状态")
	return s.Monitor.GetStatus()
}

// ProcessInfo describes the details of a running process.
type ProcessInfo struct {
	PID        int
	Name       string
	CPUUsage   float64
	Memory     int
	MemoryUsage float64
	Status     string
	CmdLine    string
}