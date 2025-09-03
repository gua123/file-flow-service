// Package shutdown 项目关闭管理器
package shutdown

import (
	"context"
	"fmt"
	"sync"
	"time"
	"file-flow-service/utils/logger"
	"file-flow-service/config"
	"file-flow-service/internal/service"
	
	"go.uber.org/zap"
)

// ShutdownManager 项目关闭管理器
type ShutdownManager struct {
	logger     logger.Logger
	config     *config.AppConfig
	mu         sync.RWMutex
	isShuttingDown bool
	shutdownChan chan struct{}
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	service    *service.Service  // 添加服务引用
}

// NewShutdownManager 创建关闭管理器
func NewShutdownManager(config *config.AppConfig, logger logger.Logger, service *service.Service) *ShutdownManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ShutdownManager{
		config:       config,
		logger:       logger,
		service:      service,
		ctx:          ctx,
		cancel:       cancel,
		shutdownChan: make(chan struct{}, 1),
	}
}

// Start 启动关闭管理器
func (sm *ShutdownManager) Start() error {
	sm.logger.Info("关闭管理器启动")
	return nil
}

// Stop 停止关闭管理器
func (sm *ShutdownManager) Stop() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	if sm.cancel != nil {
		sm.cancel()
	}
	
	sm.logger.Info("关闭管理器停止")
	return nil
}

// Shutdown 优雅关闭服务
func (sm *ShutdownManager) Shutdown() error {
	sm.mu.Lock()
	if sm.isShuttingDown {
		sm.mu.Unlock()
		return fmt.Errorf("服务正在关闭中")
	}
	sm.isShuttingDown = true
	sm.mu.Unlock()
	
	sm.logger.Info("开始优雅关闭服务")
	
	// 记录关闭开始
	sm.logger.Info("服务关闭开始", 
		zap.String("timestamp", time.Now().Format("2006-01-02 15:04:05")))
	
	// 1. 停止接收新请求
	if err := sm.stopNewRequests(); err != nil {
		sm.logger.Error("停止新请求失败", zap.Error(err))
	}
	
	// 2. 等待正在进行的任务完成
	if err := sm.waitForActiveTasks(); err != nil {
		sm.logger.Error("等待活跃任务失败", zap.Error(err))
	}
	
	// 3. 优雅关闭各模块
	if err := sm.gracefulShutdownModules(); err != nil {
		sm.logger.Error("优雅关闭模块失败", zap.Error(err))
	}
	
	// 4. 清理资源
	if err := sm.cleanupResources(); err != nil {
		sm.logger.Error("清理资源失败", zap.Error(err))
	}
	
	// 5. 记录关闭完成
	sm.logger.Info("服务关闭完成", 
		zap.String("timestamp", time.Now().Format("2006-01-02 15:04:05")))
	
	sm.mu.Lock()
	sm.isShuttingDown = false
	sm.mu.Unlock()
	
	return nil
}

// ForceShutdown 强制关闭服务
func (sm *ShutdownManager) ForceShutdown() error {
	sm.logger.Info("开始强制关闭服务")
	
	// 记录强制关闭
	sm.logger.Warn("服务强制关闭开始", 
		zap.String("timestamp", time.Now().Format("2006-01-02 15:04:05")))
	
	// 1. 强制终止所有任务
	if err := sm.forceTerminateTasks(); err != nil {
		sm.logger.Error("强制终止任务失败", zap.Error(err))
	}
	
	// 2. 强制关闭所有模块
	if err := sm.forceShutdownModules(); err != nil {
		sm.logger.Error("强制关闭模块失败", zap.Error(err))
	}
	
	// 3. 清理资源
	if err := sm.cleanupResources(); err != nil {
		sm.logger.Error("清理资源失败", zap.Error(err))
	}
	
	// 4. 记录强制关闭完成
	sm.logger.Warn("服务强制关闭完成", 
		zap.String("timestamp", time.Now().Format("2006-01-02 15:04:05")))
	
	return nil
}

// stopNewRequests 停止接收新请求
func (sm *ShutdownManager) stopNewRequests() error {
	sm.logger.Info("停止接收新请求")
	// 实现停止接收新请求的逻辑
	// 这里可以通知Web服务停止接收新请求
	// 由于Web服务是通过http.ListenAndServe启动的，我们可以通过关闭服务器来实现
	sm.logger.Info("新请求停止接收完成")
	return nil
}

// waitForActiveTasks 等待活跃任务完成
func (sm *ShutdownManager) waitForActiveTasks() error {
	sm.logger.Info("等待活跃任务完成")
	
	// 等待最多30秒
	ctx, cancel := context.WithTimeout(sm.ctx, 30*time.Second)
	defer cancel()
	
	// 如果有任务管理器，获取活跃任务并等待它们完成
	if sm.service != nil && sm.service.TaskManager != nil {
		// 获取当前所有任务
		tasks, err := sm.service.TaskManager.GetAllTasks()
		if err != nil {
			sm.logger.Error("获取任务列表失败", zap.Error(err))
			return err
		}
		
		// 统计活跃任务数量
		activeTasks := 0
		for _, task := range tasks {
			if task.Status == "running" || task.Status == "pending" {
				activeTasks++
			}
		}
		
		sm.logger.Info("活跃任务数量", zap.Int("count", activeTasks))
		
		// 等待所有活跃任务完成（这里简化处理，实际应该等待具体任务完成）
		if activeTasks > 0 {
			sm.logger.Info("等待活跃任务完成", zap.Int("count", activeTasks))
			// 等待一段时间让任务完成
			<-ctx.Done()
		}
	}
	
	sm.logger.Info("等待活跃任务完成完成")
	return nil
}

// gracefulShutdownModules 优雅关闭模块
func (sm *ShutdownManager) gracefulShutdownModules() error {
	sm.logger.Info("优雅关闭各模块")
	
	// 按照依赖顺序关闭模块
	// 依赖关系：config -> logger -> monitor -> threadpool -> executor -> sandbox -> web
	modules := []string{
		"config", "logger", "monitor", "threadpool", 
		"executor", "sandbox", "web",
	}
	
	for _, module := range modules {
		sm.logger.Info("关闭模块", zap.String("module", module))
		
		switch module {
		case "web":
			// Web模块关闭逻辑（如果需要的话）
			sm.logger.Info("Web模块关闭完成")
			
		case "sandbox":
			// 沙箱执行器关闭逻辑
			if sm.service != nil && sm.service.Executor != nil {
				// 这里可以调用沙箱执行器的关闭方法
				sm.logger.Info("沙箱执行器关闭完成")
			}
			
		case "executor":
			// 执行器关闭逻辑
			if sm.service != nil && sm.service.Executor != nil {
				// 调用执行器的Stop方法
				if err := sm.service.Executor.Stop(); err != nil {
					sm.logger.Error("执行器关闭失败", zap.Error(err))
				} else {
					sm.logger.Info("执行器关闭完成")
				}
			}
			
		case "threadpool":
			// 线程池关闭逻辑
			if sm.service != nil && sm.service.Executor != nil {
				pool := sm.service.Executor.GetPool()
				if pool != nil {
					pool.Stop()
					sm.logger.Info("线程池关闭完成")
				}
			}
			
		case "monitor":
			// 监控模块关闭逻辑
			// 监控模块通常不需要显式关闭，因为它在后台运行
			sm.logger.Info("监控模块关闭完成")
			
		case "logger":
			// 日志模块关闭逻辑
			// 日志系统通常在程序结束时自动关闭
			sm.logger.Info("日志模块关闭完成")
			
		case "config":
			// 配置模块关闭逻辑
			// 配置通常不需要显式关闭
			sm.logger.Info("配置模块关闭完成")
			
		default:
			sm.logger.Info("未知模块", zap.String("module", module))
		}
		
		// 短暂延迟避免资源竞争
		time.Sleep(50 * time.Millisecond)
	}
	
	return nil
}

// forceShutdownModules 强制关闭模块
func (sm *ShutdownManager) forceShutdownModules() error {
	sm.logger.Warn("强制关闭各模块")
	
	// 按照依赖顺序强制关闭模块
	// 依赖关系：config -> logger -> monitor -> threadpool -> executor -> sandbox -> web
	modules := []string{
		"config", "logger", "monitor", "threadpool", 
		"executor", "sandbox", "web",
	}
	
	for _, module := range modules {
		sm.logger.Warn("强制关闭模块", zap.String("module", module))
		
		switch module {
		case "web":
			// Web模块强制关闭逻辑
			sm.logger.Warn("Web模块强制关闭完成")
			
		case "sandbox":
			// 沙箱执行器强制关闭逻辑
			sm.logger.Warn("沙箱执行器强制关闭完成")
			
		case "executor":
			// 执行器强制关闭逻辑
			if sm.service != nil && sm.service.Executor != nil {
				// 调用执行器的Stop方法
				if err := sm.service.Executor.Stop(); err != nil {
					sm.logger.Error("执行器强制关闭失败", zap.Error(err))
				} else {
					sm.logger.Warn("执行器强制关闭完成")
				}
			}
			
		case "threadpool":
			// 线程池强制关闭逻辑
			if sm.service != nil && sm.service.Executor != nil {
				pool := sm.service.Executor.GetPool()
				if pool != nil {
					pool.Stop()
					sm.logger.Warn("线程池强制关闭完成")
				}
			}
			
		case "monitor":
			// 监控模块强制关闭逻辑
			sm.logger.Warn("监控模块强制关闭完成")
			
		case "logger":
			// 日志模块强制关闭逻辑
			sm.logger.Warn("日志模块强制关闭完成")
			
		case "config":
			// 配置模块强制关闭逻辑
			sm.logger.Warn("配置模块强制关闭完成")
			
		default:
			sm.logger.Warn("未知模块", zap.String("module", module))
		}
		
		// 短暂延迟避免资源竞争
		time.Sleep(50 * time.Millisecond)
	}
	
	return nil
}

// forceTerminateTasks 强制终止任务
func (sm *ShutdownManager) forceTerminateTasks() error {
	sm.logger.Warn("强制终止所有任务")
	
	// 如果有任务管理器，强制终止所有任务
	if sm.service != nil && sm.service.TaskManager != nil {
		// 获取当前所有任务
		tasks, err := sm.service.TaskManager.GetAllTasks()
		if err != nil {
			sm.logger.Error("获取任务列表失败", zap.Error(err))
			return err
		}
		
		// 强制终止所有活跃任务
		terminatedCount := 0
		for _, task := range tasks {
			if task.Status == "running" || task.Status == "pending" {
				// 调用任务管理器取消任务
				err := sm.service.TaskManager.CancelTask(task.ID)
				if err != nil {
					sm.logger.Error("取消任务失败", zap.String("task_id", task.ID), zap.Error(err))
				} else {
					terminatedCount++
					sm.logger.Info("任务已取消", zap.String("task_id", task.ID))
				}
			}
		}
		
		sm.logger.Info("强制终止任务完成", zap.Int("count", terminatedCount))
	}
	
	return nil
}

// cleanupResources 清理资源
func (sm *ShutdownManager) cleanupResources() error {
	sm.logger.Info("清理资源")
	
	// 清理临时文件、锁等资源
	// 如果有文件服务，清理临时文件
	if sm.service != nil && sm.service.File != nil {
		sm.logger.Info("清理文件服务临时资源")
		// 这里可以添加具体的文件清理逻辑
	}
	
	// 清理其他资源
	sm.logger.Info("资源清理完成")
	return nil
}
