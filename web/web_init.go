package web

import (
    "file-flow-service/config"
    "file-flow-service/utils/logger"
    "net/http"
)

// InitWebModule 初始化Web模块，返回HTTP路由处理器
func InitWebModule(logger logger.Logger, config *config.AppConfig) http.Handler {
	webInterface := NewWebInterface(logger, config)
	return webInterface.SetupRoutes()
}
