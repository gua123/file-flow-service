# File Flow Web

文件执行平台前端项目，基于Vue 3 + Element Plus构建。

## 项目结构

```
src/
├── api/                    # API接口定义
│   └── fileFlowApi.ts      # 集中管理所有API请求
├── assets/                 # 静态资源
├── components/             # 公共组件
│   └── PermissionWrapper.vue # 权限控制组件
├── services/               # 服务层
│   └── authService.ts      # 认证和权限服务
├── views/                  # 页面视图
│   ├── Home.vue            # 主页
│   ├── Login.vue           # 登录页
│   ├── FileManager.vue     # 文件管理页
│   ├── TaskManager.vue     # 任务管理页
│   └── LogViewer.vue       # 日志查看页
├── router/                 # 路由配置
│   └── index.ts            # 路由定义
├── store/                  # 状态管理
├── styles/                 # 样式文件
│   └── main.css            # 全局样式
├── App.vue                 # 根组件
└── main.ts                 # 入口文件
```

## 功能特性

- **用户认证**：登录/退出功能
- **权限控制**：基于角色的权限管理
- **文件管理**：文件上传、下载、删除
- **任务管理**：任务创建、执行、取消、删除
- **日志查看**：平台日志和执行日志查看
- **响应式设计**：适配不同屏幕尺寸

## API接口

所有API请求都通过 `src/api/fileFlowApi.ts` 文件集中管理，包括：

- 文件相关接口
- 任务相关接口  
- 日志相关接口
- 配置相关接口
- 认证相关接口

## 开发命令

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

## 权限控制

项目实现了基于角色的权限控制：

- **Admin（管理员）**：所有功能权限
- **Editor（编辑者）**：可创建、编辑任务，查看日志
- **Viewer（查看者）**：仅可查看任务和日志

权限控制通过 `PermissionWrapper` 组件实现，可以轻松地为页面元素添加权限控制。
