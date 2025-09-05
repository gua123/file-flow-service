package reporter

import (
	"html/template"
	"os"
)

// HTMLReporter HTML报告生成器
type HTMLReporter struct {
	templatePath string
}

// NewHTMLReporter 创建HTML报告生成器
func NewHTMLReporter() *HTMLReporter {
	return &HTMLReporter{
		templatePath: "internal/test/reporter/report_template.html",
	}
}

// GenerateHTMLReport 生成HTML报告
func (r *HTMLReporter) GenerateHTMLReport(report *TestReport, outputPath string) error {
	// 创建HTML模板
	tmpl := template.Must(template.New("report").Parse(htmlTemplate))
	
	// 创建输出文件
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	// 执行模板并写入文件
	return tmpl.Execute(file, report)
}

// HTML模板内容
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>测试报告</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background-color: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .summary { background: #f8f9fa; padding: 20px; border-radius: 5px; margin-bottom: 20px; border-left: 4px solid #007bff; }
        .test-result { margin: 10px 0; padding: 15px; border: 1px solid #ddd; border-radius: 5px; }
        .passed { background: #d4edda; border-color: #c3e6cb; }
        .failed { background: #f8d7da; border-color: #f5c6cb; }
        .skipped { background: #fff3cd; border-color: #ffeaa7; }
        table { border-collapse: collapse; width: 100%; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background-color: #f2f2f2; font-weight: bold; }
        .success-rate { font-size: 1.2em; font-weight: bold; color: #28a745; }
        .failure-rate { font-size: 1.2em; font-weight: bold; color: #dc3545; }
        .header { text-align: center; margin-bottom: 30px; }
        .header h1 { color: #333; }
        .header p { color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>测试报告</h1>
            <p>生成时间: {{.GeneratedAt}}</p>
        </div>
        
        <div class="summary">
            <h2>测试摘要</h2>
            <p><strong>总测试数:</strong> {{.Summary.TotalTests}}</p>
            <p><strong>通过:</strong> {{.Summary.PassedTests}}</p>
            <p><strong>失败:</strong> {{.Summary.FailedTests}}</p>
            <p><strong>跳过:</strong> {{.Summary.SkippedTests}}</p>
            <p><strong>执行时间:</strong> {{.Summary.ExecutionTime}}</p>
            <p class="success-rate">成功率: {{printf "%.2f" .Summary.SuccessRate}}%</p>
        </div>
        
        <h2>详细测试结果</h2>
        {{range .Tests}}
        <div class="test-result {{.Status}}">
            <h3>{{.Name}}</h3>
            <p><strong>状态:</strong> {{.Status}}</p>
            <p><strong>执行时间:</strong> {{.Duration}}</p>
            <p><strong>文件:</strong> {{.File}}</p>
            {{if .Description}}
            <p><strong>描述:</strong> {{.Description}}</p>
            {{end}}
        </div>
        {{end}}
        
        <h2>代码覆盖率</h2>
        <p><strong>总体覆盖率:</strong> {{printf "%.2f" .Coverage.TotalCoverage}}%</p>
        <table>
            <tr>
                <th>文件</th>
                <th>覆盖率</th>
                <th>总行数</th>
                <th>覆盖行数</th>
            </tr>
            {{range .Coverage.Files}}
            <tr>
                <td>{{.File}}</td>
                <td>{{printf "%.2f" .Coverage}}%</td>
                <td>{{.Lines}}</td>
                <td>{{.Covered}}</td>
            </tr>
            {{end}}
        </table>
    </div>
</body>
</html>
`
