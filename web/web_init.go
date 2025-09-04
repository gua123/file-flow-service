package web

import (
	"file-flow-service/config"
	"file-flow-service/internal/service"
	"file-flow-service/utils/logger"
	"net/http"
)

func InitWebModule(logger logger.Logger, config *config.AppConfig) http.Handler {
	service := service.NewService(config, logger)
	webInterface := NewWebInterface(service, logger)
	return webInterface.SetupAllRoutes()
}