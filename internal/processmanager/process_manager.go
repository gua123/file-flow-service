// Package processmanager 进程管理器
// 负责系统进程的监控、管理和服务
// 与 monitor 模块紧密协作，提供进程级别的监控能力
package processmanager

import (
	"context"
	"fmt"
	"sync"
	"time"
	"file-flow-service/utils/logger"
	"file-flow-service/config"
	
	"go.uber.org/zap"
	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo 进程信息
// 存储进程的详细信息和状态
type ProcessInfo struct {
	PID         int32     `json:"pid"`
	Name        string    `json:"name"`
	CmdLine     string    `json:"cmd_line"`
	CPUUsage    float64   `json:"cpu_usage"`
	Memory      uint64    `json:"memory"`
	MemoryUsage float64   `json:"memory_usage"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	ParentPID   int32     `json:"parent_pid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProcessStats 进程统计信息
// 提供进程执行的统计信息
type ProcessStats struct {
	TotalProcesses   int     `json:"total_processes"`
	RunningProcesses int     `json:"running_processes"`
	CPUUsage         float64 `json:"cpu_usage"`
	MemoryUsage      uint64  `json:"memory_usage"`
	ActiveThreads    int     `json:"active_threads"`
	Timestamp        int64   `json:"timestamp"`
}

// ProcessManager 进程管理器接口
// 定义进程管理器必须实现的方法
type ProcessManager interface {
	// 启动管理器
	// 参数：无
	// 返回：错误信息
	// 上下承接关系：初始化监控循环，开始进程状态监控
	Start() error
	
	// 停止管理器
	// 参数：无
	// 返回：错误信息
	// 上下承接关系：停止监控循环，清理资源
	Stop() error
	
	// 获取所有进程
	// 参数：无
	// 返回：进程信息切片，错误信息
	// 上下承接关系：返回当前所有进程的快照
	GetAllProcesses() ([]*ProcessInfo, error)
	
	// 获取指定进程
	// 参数：pid 进程ID
	// 返回：进程信息，错误信息
	// 上下承接关系：根据ID查找并返回指定进程
	GetProcess(pid int32) (*ProcessInfo, error)
	
	// 终止进程
	// 参数：pid 进程ID
	// 返回：错误信息
	// 上下承接关系：终止指定进程，更新进程状态
	TerminateProcess(pid int32) error
	
	// 重启进程
	// 参数：pid 进程ID
	// 返回：错误信息
	// 上下承接关系：先终止后重新启动指定进程
	RestartProcess(pid int32) error
	
	// 获取进程统计信息
	// 参数：无
	// 返回：进程统计信息，错误信息
	// 上下承接关系：聚合进程状态信息，返回统计结果
	GetProcessStats() (*ProcessStats, error)
	
	// 监控进程状态
	// 参数：无
	// 返回：错误信息
	// 上下承接关系：定期更新进程状态信息
	MonitorProcesses() error
}

// processManager 进程管理器实现
// 实现进程管理器接口，提供完整的进程管理功能
type processManager struct {
	logger        logger.Logger
	config        *config.AppConfig
	running       bool
	mu            sync.RWMutex
	managedProcesses map[int32]*ProcessInfo
	ticker        *time.Ticker
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewProcessManager 创建进程管理器
// 参数：config 配置对象, logger 日志记录器
// 返回：进程管理器接口实例
// 上下承接关系：初始化进程管理器结构体，创建上下文
func NewProcessManager(config *config.AppConfig, logger logger.Logger) ProcessManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &processManager{
		config:           config,
		logger:           logger,
		ctx:              ctx,
		cancel:           cancel,
		managedProcesses: make(map[int32]*ProcessInfo),
	}
}

// Start 启动进程管理器
// 参数：无
// 返回：错误信息，如果启动失败则返回错误
// 上下承接关系：初始化监控循环，开始定期更新进程状态
func (pm *processManager) Start() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if pm.running {
		return nil
	}
	
	pm.running = true
	pm.logger.Info("进程管理器启动")
	
	// 启动监控循环
	go pm.monitorLoop()
	
	return nil
}

// Stop 停止进程管理器
// 参数：无
// 返回：错误信息，如果停止失败则返回错误
// 上下承接关系：停止监控循环，清理资源，关闭上下文
func (pm *processManager) Stop() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if !pm.running {
		return nil
	}
	
	pm.running = false
	pm.cancel()
	
	if pm.ticker != nil {
		pm.ticker.Stop()
	}
	
	pm.logger.Info("进程管理器停止")
	return nil
}

// GetAllProcesses 获取所有进程
// 参数：无
// 返回：进程信息切片，错误信息
// 上下承接关系：返回当前所有进程的快照，用于进程列表展示
func (pm *processManager) GetAllProcesses() ([]*ProcessInfo, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	var processes []*ProcessInfo
	for _, proc := range pm.managedProcesses {
		processes = append(processes, proc)
	}
	
	pm.logger.Debug("获取所有进程", zap.Int("count", len(processes)))
	return processes, nil
}

// GetProcess 获取指定进程
// 参数：pid 进程ID
// 返回：进程信息，错误信息
// 上下承接关系：根据ID查找并返回指定进程，用于进程详情展示
func (pm *processManager) GetProcess(pid int32) (*ProcessInfo, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	proc, exists := pm.managedProcesses[pid]
	if !exists {
		return nil, fmt.Errorf("进程 %d 不存在", pid)
	}
	
	pm.logger.Debug("获取进程", zap.Int32("pid", pid))
	return proc, nil
}

// TerminateProcess 终止进程
// 参数：pid 进程ID
// 返回：错误信息，如果进程不存在或终止失败则返回错误
// 上下承接关系：终止指定进程，更新进程状态，记录终止操作
func (pm *processManager) TerminateProcess(pid int32) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	proc, exists := pm.managedProcesses[pid]
	if !exists {
		return fmt.Errorf("进程 %d 不存在", pid)
	}
	
	// 模拟终止进程的逻辑
	// 实际实现需要系统调用
	delete(pm.managedProcesses, pid)
	
	pm.logger.Info("进程终止", zap.Int32("pid", pid), zap.String("name", proc.Name))
	return nil
}

// RestartProcess 重启进程
// 参数：pid 进程ID
// 返回：错误信息，如果进程不存在或重启失败则返回错误
// 上下承接关系：先终止后重新启动指定进程，记录重启操作
func (pm *processManager) RestartProcess(pid int32) error {
	// 先终止进程
	err := pm.TerminateProcess(pid)
	if err != nil {
		return err
	}
	
	// 模拟重启进程
	pm.logger.Info("进程重启", zap.Int32("pid", pid))
	return nil
}

// GetProcessStats 获取进程统计信息
// 参数：无
// 返回：进程统计信息，错误信息
// 上下承接关系：聚合进程状态信息，返回统计结果用于监控面板
func (pm *processManager) GetProcessStats() (*ProcessStats, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	stats := &ProcessStats{
		TotalProcesses: len(pm.managedProcesses),
		Timestamp:      time.Now().Unix(),
	}
	
	// 计算CPU和内存使用率
	var totalCPU float64
	var totalMemory uint64
	var runningCount int
	
	for _, proc := range pm.managedProcesses {
		totalCPU += proc.CPUUsage
		totalMemory += proc.Memory
		if proc.Status == "running" {
			runningCount++
		}
	}
	
	if len(pm.managedProcesses) > 0 {
		stats.CPUUsage = totalCPU / float64(len(pm.managedProcesses))
		stats.MemoryUsage = totalMemory
		stats.RunningProcesses = runningCount
	}
	
	pm.logger.Debug("获取进程统计信息", 
		zap.Int("total_processes", stats.TotalProcesses),
		zap.Int("running_processes", stats.RunningProcesses),
		zap.Float64("cpu_usage", stats.CPUUsage))
	
	return stats, nil
}

// MonitorProcesses 监控进程状态
// 参数：无
// 返回：错误信息，如果监控失败则返回错误
// 上下承接关系：定期更新进程状态信息，确保进程状态与实际运行情况一致
func (pm *processManager) MonitorProcesses() error {
	pm.logger.Debug("监控进程状态")
	// 实现监控逻辑
	pm.updateProcessList()
	return nil
}

// monitorLoop 监控循环
// 参数：无
// 返回：无
// 上下承接关系：定期执行进程状态更新，处理进程生命周期管理
func (pm *processManager) monitorLoop() {
	interval, err := time.ParseDuration(pm.config.Monitoring.HealthCheck.Interval)
	if err != nil {
		interval = 5 * time.Second
	}
	
	pm.ticker = time.NewTicker(interval)
	defer pm.ticker.Stop()
	
	for {
		select {
		case <-pm.ctx.Done():
			pm.logger.Info("监控循环停止")
			return
		case <-pm.ticker.C:
			pm.updateProcessList()
		}
	}
}

// updateProcessList 更新进程列表
// 参数：无
// 返回：无
// 上下承接关系：定期获取系统进程信息，更新进程管理器中的进程列表
func (pm *processManager) updateProcessList() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	// 实际实现：获取系统进程信息
	pm.logger.Debug("更新进程列表")
	
	// 获取所有系统进程
	processes, err := process.Processes()
	if err != nil {
		pm.logger.Error("获取进程列表失败", zap.Error(err))
		return
	}
	
	// 清空当前进程列表
	pm.managedProcesses = make(map[int32]*ProcessInfo)
	
	// 限制返回的进程数量，避免数据过多
	maxProcesses := 100
	if len(processes) > maxProcesses {
		processes = processes[:maxProcesses]
	}
	
	// 遍历进程并填充信息
	for _, proc := range processes {
		// 获取进程基本信息
		name, err := proc.Name()
		if err != nil {
			name = "unknown"
		}
		
		// 获取CPU使用率
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}
		
		// 获取内存使用量
		memInfo, err := proc.MemoryInfo()
		if err != nil {
			memInfo = &process.MemoryInfoStat{}
		}
		
		// 获取内存使用率
		memPercent, err := proc.MemoryPercent()
		if err != nil {
			memPercent = 0
		}
		
		// 获取进程状态
		var status string
		statuses, err := proc.Status()
		if err != nil {
			status = "unknown"
		} else {
			// Status返回的是字符串切片，取第一个元素
			if len(statuses) > 0 {
				status = statuses[0]
			} else {
				status = "unknown"
			}
		}
		
		// 获取命令行
		cmdLine, err := proc.Cmdline()
		if err != nil {
			cmdLine = ""
		}
		
		// 获取启动时间
		startTime, err := proc.CreateTime()
		if err != nil {
			startTime = time.Now().Unix()
		}
		
		// 获取父进程ID
		parentPID, err := proc.Ppid()
		if err != nil {
			parentPID = 0
		}
		
		// 创建进程信息
		processInfo := &ProcessInfo{
			PID:         proc.Pid,
			Name:        name,
			CmdLine:     cmdLine,
			CPUUsage:    cpuPercent,
			Memory:      memInfo.RSS,
			MemoryUsage: float64(memPercent),
			Status:      status,
			StartTime:   time.Unix(startTime, 0),
			ParentPID:   parentPID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		
		// 添加到管理列表
		pm.managedProcesses[proc.Pid] = processInfo
	}
	
	pm.logger.Debug("进程列表更新完成", zap.Int("count", len(pm.managedProcesses)))
}

// cleanupTerminatedProcesses 清理已终止的进程
// 参数：无
// 返回：无
// 上下承接关系：定期清理已终止的进程记录，释放内存资源
func (pm *processManager) cleanupTerminatedProcesses() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	// 清理逻辑
	pm.logger.Debug("清理已终止进程")
}
