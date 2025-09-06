package initialization

import (
	"database/sql"
	"log"
	"os"
	"path"
)

// IsInitialRun 检查项目是否为首次运行
func IsInitialRun() bool {
	// 检查项目是否为首次运行的逻辑（例如：检查config目录下是否存在标记文件）
	if _, err := os.Stat("config/.initialized"); os.IsNotExist(err) {
		return true
	}
	return false
}

// InitApp 执行所有初始化操作
func InitApp() error {
	// 创建日志目录
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

	// 初始化数据库表结构
	db, err := sql.Open("sqlite3", "database/database.db")
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

	log.Println("项目初始化完成，所有必要资源已就绪")
	return nil
}