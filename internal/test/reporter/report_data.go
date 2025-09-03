package reporter

import (
	"time"
)

// TestReport 测试报告结构
type TestReport struct {
	Summary     TestSummary `json:"summary"`
	Tests       []TestResult `json:"tests"`
	Coverage    CoverageInfo `json:"coverage"`
	GeneratedAt string       `json:"generatedAt"`
}

// TestSummary 测试摘要
type TestSummary struct {
	TotalTests    int     `json:"totalTests"`
	PassedTests   int     `json:"passedTests"`
	FailedTests   int     `json:"failedTests"`
	SkippedTests  int     `json:"skippedTests"`
	ExecutionTime string  `json:"executionTime"`
	SuccessRate   float64 `json:"successRate"`
}

// TestResult 单个测试结果
type TestResult struct {
	Name        string `json:"name"`
	Status      string `json:"status"`  // passed, failed, skipped
	Duration    string `json:"duration"`
	File        string `json:"file"`
	Description string `json:"description,omitempty"`
}

// CoverageInfo 代码覆盖率信息
type CoverageInfo struct {
	TotalCoverage float64      `json:"totalCoverage"`
	Files         []FileCoverage `json:"files"`
}

// FileCoverage 文件覆盖率
type FileCoverage struct {
	File        string  `json:"file"`
	Coverage    float64 `json:"coverage"`
	Lines       int     `json:"lines"`
	Covered     int     `json:"covered"`
}

// CreateTestReport 创建测试报告
func CreateTestReport() *TestReport {
	return &TestReport{
		Summary: TestSummary{
			TotalTests:    0,
			PassedTests:   0,
			FailedTests:   0,
			SkippedTests:  0,
			ExecutionTime: "0s",
			SuccessRate:   0.0,
		},
		Tests: make([]TestResult, 0),
		Coverage: CoverageInfo{
			TotalCoverage: 0.0,
			Files: make([]FileCoverage, 0),
		},
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
}
