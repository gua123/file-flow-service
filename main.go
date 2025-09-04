// main.go
// 服务入口文件
// 初始化配置、日志和启动服务
package main

import (
	"file-flow-service/config"
	"file-flow-service/internal/restart"
	"file-flow-service/internal/service"
	"file-flow-service/sandbox/environments"
	"file-flow-service/sandbox/execution"
	"file-flow-service/utils/logger"
	"file-flow-service/web"
	"log"
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
	appLogger := logger.GetLogger()

	// 3. 初始化环境管理模块
	envManager := environments.NewEnvironmentManager()
	if err := envManager.Init(appConfig, appLogger); err != nil {
		log.Fatalf("环境管理模块初始化失败: %v", err)
	}

	// 4. 初始化沙盒执行模块
	sandboxExecutor := execution.NewSandboxExecutor()
	if err := sandboxExecutor.Init(appConfig, appLogger, envManager); err != nil {
		log.Fatalf("沙盒执行模块初始化失败: %v", err)
	}

	// 5. 创建全局service实例
	serviceInstance := service.NewService(appConfig, appLogger)

	// 6. 创建重启管理器
	restartManager := restart.NewRestartManager(appConfig, appLogger, serviceInstance)

	// 启动重启管理器
	if err := restartManager.Start(); err != nil {
		appLogger.Error("重启管理器启动失败")
	}

	// 7. 启动服务器
	web.StartServer()
}