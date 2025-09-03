# 文件执行平台路由概览

## 后端路由结构

### API路由前缀
所有API接口都使用 `/api/` 前缀

### 主要API端点

| 路径 | 方法 | 描述 |
|------|------|------|
| `/api/health` | GET | 健康检查接口 |
| `/api/download` | GET | 文件下载接口 |
| `/api/tasks` | POST/GET | 任务管理接口 |
| `/api/status` | GET | 获取执行器状态接口 |
| `/api/logs` | GET | 获取日志接口 |
| `/api/config` | GET/POST | 配置管理接口 |
| `/api/login` | POST | 用户登录接口 |

### 配置路由
根据配置文件 `config/config.yaml` 中的设置：
- 健康检查路径: `/health`
- 任务端点: `/task`
- 实际API路径: `/api/task` (通过配置文件设置)

## 前端路由结构

### Vue Router 路由

| 路径 | 组件 | 描述 | 认证要求 |
|------|------|------|----------|
| `/` | 重定向到 `/login` | 根路径重定向 | 无需认证 |
| `/login` | Login.vue | 登录页面 | 无需认证 |
| `/home` | Home.vue | 主页 | 需要认证 |
| `/files` | FileManager.vue | 文件管理页面 | 需要认证 |
| `/tasks` | TaskManager.vue | 任务管理页面 | 需要认证 |
| `/logs` | LogViewer.vue | 日志查看页面 | 需要认证 |

## 路由关系图

```
前端路由访问
    ↓
/login → 登录接口 → 获取token
    ↓
/home, /files, /tasks, /logs → 需要认证的API接口
    ↓
API接口调用
    ↓
后端处理业务逻辑
```

## 认证流程

1. 用户访问 `/login` 路由
2. 用户提交登录信息到 `/api/login` 接口
3. 登录成功后，前端存储JWT token
4. 后续所有需要认证的API请求都需包含:
   ```
   Authorization: Bearer <token>
   ```

## 权限等级
- **Admin**: 所有权限
- **Editor**: 提交任务、上传文件、查看任务状态
- **Viewer**: 查看任务状态、下载结果文件

## 静态资源路由
- 所有静态文件通过 `/` 路径访问
- 前端构建文件位于 `web/file-flow-web/dist/`
- Vue Router 客户端路由处理