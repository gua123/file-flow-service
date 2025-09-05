package testutils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// MockLogger 模拟日志记录器
type MockLogger struct {
	mu      sync.Mutex
	entries []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		entries: make([]string, 0),
	}
}

func (m *MockLogger) Debug(msg string, fields ...zap.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, "DEBUG: "+msg)
}

func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, "INFO: "+msg)
}

func (m *MockLogger) Warn(msg string, fields ...zap.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, "WARN: "+msg)
}

func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, "ERROR: "+msg)
}

func (m *MockLogger) Fatal(msg string, fields ...zap.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, "FATAL: "+msg)
}

func (m *MockLogger) Panic(msg string, fields ...zap.Field) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, "PANIC: "+msg)
}

func (m *MockLogger) Sync() error {
	return nil
}

func (m *MockLogger) Core() zapcore.Core {
	return nil
}

func (m *MockLogger) With(fields ...zap.Field) *zap.Logger {
	return nil
}

func (m *MockLogger) Named(s string) *zap.Logger {
	return nil
}

func (m *MockLogger) Check(level zapcore.Level, msg string) *zapcore.CheckedEntry {
	return nil
}

func (m *MockLogger) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, string(p))
	return len(p), nil
}

func (m *MockLogger) GetEntries() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.entries...)
}

func (m *MockLogger) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = make([]string, 0)
}
