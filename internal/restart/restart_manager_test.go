package restart

import (
	"testing"
	"file-flow-service/internal/service/interfaces"
	"file-flow-service/config"
	"file-flow-service/utils/logger"
	"github.com/stretchr/testify/assert"
)

func TestRestartManagerCreation(t *testing.T) {
	// 测试重启管理器创建
	cfg := config.GetConfig()
	logger := logger.GetLogger()
	
	// 创建服务实例（简化测试）
	serviceInstance := &service.Service{}
	
	// 创建重启管理器
	rm := NewRestartManager(cfg, logger, serviceInstance)
	
	// 验证重启管理器创建成功
	assert.NotNil(t, rm)
}

// 测试重启管理器创建功能（避免复杂的依赖初始化问题）
func TestRestartManagerCreationOnly(t *testing.T) {
	// 测试重启管理器创建功能
	cfg := config.GetConfig()
	logger := logger.GetLogger()
	
	// 创建服务实例（简化测试）
	serviceInstance := &service.Service{}
	
	// 创建重启管理器
	rm := NewRestartManager(cfg, logger, serviceInstance)
	
	// 验证重启管理器创建成功
	assert.NotNil(t, rm)
}