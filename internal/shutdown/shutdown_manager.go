package shutdown

import (
	"context"
	"sync"
	"file-flow-service/internal/service/interfaces"
	"file-flow-service/config"
	"file-flow-service/utils/logger"
)

type ShutdownManager struct {
	logger         logger.Logger
	config         *config.AppConfig
	mu             sync.RWMutex
	isShuttingDown bool
	shutdownChan   chan struct{}
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	service        interfaces.Service
}

func NewShutdownManager(service interfaces.Service, logger logger.Logger, config *config.AppConfig) *ShutdownManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &ShutdownManager{
		service:        service,
		logger:         logger,
		config:         config,
		shutdownChan:   make(chan struct{}, 1),
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (sm *ShutdownManager) GracefulShutdown() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.isShuttingDown {
		return
	}
	close(sm.shutdownChan)
	sm.isShuttingDown = true
	sm.wg.Wait()
	sm.cancel()
}

func (sm *ShutdownManager) ForceShutdown() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.isShuttingDown {
		return
	}
	sm.cancel()
	sm.isShuttingDown = true
}