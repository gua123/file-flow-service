package monitor

import (
	"file-flow-service/config"
	"file-flow-service/utils/logger"
	"context"
)

type MonitorImpl struct {
	logger logger.Logger
	config *config.AppConfig
}

func NewMonitorImpl(logger logger.Logger, config *config.AppConfig) *MonitorImpl {
	return &MonitorImpl{
		logger: logger,
		config: config,
	}
}

func (m *MonitorImpl) Start(ctx context.Context) {
	// 启动监控逻辑
	m.logger.Info("Monitoring started")
}

func (m *MonitorImpl) Stop(ctx context.Context) {
	// 停止监控逻辑
	m.logger.Info("Monitoring stopped")
}