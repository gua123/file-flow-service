package web

import (
	"encoding/json"
	"file-flow-service/internal/service"
	"file-flow-service/internal/service/api"
	"net/http"

	"github.com/go-chi/chi"
	"file-flow-service/utils/logger"
)

type webInterface struct {
	logger  *logger.Logger
	service *service.Service
	// other fields
}

// sendErrorResponse 发送错误响应
// 用于统一处理错误响应，避免代码重复
func (w *webInterface) sendErrorResponse(writer http.ResponseWriter, code int, message string) {
	w.logger.LogError(message)
	w.sendJSONResponse(writer, code, message, nil)
}

// sendJSONResponse 发送JSON响应
func (w *webInterface) sendJSONResponse(writer http.ResponseWriter, code int, message string, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	response := map[string]interface{}{
		"message": message,
		"data":    data,
	}
	json.NewEncoder(writer).Encode(response)
}

// UpdateTaskHandler 处理更新任务请求
func (w *webInterface) UpdateTaskHandler(writer http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		w.sendErrorResponse(writer, http.StatusBadRequest, "任务ID不能为空")
		return
	}
	var req api.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.sendErrorResponse(writer, http.StatusBadRequest, "无效的请求体: "+err.Error())
		return
	}
	if err := w.service.UpdateTask(taskID, req); err != nil {
		w.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	w.sendJSONResponse(writer, http.StatusOK, "任务更新成功", nil)
}

// DeleteTaskHandler 处理删除任务请求
func (w *webInterface) DeleteTaskHandler(writer http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		w.sendErrorResponse(writer, http.StatusBadRequest, "任务ID不能为空")
		return
	}
	if err := w.service.DeleteTask(taskID); err != nil {
		w.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	w.sendJSONResponse(writer, http.StatusOK, "任务删除成功", nil)
}

// GetLogsHandler 处理获取日志请求
func (w *webInterface) GetLogsHandler(writer http.ResponseWriter, r *http.Request) {
	logType := r.URL.Query().Get("logType")
	since := r.URL.Query().Get("since") // 同步修改参数名保持一致性
	if logType == "" || since == "" {
		w.sendErrorResponse(writer, http.StatusBadRequest, "缺少日志类型或时间范围参数")
		return
	}
	logs, err := w.service.GetLogs(logType, since)
	if err != nil {
		w.sendErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}
	w.sendJSONResponse(writer, http.StatusOK, "日志获取成功", logs)
}

// GetProcessListHandler 处理获取进程列表请求
func (w *webInterface) GetProcessListHandler(writer http.ResponseWriter, r *http.Request) {
	processes, err := w.service.GetProcessList()
	if err != nil {
		w.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	w.sendJSONResponse(writer, http.StatusOK, "进程列表获取成功", processes)
}

// UpdateConfigHandler 处理更新配置请求
func (w *webInterface) UpdateConfigHandler(writer http.ResponseWriter, r *http.Request) {
	type ConfigUpdateRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	var req ConfigUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.sendErrorResponse(writer, http.StatusBadRequest, "无效的请求体")
		return
	}
	if err := w.service.UpdateConfig(req.Key, req.Value); err != nil {
		w.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	w.sendJSONResponse(writer, http.StatusOK, "配置更新成功", nil)
}

// DownloadFileHandler 处理文件下载请求
func (w *webInterface) DownloadFileHandler(writer http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "fileID")
	if fileID == "" {
		w.sendErrorResponse(writer, http.StatusBadRequest, "文件ID不能为空")
		return
	}
	filePath, err := w.service.DownloadFile(fileID)
	if err != nil {
		w.sendErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	w.sendJSONResponse(writer, http.StatusOK, "文件下载成功", map[string]string{"path": filePath})
}
