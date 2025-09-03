// Package websocket 提供WebSocket通信功能，用于文件模块的实时交互。
// 该包处理WebSocket连接、消息接收和发送，支持文件状态更新和进度通知。
package web

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader 用于升级HTTP连接到WebSocket。
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleWebSocket 处理WebSocket连接请求。
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// 处理WebSocket消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// 处理消息
		conn.WriteMessage(websocket.TextMessage, []byte("Received: "+string(message)))
	}
}
