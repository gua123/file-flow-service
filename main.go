// main.go
// 服务入口文件
// 初始化配置、日志和启动服务
package main

import (
	"file-flow-service/config"
	"file-flow-service/internal/restart"
	"file-flow-service/internal/service"
	"file-flow-service/web"
	"file-flow-service/utils/logger"
	"file-flow-service/sandbox/environments"
	"file-flow-service/sandbox/execution"
	"log"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	configPath := "config/config.yaml"
	if err := config.InitConfig(configPath); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	appConfig := config.GetConfig()

	// 2. 初始化日志模块
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	
	// 3. 初始化模块日志记录器
	if err := logger.InitModuleLoggers(); err != nil {
		log.Fatalf("模块日志初始化失败: %v", err)
	}
	
	appLogger := logger.GetLogger()

	// 4. 初始化环境管理模块
	envManager := environments.NewEnvironmentManager()
	if err := envManager.Init(appConfig, logger.GetEnvironmentLogger()); err != nil {
		log.Fatalf("环境管理模块初始化失败: %v", err)
	}

	// 5. 初始化沙盒执行模块
	sandboxExecutor := execution.NewSandboxExecutor()
	if err := sandboxExecutor.Init(appConfig, logger.GetExecutionLogger(), envManager); err != nil {
		log.Fatalf("沙盒执行模块初始化失败: %v", err)
	}

	// 6. 创建全局service实例
	serviceInstance := service.NewService(appConfig, logger.GetServiceLogger())

	// 7. 创建重启管理器
	restartManager := restart.NewRestartManager(appConfig, logger.GetLogger(), serviceInstance)
	
	// 启动重启管理器
	if err := restartManager.Start(); err != nil {
		appLogger.Error("重启管理器启动失败", zap.Error(err))
	}

	// 8. 初始化Web模块
	webModule := web.InitWebModule(logger.GetWebLogger(), appConfig)
	
	// 设置路由
	web.SetupAllRoutes(webModule)

	// 启动服务器
	web.StartServer()
}