package initializer

import (
	"database/sql"
	"log"
	"os"
	"path"
)

// InitApp executes all initialization operations
func InitApp() error {
	// Create log directories
	logDirs := []string{
		path.Join("log", "environment"),
		path.Join("log", "execution"),
		path.Join("log", "executor"),
		path.Join("log", "file"),
		path.Join("log", "flow"),
		path.Join("log", "permission"),
		path.Join("log", "service"),
		path.Join("log", "web"),
	}

	for _, dir := range logDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("创建日志目录失败 %s: %v", dir, err)
			return err
		}
	}

	// Create tasks table
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}
	defer db.Close()

	createTableSQL := `
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT,
    creator TEXT,
    createdAt TEXT,
    assignedTo TEXT,
    description TEXT,
    resultPath TEXT,
    progress INTEGER
)`
	if _, err := db.Exec(createTableSQL); err != nil {
		return err
	}

	log.Println("初始化完成，所有必要资源已就绪")
	return nil
}

// CheckInitializationRequired 检查是否需要进行初始化
func CheckInitializationRequired() bool {
	// 检查数据库文件是否存在
	if _, err := os.Stat("database.db"); os.IsNotExist(err) {
		return true
	}

	// 检查日志目录是否存在
	logDirs := []string{
		"log/environment",
		"log/execution",
		"log/executor",
		"log/file",
		"log/flow",
		"log/permission",
		"log/service",
		"log/web",
	}
	for _, dir := range logDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return true
		}
	}

	// 检查任务表是否存在
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return true // 无法打开数据库
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='tasks'").Scan(&count)
	if err != nil || count == 0 {
		return true
	}

	return false
}