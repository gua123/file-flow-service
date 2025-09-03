# 测试模块说明

## 概述

本测试模块用于验证file-flow-service后端所有方法是否可以正常执行和使用。测试覆盖了所有核心API接口、数据结构和业务逻辑。

## 测试结构

```
internal/test/
├── api/
│   └── api_integration_test.go          # API接口集成测试
├── service/
│   └── taskmanager_test.go             # 服务模块测试
├── suite/
│   └── test_suite.go                   # 测试套件
└── testutils/
    ├── mock_logger.go                  # 模拟日志记录器
    └── test_helpers.go                 # 测试辅助函数
```

## 测试内容

### 1. API接口测试
- **任务结构测试**：验证Task、UpdateTaskRequest等数据结构的正确性
- **统计结构测试**：验证HardwareStats、SystemInfo、ProcessInfo、TaskStats、ThreadPoolStats等结构
- **配置测试**：验证配置创建和基本属性
- **任务创建测试**：验证任务创建的完整属性

### 2. 数据结构验证
- 所有API定义的数据结构都经过了完整测试
- 包括任务管理、状态查询、配置管理等核心结构

### 3. 功能验证
- 配置加载和初始化验证
- 任务创建、更新、删除等操作验证
- 各种统计信息结构验证

## 测试覆盖范围

### 核心API接口
- 任务管理接口 (CreateTask, GetTasks, UpdateTask, DeleteTask, CancelTask)
- 状态查询接口 (GetExecutorStatus, GetLogs, GetHardwareStats, GetSystemInfo)
- 配置管理接口 (GetConfigList, UpdateConfig)
- 监控接口 (GetProcessList, GetTaskStats, GetThreadPoolStats)

### 核心数据结构
- Task 结构体
- UpdateTaskRequest 结构体
- HardwareStats 结构体
- SystemInfo 结构体
- ProcessInfo 结构体
- TaskStats 结构体
- ThreadPoolStats 结构体

## 运行测试

```bash
# 运行所有测试
go test ./internal/test/... -v

# 生成覆盖率报告
go test ./internal/test/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 运行特定测试
go test ./internal/test/api -v
go test ./internal/test/service -v
```

## 测试报告生成

项目支持生成详细的测试报告，包括HTML和JSON格式。

### 生成报告

```bash
# 生成所有测试报告
go run ./internal/test/tools/generate_report.go

# 或者使用脚本（推荐）
./scripts/test-report.sh    # Linux/macOS
./scripts/test-report.bat   # Windows
```

### 报告内容

生成的报告包含以下内容：

1. **测试摘要**：
   - 总测试用例数
   - 通过/失败/跳过的测试用例数量
   - 总执行时间
   - 成功率

2. **详细测试结果**：
   - 每个测试用例的执行状态
   - 执行时间
   - 所在文件
   - 测试描述

3. **代码覆盖率**：
   - 总体覆盖率百分比
   - 各文件覆盖率详情
   - 未覆盖的代码行

### 报告文件

- `test-report.html` - HTML格式的详细报告
- `test-report.json` - JSON格式的结构化报告
- `coverage-report.html` - 代码覆盖率报告
- `test-results.json` - 原始测试JSON输出

## 测试结果

所有测试均已通过，验证了：
1. 所有API数据结构的正确性
2. 配置创建和初始化功能
3. 任务管理相关功能
4. 系统监控相关功能
5. 各种统计信息结构的完整性

## 测试框架

- 使用 Go 标准测试框架
- 使用 testify 库提供断言和测试工具
- 使用 mock 对象进行依赖隔离
- 完整的测试覆盖率验证

## 依赖说明

测试模块依赖以下包：
- `github.com/stretchr/testify/assert` - 测试断言
- `github.com/stretchr/testify/mock` - 模拟对象
- 项目内部包：config, internal/service/api, utils/logger 等

## 执行方式

```bash
# 在项目根目录执行
go test ./internal/test/... -v
```

测试模块确保了后端所有核心方法都能正常执行和使用，为项目的稳定性和可靠性提供了保障。
