# API 路由列表

## 任务管理模块

| 接口名 | 调用路由路径 | 示例参数 | 返回结果 |
|--------|--------------|----------|----------|
| 获取任务列表 | GET /api/tasks | page=1&size=10 | { "items": [{"id": "task_123", "name": "test_task"}], "total": 1 } |
| 更新任务 | POST /api/tasks/{id} | { "status": "completed" } | { "id": "task_123", "status": "completed" } |
| 删除任务 | POST /api/tasks/{id} | - | { "message": "success" } |

## 认证模块

| 接口名 | 调用路由路径 | 示例参数 | 返回结果 |
|--------|--------------|----------|----------|

## 执行器模块

| 接口名 | 调用路由路径 | 示例参数 | 返回结果 |
|--------|--------------|----------|----------|

## 日志模块

| 接口名 | 调用路由路径 | 示例参数 | 返回结果 |
|--------|--------------|----------|----------|
| 获取日志 | GET /api/logs | logType=execution&since=2023-01-01 | [ { "id": "log_1", "message": "Task started" } ] |

## 配置模块

| 接口名 | 调用路由路径 | 示例参数 | 返回结果 |
|--------|--------------|----------|----------|
| 更新配置 | POST /api/config/{key} | { "value": "20" } | { "key": "max_threads", "value": "20" } |

## 文件管理模块

| 接口名 | 调用路由路径 | 示例参数 | 返回结果 |
|--------|--------------|----------|----------|
| 文件下载 | GET /api/download/{id} | - | (binary file) |

## 系统监控模块

| 接口名 | 调用路由路径 | 示例参数 | 返回结果 |
|--------|--------------|----------|----------|
| 硬件统计 | GET /api/hardware/stats | - | { "cpu_usage": 45, "memory_usage": 60 } |
| 系统信息 | GET /api/system/info | - | { "os": "Windows 10", "version": "1.0.0" } |
| 进程列表 | GET /api/processes | - | [ { "pid": 123, "name": "file-flow-service" } ] |
| 任务统计 | GET /api/tasks/stats | - | { "total_tasks": 100, "completed": 80 } |
| 线程池统计 | GET /api/threadpool/stats | - | { "active_threads": 5, "queue_size": 2 } |