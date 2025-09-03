// config.go
// 配置模块实现，加载config.yaml并提供全局配置对象
// 通过Config全局变量访问配置参数
// 通过Config全局变量访问配置参数
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
	
	"gopkg.in/yaml.v3"
)

var GlobalConfig *AppConfig

type AppConfig struct {
	mu                   sync.Mutex
	App                  App                  `yaml:"app"`
	LoggerConf           LoggerConf           `yaml:"logger"`
	Threadpool           Threadpool           `yaml:"threadpool"`
	File                 File                 `yaml:"file"`
	Internal             Internal             `yaml:"internal"`
	HotReload            HotReload            `yaml:"hot_reload"`
	CsrfEnabled          bool                 `yaml:"csrf_enabled"`
	Database             Database             `yaml:"database"`
	ThreadpoolMonitoring ThreadpoolMonitoring `yaml:"threadpool_monitoring"`
	Sandbox              Sandbox              `yaml:"sandbox"`
	Monitoring           Monitoring           `yaml:"monitoring"`
	Clients              Clients              `yaml:"clients"`
	Dependencies         Dependencies         `yaml:"dependencies"`
	EnvOverrides         EnvOverrides         `yaml:"env_overrides"`
	Permissions          []string             `json:"permissions"`
	FileManagement       FileManagement       `yaml:"file_management"`
	MonitorInterval      string               `yaml:"monitor_interval"`
	Web                  Web                  `yaml:"web"`
	Logging              Logging              `yaml:"logging"`
	History              []*ConfigSnapshot    `yaml:"-"`
	current              *ConfigSnapshot
	previous             *ConfigSnapshot
	initialized          bool
}

type ConfigSnapshot struct {
	MonitorInterval string
	LoggerLevels    map[string]bool
	MaxWorkers      int
	MaxQueue        int
	MemoryLimit     string
	TaskTimeout     string
	Port            int
	BaseURL         string
	StoragePath     string
	MaxUploadSize   int64
	AllowedPaths    []string
}

type App struct {
	Port               int    `yaml:"port"`
	Name               string `yaml:"name"`
	Env                string `yaml:"env"`
	Version            string `yaml:"version"`
	EnableEnvOverrides bool   `yaml:"enable_env_overrides"`
	BaseURL            string `yaml:"base_url"`
}

type LoggerConf struct {
	BasePath string          `yaml:"base_path"`
	Levels   map[string]bool `yaml:"levels"`
	Format   string          `yaml:"format"`
	Rotation Rotation        `yaml:"rotation"`
	Outputs  []string        `yaml:"outputs"`
}

type Rotation struct {
	MaxAgeDays   int    `yaml:"max_age_days"`
	MaxSizeMB    int    `yaml:"max_size_mb"`
	TimeInterval string `yaml:"time_interval"`
}

type File struct {
	StoragePath   string `yaml:"storage_path"`
	MaxUploadSize int64  `yaml:"max_upload_size"`
}

type Internal struct {
	Service Service `yaml:"service"`
	Monitor Monitor `yaml:"monitor"`
}

type Service struct {
	SandboxTimeout   int `yaml:"sandbox_timeout"`
	MaxParallelTasks int `yaml:"max_parallel_tasks"`
}

type Monitor struct {
	HeartbeatInterval int `yaml:"heartbeat_interval"`
}

type HotReload struct {
	Enabled        bool     `yaml:"enabled"`
	UpdateEndpoint string   `yaml:"update_endpoint"`
	AllowedPaths   []string `yaml:"allowed_paths"`
}

type Threadpool struct {
	MaxWorkers  int    `yaml:"max_workers"`
	MaxQueue    int    `yaml:"max_queue"`
	TaskTimeout string `yaml:"task_timeout"`
	AutoScale   bool   `yaml:"auto_scale"`
	MinWorkers  int    `yaml:"min_workers"`
	MemoryLimit string `yaml:"memory_limit"`
}

type ThreadpoolMonitoring struct {
	StatsInterval string `yaml:"stats_interval"`
}

type Web struct {
	Middleware Middleware `yaml:"middleware"`
	Routes     Routes     `yaml:"routes"`
}

type Middleware struct {
}

type Routes struct {
	HealthCheck  string `yaml:"health_check"`
	TaskEndpoint string `yaml:"task_endpoint"`
}

type Sandbox struct {
	Isolation        Isolation      `yaml:"isolation"`
	ResourceLimits   ResourceLimits `yaml:"resource_limits"`
	ExecutionTimeout string         `yaml:"execution_timeout"`
	Environments     Environments   `yaml:"environments"`
	Execution        Execution      `yaml:"execution"`
}

type Environments struct {
	BasePath string `yaml:"base_path"`
	Python   Python `yaml:"python"`
	Java     Java   `yaml:"java"`
}

type Python struct {
	BasePath       string `yaml:"base_path"`
	InstallersPath string `yaml:"installers_path"`
	VersionsPath   string `yaml:"versions_path"`
}

type Java struct {
	BasePath       string `yaml:"base_path"`
	InstallersPath string `yaml:"installers_path"`
	VersionsPath   string `yaml:"versions_path"`
}

type Execution struct {
	BasePath  string `yaml:"base_path"`
	TasksPath string `yaml:"tasks_path"`
	TempPath  string `yaml:"temp_path"`
	LocksPath string `yaml:"locks_path"`
}

type Isolation struct {
	Chroot bool   `yaml:"chroot"`
	User   string `yaml:"user"`
	Group  string `yaml:"group"`
}

type ResourceLimits struct {
	Memory   string `yaml:"memory"`
	CpuCores int    `yaml:"cpu_cores"`
}

type Monitoring struct {
	StatusPushInterval string             `yaml:"status_push_interval"`
	HealthCheck        HealthCheck        `yaml:"health_check"`
	ResourceThresholds ResourceThresholds `yaml:"resource_thresholds"`
	HardwareMonitoring HardwareMonitoring `yaml:"hardware_monitoring"`
	ProcessMonitoring  ProcessMonitoring  `yaml:"process_monitoring"`
}

type HardwareMonitoring struct {
	Enabled      bool   `yaml:"enabled"`
	Interval     string `yaml:"interval"`
	ProcessLimit int    `yaml:"process_limit"`
}

type ProcessMonitoring struct {
	Enabled      bool   `yaml:"enabled"`
	Interval     string `yaml:"interval"`
	MaxProcesses int    `yaml:"max_processes"`
}

type HealthCheck struct {
	Path     string `yaml:"path"`
	Interval string `yaml:"interval"`
}

type ResourceThresholds struct {
	MemoryUsagePercent int `yaml:"memory_usage_percent"`
	CpuUsagePercent    int `yaml:"cpu_usage_percent"`
}

type Clients struct {
	Web WebClient `yaml:"web"`
}

type WebClient struct {
	MaxConnections int `yaml:"max_connections"`
}

type Dependencies struct {
	Rclone   Rclone   `yaml:"rclone"`
	Template Template `yaml:"template"`
}

type Rclone struct {
	ChunkSize   string `yaml:"chunk_size"`
	Concurrence int    `yaml:"concurrency"`
}

type Template struct {
	StrictMode bool `yaml:"strict_mode"`
}

type EnvOverrides struct {
	Priority    string   `yaml:"priority"`
	AllowedVars []string `yaml:"allowed_vars"`
}

type FileManagement struct {
	TaskDir   TaskDir   `yaml:"task_dir"`
	Results   Results   `yaml:"results"`
	Locks     Locks     `yaml:"locks"`
	BasePaths BasePaths `yaml:"base_paths"`
}

type TaskDir struct {
	CleanupDelay     string `yaml:"cleanup_delay"`
	MaxRetentionDays int    `yaml:"max_retention_days"`
}

type Results struct {
	MaxVersions     int    `yaml:"max_versions"`
	RetentionPolicy string `yaml:"retention_policy"`
}

type Locks struct {
	Timeout string `yaml:"timeout"`
}

type BasePaths struct {
	Tasks   string `yaml:"tasks"`
	Results string `yaml:"results"`
	Logs    string `yaml:"logs"`
}

type Database struct {
	Connection string `yaml:"connection"`
}

type Logging struct {
	RotateSize  int `yaml:"rotate_size"`
	RotateCount int `yaml:"rotate_count"`
}

func (c *AppConfig) LoadConfig(configPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	newCfg, err := loadAndValidateConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	snapshot := &ConfigSnapshot{
		MonitorInterval: newCfg.MonitorInterval,
		LoggerLevels:    newCfg.LoggerConf.Levels,
		MaxWorkers:      newCfg.Threadpool.MaxWorkers,
		MaxQueue:        newCfg.Threadpool.MaxQueue,
		MemoryLimit:     newCfg.Threadpool.MemoryLimit,
		TaskTimeout:     newCfg.Threadpool.TaskTimeout,
		Port:            newCfg.App.Port,
		BaseURL:         newCfg.App.BaseURL,
		StoragePath:     newCfg.File.StoragePath,
		MaxUploadSize:   newCfg.File.MaxUploadSize,
		AllowedPaths:    newCfg.HotReload.AllowedPaths,
	}
	c.current = snapshot
	c.previous = snapshot
	c.initialized = true
	c.History = append(c.History, snapshot)
	return nil
}

func (c *AppConfig) ReloadConfig(configPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	newCfg, err := loadAndValidateConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载新配置失败: %v", err)
	}

	newSnapshot := &ConfigSnapshot{
		MonitorInterval: newCfg.MonitorInterval,
		LoggerLevels:    newCfg.LoggerConf.Levels,
		MaxWorkers:      newCfg.Threadpool.MaxWorkers,
		MaxQueue:        newCfg.Threadpool.MaxQueue,
		MemoryLimit:     newCfg.Threadpool.MemoryLimit,
		TaskTimeout:     newCfg.Threadpool.TaskTimeout,
		Port:            newCfg.App.Port,
		BaseURL:         newCfg.App.BaseURL,
		StoragePath:     newCfg.File.StoragePath,
		MaxUploadSize:   newCfg.File.MaxUploadSize,
		AllowedPaths:    newCfg.HotReload.AllowedPaths,
	}

	c.previous = c.current
	c.current = newSnapshot
	return nil
}

func (c *AppConfig) loadAndValidateConfig(configPath string) (*AppConfig, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	newCfg := &AppConfig{}
	if err := yaml.Unmarshal(content, newCfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	newCfg.LoggerConf.BasePath = filepath.Join(getProjectRoot(), newCfg.LoggerConf.BasePath)
	newCfg.File.StoragePath = filepath.Join(getProjectRoot(), newCfg.File.StoragePath)

	err = os.MkdirAll(newCfg.LoggerConf.BasePath, 0755)
	if err != nil {
		return nil, err
	}

	if err := newCfg.validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	return newCfg, nil
}

func (c *AppConfig) validate() error {
	// Logger验证
	for level := range c.LoggerConf.Levels {
		switch level {
		case "debug", "info", "warn", "error", "critical":
		default:
			return fmt.Errorf("无效的日志级别 %q", level)
		}
	}
	if _, err := time.ParseDuration(c.LoggerConf.Rotation.TimeInterval); err != nil {
		return fmt.Errorf("日志轮转时间间隔 %q 格式不合法: %v", c.LoggerConf.Rotation.TimeInterval, err)
	}

	// Threadpool验证
	if c.Threadpool.MaxWorkers <= 0 {
		return fmt.Errorf("线程池最大工作线程数 %d 不合法", c.Threadpool.MaxWorkers)
	}
	if c.Threadpool.MaxQueue < 0 {
		return fmt.Errorf("任务队列最大容量 %d 不合法", c.Threadpool.MaxQueue)
	}
	if !isValidSize(c.Threadpool.MemoryLimit) {
		return fmt.Errorf("线程池内存限制 %q 格式不合法", c.Threadpool.MemoryLimit)
	}

	// File验证
	if c.File.MaxUploadSize <= 0 {
		return fmt.Errorf("最大上传文件大小 %d 不合法", c.File.MaxUploadSize)
	}

	// Sandbox验证
	if !isValidSize(c.Sandbox.ResourceLimits.Memory) {
		return fmt.Errorf("沙箱内存限制 %q 格式不合法", c.Sandbox.ResourceLimits.Memory)
	}
	if c.Sandbox.ResourceLimits.CpuCores < 1 {
		return fmt.Errorf("CPU核心数 %d 不合法", c.Sandbox.ResourceLimits.CpuCores)
	}

	// Monitoring验证
	if _, err := time.ParseDuration(c.Monitoring.HealthCheck.Interval); err != nil {
		return fmt.Errorf("健康检查间隔 %q 格式不合法: %v", c.Monitoring.HealthCheck.Interval, err)
	}

	// Clients验证
	// 移除桌面客户端验证，因为项目中没有桌面端

	// Dependencies验证
	if !isValidSize(c.Dependencies.Rclone.ChunkSize) {
		return fmt.Errorf("rclone分片大小 %q 格式不合法", c.Dependencies.Rclone.ChunkSize)
	}

	// FileManagement验证
	for _, path := range []string{
		c.FileManagement.BasePaths.Tasks,
		c.FileManagement.BasePaths.Results,
		c.FileManagement.BasePaths.Logs,
	} {
		if _, err := filepath.Abs(path); err != nil {
			return fmt.Errorf("基础路径 %q 不合法: %v", path, err)
		}
	}

	// 其他通用验证
	if c.App.Port <= 0 || c.App.Port > 65535 {
		return fmt.Errorf("端口 %d 不合法", c.App.Port)
	}
	if _, err := time.ParseDuration(c.MonitorInterval); err != nil {
		return fmt.Errorf("监控间隔 %q 格式不合法: %v", c.MonitorInterval, err)
	}
	return nil
}

func getProjectRoot() string {
	dir, _ := os.Getwd()
	return dir
}

func InitConfig(configPath string) error {
	if configPath == "" {
		configPath = "config/config.yaml"
	}
	newCfg, err := loadAndValidateConfig(configPath)
	if err != nil {
		return err
	}
	if GlobalConfig != nil {
		return fmt.Errorf("global config already initialized")
	}
	GlobalConfig = newCfg
	return nil
}

func isValidSize(sizeStr string) bool {
	matched, _ := regexp.MatchString(`^\d+[kKmMgGtTpPeE]?[bB]$`, sizeStr)
	return matched
}

func GetConfig() *AppConfig {
	return GlobalConfig
}

func loadAndValidateConfig(path string) (*AppConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	newCfg := &AppConfig{}
	if err := yaml.Unmarshal(content, newCfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	newCfg.LoggerConf.BasePath = filepath.Join(getProjectRoot(), newCfg.LoggerConf.BasePath)
	newCfg.File.StoragePath = filepath.Join(getProjectRoot(), newCfg.File.StoragePath)

	err = os.MkdirAll(newCfg.LoggerConf.BasePath, 0755)
	if err != nil {
		return nil, err
	}

	if err := newCfg.validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	return newCfg, nil
}

func (c *AppConfig) AllowPath(path string) bool {
	for _, p := range c.HotReload.AllowedPaths {
		matched, _ := filepath.Match(p, path)
		if matched {
			return true
		}
	}
	return false
}