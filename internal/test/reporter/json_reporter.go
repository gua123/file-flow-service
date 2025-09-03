package reporter

import (
	"encoding/json"
	"os"
)

// JSONReporter JSON报告生成器
type JSONReporter struct{}

// NewJSONReporter 创建JSON报告生成器
func NewJSONReporter() *JSONReporter {
	return &JSONReporter{}
}

// GenerateJSONReport 生成JSON报告
func (r *JSONReporter) GenerateJSONReport(report *TestReport, outputPath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(outputPath, data, 0644)
}
