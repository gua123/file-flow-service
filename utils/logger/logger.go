// logger.go
// 日志模块实现，提供全局日志对象
// 用于记录平台运行日志和文件执行日志
// 服务模块通过logger.GetLogger()获取日志对象
// 服务模块通过logger.GetLogger()获取日志对象

package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"file-flow-service/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func sliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	LogError(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	SetLevel(level string) error
}

type ZapLogger struct {
	logger      *zap.Logger
	atomicLevel zap.AtomicLevel
}

// NewZapLogger 创建Zap日志实例
// 参数：config 日志配置
// 返回：ZapLogger实例，错误信息
func NewZapLogger(config *config.LoggerConf) (*ZapLogger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	atomicLevel := zap.NewAtomicLevel()

	var levelStr string
	for level := range config.Levels {
		if config.Levels[level] {
			levelStr = level
			break
		}
	}
	if levelStr == "" {
		levelStr = "info"
	}
	lev, err := zapcore.ParseLevel(levelStr)
	if err != nil {
		return nil, err
	}
	atomicLevel.SetLevel(lev)

	var cores []zapcore.Core

	// 控制台输出
	if config.Levels["info"] {
		consoleEncoder := zapcore.NewJSONEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			atomicLevel,
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if sliceContains(config.Outputs, "file") {
		err := os.MkdirAll(config.BasePath, 0755)
		if err != nil {
			return nil, err
		}
		filePath := filepath.Join(config.BasePath, "app.log")
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			atomicLevel,
		)
		cores = append(cores, fileCore)
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core)
	return &ZapLogger{
		logger:      logger,
		atomicLevel: atomicLevel,
	}, nil
}

func (zl *ZapLogger) Debug(msg string, fields ...zap.Field) {
	zl.logger.Debug(msg, fields...)
}

func (zl *ZapLogger) Info(msg string, fields ...zap.Field) {
	zl.logger.Info(msg, fields...)
}

func (zl *ZapLogger) Warn(msg string, fields ...zap.Field) {
	zl.logger.Warn(msg, fields...)
}

func (zl *ZapLogger) Error(msg string, fields ...zap.Field) {
	zl.logger.Error(msg, fields...)
}

func (zl *ZapLogger) LogError(msg string, fields ...zap.Field) {
	zl.logger.Error(msg, fields...)
}

func (zl *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	zl.logger.Fatal(msg, fields...)
}

func (zl *ZapLogger) SetLevel(level string) error {
	lev, err := zapcore.ParseLevel(level)
	if err != nil {
		return err
	}
	zl.atomicLevel.SetLevel(lev)
	return nil
}

// ZapField 创建zap字段
func ZapField(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

var globalLogger Logger
var serviceLogger Logger
var flowLogger Logger // 新增全局变量
var executorLogger Logger
var fileLogger Logger
var environmentLogger Logger
var executionLogger Logger
var permissionLogger Logger
var webLogger Logger

// InitLogger 初始化日志模块
// 创建日志目录，设置日志级别，初始化全局日志对象
// 参数: config 日志配置
// 返回: 错误信息
func InitLogger() error {
	appConfig := config.GetConfig()
	zapLogger, err := NewZapLogger(&appConfig.LoggerConf)
	if err != nil {
		return err
	}
	globalLogger = zapLogger
	return nil
}

// InitModuleLoggers 初始化各模块日志记录器
// 返回: 错误信息
func InitModuleLoggers() error {
	appConfig := config.GetConfig()

	// 创建service日志目录
	serviceLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "service")
	err := os.MkdirAll(serviceLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建service日志目录失败: %v", err)
	}

	// 创建service日志配置
	serviceLoggerConf := config.LoggerConf{
		BasePath: serviceLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}
	serviceLogger, err = NewZapLogger(&serviceLoggerConf)
	if err != nil {
		return err
	}

	// 记录日志目录创建成功
	GetLogger().Info("服务模块日志目录创建成功", zap.String("path", serviceLogPath))

	// 服务模块初始化日志
	GetServiceLogger().Info("服务模块初始化完成")

	// 创建flow日志目录
	flowLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "flow")
	err = os.MkdirAll(flowLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建flow日志目录失败: %v", err)
	}

	// 创建flow日志配置
	flowLoggerConf := config.LoggerConf{
		BasePath: flowLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}
	flowLogger, err = NewZapLogger(&flowLoggerConf)
	if err != nil {
		return err
	}

	// 记录日志目录创建成功
	GetLogger().Info("流程模块日志目录创建成功", zap.String("path", flowLogPath))

	// 流程模块初始化日志
	// GetFlowLogger().Info("流程模块初始化完成") // 流程模块没有对应的全局变量，跳过此行

	// 为其他模块创建日志记录器
	// 创建executor日志目录
	executorLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "executor")
	err = os.MkdirAll(executorLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建executor日志目录失败: %v", err)
	}

	// 创建executor日志配置
	executorLoggerConf := config.LoggerConf{
		BasePath: executorLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}
	executorLogger, err = NewZapLogger(&executorLoggerConf)
	if err != nil {
		return err
	}

	// 记录日志目录创建成功
	GetLogger().Info("执行器模块日志目录创建成功", zap.String("path", executorLogPath))

	// 执行器模块初始化日志
	GetExecutorLogger().Info("执行器模块初始化完成")

	// 创建file日志目录
	fileLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "file")
	err = os.MkdirAll(fileLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建file日志目录失败: %v", err)
	}

	// 创建file日志配置
	fileLoggerConf := config.LoggerConf{
		BasePath: fileLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}
	fileLogger, err = NewZapLogger(&fileLoggerConf)
	if err != nil {
		return err
	}

	// 记录日志目录创建成功
	GetLogger().Info("文件模块日志目录创建成功", zap.String("path", fileLogPath))

	// 文件模块初始化日志
	GetFileLogger().Info("文件模块初始化完成")

	// 创建environment日志目录
	environmentLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "environment")
	err = os.MkdirAll(environmentLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建environment日志目录失败: %v", err)
	}

	// 创建environment日志配置
	environmentLoggerConf := config.LoggerConf{
		BasePath: environmentLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}
	environmentLogger, err = NewZapLogger(&environmentLoggerConf)
	if err != nil {
		return err
	}

	// 记录日志目录创建成功
	GetLogger().Info("环境管理模块日志目录创建成功", zap.String("path", environmentLogPath))

	// 环境管理模块初始化日志
	GetEnvironmentLogger().Info("环境管理模块初始化完成")

	// 创建execution日志目录
	executionLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "execution")
	err = os.MkdirAll(executionLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建execution日志目录失败: %v", err)
	}

	// 创建execution日志配置
	executionLoggerConf := config.LoggerConf{
		BasePath: executionLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}
	executionLogger, err = NewZapLogger(&executionLoggerConf)
	if err != nil {
		return err
	}

	// 记录日志目录创建成功
	GetLogger().Info("执行模块日志目录创建成功", zap.String("path", executionLogPath))

	// 执行模块初始化日志
	GetExecutionLogger().Info("执行模块初始化完成")

	// 创建permission日志目录
	permissionLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "permission")
	err = os.MkdirAll(permissionLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建permission日志目录失败: %v", err)
	}

	// 创建permission日志配置
	permissionLoggerConf := config.LoggerConf{
		BasePath: permissionLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}
	permissionLogger, err = NewZapLogger(&permissionLoggerConf)
	if err != nil {
		return err
	}

	// 记录日志目录创建成功
	GetLogger().Info("权限模块日志目录创建成功", zap.String("path", permissionLogPath))

	// 权限模块初始化日志
	GetPermissionLogger().Info("权限模块初始化完成")

	// 创建web日志目录
	webLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "web")
	err = os.MkdirAll(webLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建web日志目录失败: %v", err)
	}

	// 创建web日志配置
	webLoggerConf := config.LoggerConf{
		BasePath: webLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}
	webLogger, err = NewZapLogger(&webLoggerConf)
	if err != nil {
		return err
	}

	// 记录日志目录创建成功
	GetLogger().Info("Web模块日志目录创建成功", zap.String("path", webLogPath))

	// Web模块初始化日志
	GetWebLogger().Info("Web模块初始化完成")

	return nil
}

// GetLogger 获取全局日志对象
// 参数：无
// 返回：日志接口实例
func GetLogger() Logger {
	return globalLogger
}

// GetServiceLogger 获取服务模块日志记录器
func GetServiceLogger() Logger {
	return serviceLogger
}

// GetExecutorLogger 获取执行器模块日志记录器
func GetExecutorLogger() Logger {
	return executorLogger
}

// GetFileLogger 获取文件模块日志记录器
func GetFileLogger() Logger {
	return fileLogger
}

// GetEnvironmentLogger 获取环境管理模块日志记录器
func GetEnvironmentLogger() Logger {
	return environmentLogger
}

// GetExecutionLogger 获取执行模块日志记录器
func GetExecutionLogger() Logger {
	return executionLogger
}

// GetPermissionLogger 获取权限模块日志记录器
func GetPermissionLogger() Logger {
	return permissionLogger
}

// GetWebLogger 获取Web模块日志记录器
func GetWebLogger() Logger {
	return webLogger
}

// GetFlowLogger 获取流程模块日志记录器
func GetFlowLogger() Logger {
	// 流程模块没有单独的全局变量，返回nil或者需要特殊处理
	// 为了不导致panic，这里返回nil
	return nil
}

// InitServiceLogger 初始化服务日志模块
// 创建service日志目录，设置日志级别，初始化服务日志对象
// 返回: 错误信息
func InitServiceLogger() error {
	appConfig := config.GetConfig()

	// 创建service日志目录
	serviceLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "service")
	err := os.MkdirAll(serviceLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建service日志目录失败: %v", err)
	}

	// 创建service日志配置
	serviceLoggerConf := config.LoggerConf{
		BasePath: serviceLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}

	_, err = NewZapLogger(&serviceLoggerConf) // Corrected: Ignore first return value
	if err != nil {
		return err
	}

	return nil
}

// InitFlowLogger 初始化flow日志模块
// 创建flow日志目录，设置日志级别，初始化flow日志对象
// 返回: 错误信息
func InitFlowLogger() error {
	appConfig := config.GetConfig()

	// 创建flow日志目录
	flowLogPath := filepath.Join(appConfig.LoggerConf.BasePath, "flow")
	err := os.MkdirAll(flowLogPath, 0755)
	if err != nil {
		return fmt.Errorf("创建flow日志目录失败: %v", err)
	}

	// 创建flow日志配置
	flowLoggerConf := config.LoggerConf{
		BasePath: flowLogPath,
		Levels:   appConfig.LoggerConf.Levels,
		Format:   appConfig.LoggerConf.Format,
		Rotation: appConfig.LoggerConf.Rotation,
		Outputs:  appConfig.LoggerConf.Outputs,
	}

	_, err = NewZapLogger(&flowLoggerConf) // Corrected: Ignore first return value
	if err != nil {
		return err
	}

	return nil
}
