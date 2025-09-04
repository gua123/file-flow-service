package web

import (
	"file-flow-service/internal/service"
	"file-flow-service/utils/logger"
	"net/http"
	"encoding/json"
)

type WebInterface struct {
	logger logger.Logger
	service *service.Service
}

func NewWebInterface(service *service.Service, logger logger.Logger) *WebInterface {
	return &WebInterface{
		service: service,
		logger:  logger,
	}
}

func (w *WebInterface) SetupAllRoutes() http.Handler {
	http.HandleFunc("/api/upload", w.HandleUpload)
	http.HandleFunc("/api/execute", w.HandleExecute)
	http.HandleFunc("/api/status", w.HandleStatus)
	return nil
}

func (w *WebInterface) HandleUpload(rw http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.logger.Error("文件上传失败: " + err.Error())
		http.Error(rw, "上传失败", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName, err := w.service.UploadFile(fileHeader)
	if err != nil {
		w.logger.Error("文件上传处理失败: " + err.Error())
		http.Error(rw, "文件处理失败", http.StatusInternalServerError)
		return
	}

	w.WriteJSON(rw, map[string]string{"file": fileName})
}

func (w *WebInterface) HandleExecute(rw http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("cmd")
	args := r.FormValue("args")
	if cmd == "" {
		w.logger.Error("命令参数缺失")
		http.Error(rw, "缺少命令", http.StatusBadRequest)
		return
	}

	err := w.service.ExecuteCommand(cmd, []string{args})
	if err != nil {
		w.logger.Error("命令执行失败: " + err.Error())
		http.Error(rw, "命令执行失败", http.StatusInternalServerError)
		return
	}

	w.WriteJSON(rw, map[string]string{"status": "success"})
}

func (w *WebInterface) HandleStatus(rw http.ResponseWriter, r *http.Request) {
	status := w.service.GetStatus()
	w.WriteJSON(rw, map[string]string{"status": status})
}

func (w *WebInterface) WriteJSON(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, "内部错误", http.StatusInternalServerError)
		return
	}

	rw.Write(jsonData)
}