// Package monitor 性能监控模块，提供硬件和任务监控功能
// 该模块负责监控系统性能、任务执行状态和资源使用情况
// 与 service 模块紧密协作，为系统提供实时监控数据
package monitor

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
	"bufio"
	"regexp"
	"strings"

	"file-flow-service/config"
	"file-flow-service/internal/processmanager"
	"file-flow-service/internal/service/api"
	"file-flow-service/utils/logger"

	"go.uber.org/zap"
	
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// MonitorImpl 监控服务实现
// 实现了监控接口，提供硬件统计、系统信息、进程列表等监控功能
type MonitorImpl struct {
	interval     time.Duration
	logger       logger.Logger
	service      api.Service
	processManager processmanager.ProcessManager
	mu           sync.Mutex // 用于保护监控配置的并发访问
}

// NewMonitor 创建监控实例
// 参数：interval 监控间隔, logger 日志记录器, svc API服务实例
// 返回：监控实例
func NewMonitor(interval time.Duration, logger logger.Logger, svc api.Service) *MonitorImpl {
	return &MonitorImpl{
		interval: interval,
		logger:   logger,
		service:  svc,
	}
}

// Start 启动监控服务
// 参数：ctx 上下文，用于控制监控服务的生命周期
// 返回：无
// 上下承接关系：调用 service.GetStatus() 获取服务状态，调用 GetHardwareStats() 获取硬件统计
func (m *MonitorImpl) Start(ctx context.Context) {
	m.logger.Info("监控服务启动")
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("监控服务停止")
			return
		case <-ticker.C:
			m.logger.Debug("执行监控任务")
			// 获取服务状态
			status := m.service.GetStatus()
			m.logger.Info("服务状态", zap.String("status", status))
			
			// 获取硬件统计信息
			hardwareStats, err := m.GetHardwareStats()
			if err != nil {
				m.logger.Error("获取硬件统计信息失败", zap.Error(err))
			} else {
				m.logger.Info("硬件统计信息", 
					zap.Float64("cpu_usage", hardwareStats.CPUUsage),
					zap.Uint64("memory_used", hardwareStats.MemoryUsed),
					zap.Uint64("disk_used", hardwareStats.DiskUsed))
			}
		}
	}
}

// UpdateInterval 更新监控间隔
// 参数：newValue 新的间隔时间字符串，格式如 "5s", "1m"
// 返回：错误信息，如果解析失败则返回错误
// 上下承接关系：更新内部的监控间隔配置，影响监控频率
func (m *MonitorImpl) UpdateInterval(newValue string) error {
	interval, err := time.ParseDuration(newValue)
	if err != nil {
		return fmt.Errorf("解析监控间隔失败: %v", err)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.interval = interval
	m.logger.Info("监控间隔已更新", zap.Duration("interval", m.interval))
	return nil
}

// UpdateLogLevel 更新日志级别
// 参数：level 日志级别，如 "debug", "info", "warn", "error"
// 返回：错误信息，如果级别无效则返回错误
// 上下承接关系：更新内部日志记录器的日志级别
func (m *MonitorImpl) UpdateLogLevel(level string) error {
	err := m.logger.SetLevel(level)
	if err != nil {
		return fmt.Errorf("设置日志级别失败: %v", err)
	}
	m.logger.Info("日志级别已更新", zap.String("level", level))
	return nil
}

// GetLogs 获取日志
// 参数：logType 日志类型，如 "info", "error", "debug"，since 时间点过滤
// 返回：日志列表，错误信息
// 上下承接关系：调用日志系统获取指定类型和时间范围的日志
func (m *MonitorImpl) GetLogs(logType string, since string) ([]string, error) {
	m.logger.Info("获取日志", zap.String("type", logType), zap.String("since", since))
	// 读取日志文件并过滤
	logs, err := m.readAndFilterLogs(logType, since)
	if err != nil {
		return nil, err
	}
	return logs, nil
}


// readAndFilterLogs 读取日志文件并根据类型和时间过滤
func (m *MonitorImpl) readAndFilterLogs(logType string, since string) ([]string, error) {
	// 解析since时间
	var sinceTime time.Time
	if since != "" {
		var err error
		sinceTime, err = time.Parse("2006-01-02 15:04:05", since)
		if err != nil {
			return nil, fmt.Errorf("无效的since时间格式: %v", err)
		}
	}

	// 打开日志文件
	file, err := os.Open("log/app.log")
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	scanner := bufio.NewScanner(file)
	var logs []string
	for scanner.Scan() {
		line := scanner.Text()
		// 解析时间戳和日志内容
		re := regexp.MustCompile(`\[(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] (.*)`)
		matches := re.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue // 跳过格式不匹配的行
		}
		timestampStr := matches[1]
		logContent := matches[2]

		// 解析时间戳
		logTime, err := time.Parse("2006-01-02 15:04:05", timestampStr)
		if err != nil {
			continue // 跳过解析失败的行
		}

		// 检查时间是否在since之后
		if logTime.Before(sinceTime) {
			continue
		}

		// 检查日志类型
		var logTypeMatch bool
		switch logType {
		case "info":
			logTypeMatch = strings.Contains(logContent, "信息")
		case "error":
			logTypeMatch = strings.Contains(logContent, "错误")
		case "warn":
			logTypeMatch = strings.Contains(logContent, "警告")
		case "debug":
			logTypeMatch = strings.Contains(logContent, "调试")
		default:
			logTypeMatch = true // 不过滤类型
		}

		if logTypeMatch {
			logs = append(logs, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取日志文件失败: %v", err)
	}

	return logs, nil
}

// GetHardwareStats 获取硬件统计信息
// 参数：无
// 返回：硬件统计信息，错误信息
// 上下承接关系：收集CPU、内存、磁盘等系统资源使用情况
func (m *MonitorImpl) GetHardwareStats() (*api.HardwareStats, error) {
	m.logger.Debug("获取硬件统计信息")
	
	// 获取系统资源使用情况
	var stats api.HardwareStats
	
	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		m.logger.Warn("获取CPU使用率失败", zap.Error(err))
		stats.CPUUsage = 0
	} else {
		stats.CPUUsage = cpuPercent[0]
	}
	
	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		m.logger.Warn("获取内存信息失败", zap.Error(err))
		stats.MemoryTotal = 0
		stats.MemoryUsed = 0
		stats.MemoryFree = 0
		stats.MemoryUsage = 0
	} else {
		stats.MemoryTotal = memInfo.Total
		stats.MemoryUsed = memInfo.Used
		stats.MemoryFree = memInfo.Free
		stats.MemoryUsage = memInfo.UsedPercent
	}
	
	// 获取磁盘信息
	diskInfo, err := disk.Usage("/")
	if err != nil {
		m.logger.Warn("获取磁盘信息失败", zap.Error(err))
		stats.DiskTotal = 0
		stats.DiskUsed = 0
		stats.DiskFree = 0
		stats.DiskUsage = 0
	} else {
		stats.DiskTotal = diskInfo.Total
		stats.DiskUsed = diskInfo.Used
		stats.DiskFree = diskInfo.Free
		stats.DiskUsage = diskInfo.UsedPercent
	}
	
	// 获取系统负载
	// 注意：LoadAvg 在某些系统上可能不可用，使用0作为默认值
	stats.LoadAverage = 0
	
	// 获取系统运行时间
	hostInfo, err := host.Info()
	if err != nil {
		m.logger.Warn("获取主机信息失败", zap.Error(err))
		stats.Uptime = 0
	} else {
		stats.Uptime = hostInfo.Uptime
	}
	
	// 获取进程数量
	processes, err := process.Processes()
	if err != nil {
		m.logger.Warn("获取进程列表失败", zap.Error(err))
		stats.ProcessCount = 0
	} else {
		stats.ProcessCount = len(processes)
	}
	
	// 设置时间戳
	stats.Timestamp = time.Now().Unix()
	
	m.logger.Debug("硬件统计信息获取完成", 
		zap.Float64("cpu_usage", stats.CPUUsage),
		zap.Uint64("memory_used", stats.MemoryUsed),
		zap.Uint64("disk_used", stats.DiskUsed))
	
	return &stats, nil
}

// GetSystemInfo 获取系统信息
// 参数：无
// 返回：系统信息，错误信息
// 上下承接关系：获取操作系统、架构、Go版本等系统基本信息
func (m *MonitorImpl) GetSystemInfo() (*api.SystemInfo, error) {
	m.logger.Debug("获取系统信息")
	
	var info api.SystemInfo
	
	// 获取主机名
	info.Hostname, _ = os.Hostname()
	
	// 获取操作系统信息
	info.OS = runtime.GOOS
	info.Platform = runtime.GOARCH
	info.Architecture = runtime.GOARCH
	
	// 获取内核版本
	hostInfo, err := host.Info()
	if err != nil {
		m.logger.Warn("获取主机信息失败", zap.Error(err))
		info.Kernel = "Unknown"
	} else {
		info.Kernel = hostInfo.KernelVersion
	}
	
	// 获取Go版本
	info.GoVersion = runtime.Version()
	
	// 获取启动时间
	hostInfo2, err := host.Info()
	if err != nil {
		m.logger.Warn("获取主机信息失败", zap.Error(err))
		info.StartTime = time.Now().Unix()
	} else {
		info.StartTime = int64(hostInfo2.BootTime)
	}
	
	m.logger.Debug("系统信息获取完成", 
		zap.String("hostname", info.Hostname),
		zap.String("os", info.OS),
		zap.String("go_version", info.GoVersion))
	
	return &info, nil
}

// GetProcessList 获取进程列表
// 参数：无
// 返回：进程列表，错误信息
// 上下承接关系：调用进程管理模块获取当前运行的进程信息
func (m *MonitorImpl) GetProcessList() ([]*api.ProcessInfo, error) {
	m.logger.Debug("获取进程列表")
	
	// 获取所有进程
	processes, err := process.Processes()
	if err != nil {
		m.logger.Error("获取进程列表失败", zap.Error(err))
		return nil, err
	}
	
	// 转换为API格式
	var result []*api.ProcessInfo
	
	// 限制返回的进程数量，避免数据过多
	maxProcesses := 100
	if len(processes) > maxProcesses {
		processes = processes[:maxProcesses]
	}
	
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
		
		result = append(result, &api.ProcessInfo{
			PID:         proc.Pid,
			Name:        name,
			CPUUsage:    cpuPercent,
			Memory:      memInfo.RSS,
			MemoryUsage: float64(memPercent),
			Status:      status,
			CmdLine:     cmdLine,
		})
	}
	
	m.logger.Debug("进程列表获取完成", zap.Int("count", len(result)))
	
	return result, nil
}

// SetProcessManager 设置进程管理器
// 参数：pm 进程管理器实例
// 返回：无
// 上下承接关系：为监控模块提供进程管理能力
func (m *MonitorImpl) SetProcessManager(pm processmanager.ProcessManager) {
	m.processManager = pm
}

type MonitorStruct struct {
	Config *config.AppConfig
	Logger logger.Logger
	// 其他监控相关字段
}

func NewMonitorStruct(config *config.AppConfig, logger logger.Logger) *MonitorStruct {
	return &MonitorStruct{
		Config: config,
		Logger: logger,
	}
}
