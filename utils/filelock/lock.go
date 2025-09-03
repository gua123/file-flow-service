package filelock

import (
	"sync"
)

// 文件锁管理器
type FileLockManager struct {
	locks map[string]*sync.Mutex
	mu    sync.RWMutex // 保护锁集合的并发访问
}

// 新建锁管理器
func NewFileLockManager() *FileLockManager {
	return &FileLockManager{
		locks: make(map[string]*sync.Mutex),
	}
}

// 获取文件锁（自动创建）
func (m *FileLockManager) GetLock(filePath string) *sync.Mutex {
	m.mu.RLock()
	lock, exists := m.locks[filePath]
	m.mu.RUnlock()

	if exists {
		return lock
	}

	// 不存在则加写锁创建
	m.mu.Lock()
	defer m.mu.Unlock()
	lock, exists = m.locks[filePath]
	if !exists {
		lock = &sync.Mutex{}
		m.locks[filePath] = lock
	}
	return lock
}

// 释放锁（当文件被删除时调用）
func (m *FileLockManager) ReleaseLock(filePath string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.locks, filePath)
}
