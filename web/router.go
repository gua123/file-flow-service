package web

import (
	"net/http"
	"os"
	"strings"
	"github.com/rs/cors"
)

// SetupAllRoutes 设置所有HTTP路由
func SetupAllRoutes() {
	// 1. 首先处理API路由（优先级高）
	// 正确移除/api/前缀
	apiHandler := http.StripPrefix("/api/", http.DefaultServeMux)
	http.Handle("/api/", apiHandler)
	
	// 2. 处理静态文件和前端路由
	// 对于所有其他请求，先检查是否为静态文件，否则返回前端index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 如果是API路径，让API处理
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// API路径应该由上面的http.Handle处理，这里只是防止冲突
			return
		}
		
		// 检查请求的文件是否存在（包括前端构建文件）
		// 先检查是否为静态资源文件
		if strings.Contains(r.URL.Path, ".") {
			// 如果包含点号，可能是静态文件
			filePath := "./web/file-flow-web/dist/fileflow" + r.URL.Path
			if _, err := os.Stat(filePath); err == nil {
				http.ServeFile(w, r, filePath)
				return
			}
		}
		
		// 对于前端路由，返回前端index.html让Vue Router处理
		// 这样可以避免404错误，让Vue Router处理客户端路由
		frontendIndexPath := "./web/file-flow-web/dist/fileflow/index.html"
		if _, err := os.Stat(frontendIndexPath); err == nil {
			http.ServeFile(w, r, frontendIndexPath)
			return
		}
		
		// 如果找不到前端index.html，返回后端的index.html（用于测试）
		http.ServeFile(w, r, "./web/index.html")
	})
}

// StartServer 启动HTTP服务器
func StartServer() {
	SetupAllRoutes() // 确保路由已经设置
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}).Handler(http.DefaultServeMux)
	http.ListenAndServe(":8080", handler)
}