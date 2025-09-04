// Package environments 环境管理模块
// 管理各种执行环境（Python、Java等）
package environments

import (
	"fmt"
	"os"
	"path/filepath"
	"file-flow-service/config"
	"file-flow-service/utils/logger"
)

// EnvironmentManager 环境管理器接口
type EnvironmentManager interface {
	// Init 初始化环境管理器
	Init(config *config.AppConfig, logger logger.Logger) error
	
	// GetPythonPath 获取Python环境路径
	GetPythonPath(version string) (string, error)
	
	// GetJavaPath 获取Java环境路径
	GetJavaPath(version string) (string, error)
	
	// InstallEnvironment 安装指定环境
	InstallEnvironment(envType, version, installerPath string) error
	
	// ValidateEnvironment 验证环境是否有效
	ValidateEnvironment(envType, version string) bool
}

// environmentManager 环境管理器实现
type environmentManager struct {
	config *config.AppConfig
	logger logger.Logger
}

// NewEnvironmentManager 创建环境管理器实例
func NewEnvironmentManager() EnvironmentManager {
	return &environmentManager{}
}

// Init 初始化环境管理器
// 参数: config 配置对象, logger 日志对象
// 返回: 错误信息
func (em *environmentManager) Init(config *config.AppConfig, logger logger.Logger) error {
	em.config = config
	em.logger = logger
	
	// 创建环境目录
	environmentsPath := config.Sandbox.Environments.BasePath
	if err := os.MkdirAll(environmentsPath, 0755); err != nil {
		return fmt.Errorf("创建环境目录失败: %v", err)
	}
	
	// 创建Python环境目录
	pythonPath := config.Sandbox.Environments.Python.BasePath
	if err := os.MkdirAll(pythonPath, 0755); err != nil {
		return fmt.Errorf("创建Python环境目录失败: %v", err)
	}
	
	// 创建Java环境目录
	javaPath := config.Sandbox.Environments.Java.BasePath
	if err := os.MkdirAll(javaPath, 0755); err != nil {
		return fmt.Errorf("创建Java环境目录失败: %v", err)
	}
	
	em.logger.Info("环境管理器初始化完成")
	return nil
}

// GetPythonPath 获取Python环境路径
// 参数: version Python版本
// 返回: 环境路径，错误信息
func (em *environmentManager) GetPythonPath(version string) (string, error) {
	if em.config == nil {
		return "", fmt.Errorf("环境管理器未初始化")
	}
	
	pythonVersionsPath := em.config.Sandbox.Environments.Python.VersionsPath
	versionPath := filepath.Join(pythonVersionsPath, version)
	
	// 检查路径是否存在
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		return "", fmt.Errorf("Python版本 %s 不存在", version)
	}
	
	return versionPath, nil
}

// GetJavaPath 获取Java环境路径
// 参数: version Java版本
// 返回: 环境路径，错误信息
func (em *environmentManager) GetJavaPath(version string) (string, error) {
	if em.config == nil {
		return "", fmt.Errorf("环境管理器未初始化")
	}
	
	javaVersionsPath := em.config.Sandbox.Environments.Java.VersionsPath
	versionPath := filepath.Join(javaVersionsPath, version)
	
	// 检查路径是否存在
	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		return "", fmt.Errorf("Java版本 %s 不存在", version)
	}
	
	return versionPath, nil
}

// InstallEnvironment 安装指定环境
// 参数: envType 环境类型, version 版本, installerPath 安装包路径
// 返回: 错误信息
func (em *environmentManager) InstallEnvironment(envType, version, installerPath string) error {
	if em.config == nil {
		return fmt.Errorf("环境管理器未初始化")
	}
	
	// 根据环境类型处理安装
	switch envType {
	case "python":
		return em.installPython(version, installerPath)
	case "java":
		return em.installJava(version, installerPath)
	default:
		return fmt.Errorf("不支持的环境类型: %s", envType)
	}
}

// installPython 安装Python环境
// 参数: version 版本, installerPath 安装包路径
// 返回: 错误信息
func (em *environmentManager) installPython(version, installerPath string) error {
	// 这里应该实现Python安装逻辑
	// 目前只是示例，实际需要根据具体需求实现
	em.logger.Info("安装Python环境 version=" + version + " installer=" + installerPath)
	
	// 实现基本的安装逻辑
	// 1. 检查安装包是否存在
	if _, err := os.Stat(installerPath); os.IsNotExist(err) {
		return fmt.Errorf("安装包不存在: %s", installerPath)
	}
	
	// 2. 创建版本目录
	pythonVersionsPath := em.config.Sandbox.Environments.Python.VersionsPath
	versionPath := filepath.Join(pythonVersionsPath, version)
	
	if err := os.MkdirAll(versionPath, 0755); err != nil {
		return fmt.Errorf("创建Python版本目录失败: %v", err)
	}
	
	// 3. 执行安装（这里只是示例）
	em.logger.Info("Python环境安装完成 version=" + version + " path=" + versionPath)
	return nil
}

// installJava 安装Java环境
// 参数: version 版本, installerPath 安装包路径
// 返回: 错误信息
func (em *environmentManager) installJava(version, installerPath string) error {
	// 这里应该实现Java安装逻辑
	// 目前只是示例，实际需要根据具体需求实现
	em.logger.Info("安装Java环境 version=" + version + " installer=" + installerPath)
	
	// 实现基本的安装逻辑
	// 1. 检查安装包是否存在
	if _, err := os.Stat(installerPath); os.IsNotExist(err) {
		return fmt.Errorf("安装包不存在: %s", installerPath)
	}
	
	// 2. 创建版本目录
	javaVersionsPath := em.config.Sandbox.Environments.Java.VersionsPath
	versionPath := filepath.Join(javaVersionsPath, version)
	
	if err := os.MkdirAll(versionPath, 0755); err != nil {
		return fmt.Errorf("创建Java版本目录失败: %v", err)
	}
	
	// 3. 执行安装（这里只是示例）
	em.logger.Info("Java环境安装完成 version=" + version + " path=" + versionPath)
	return nil
}

// ValidateEnvironment 验证环境是否有效
// 参数: envType 环境类型, version 版本
// 返回: 是否有效
func (em *environmentManager) ValidateEnvironment(envType, version string) bool {
	if em.config == nil {
		return false
	}
	
	var versionPath string
	switch envType {
	case "python":
		versionPath = filepath.Join(em.config.Sandbox.Environments.Python.VersionsPath, version)
	case "java":
		versionPath = filepath.Join(em.config.Sandbox.Environments.Java.VersionsPath, version)
	default:
		return false
	}
	
	// 检查路径是否存在
	_, err := os.Stat(versionPath)
	return !os.IsNotExist(err)
}