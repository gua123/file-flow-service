/**
 * 文件执行平台API接口
 * 集中管理所有前后端交互逻辑
 */

import { ElMessage } from 'element-plus'

// API基础URL（可以根据需要修改）
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

/**
 * API请求基础类
 */
class ApiClient {
  /**
   * 发送API请求
   */
  static async request<T>(url: string, options: RequestInit = {}): Promise<T> {
    const token = localStorage.getItem('token')
    
    // 设置默认请求头
    const defaultHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    
    // 如果有token，添加到请求头
    if (token) {
      defaultHeaders['Authorization'] = `Bearer ${token}`
    }
    
    // 合并请求头
    const config = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    }
    
    try {
      const response = await fetch(`${API_BASE_URL}${url}`, config)
      
      // 检查响应状态
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const data = await response.json()
      
      // 检查API响应状态
      if (data.code !== 0) {
        throw new Error(data.message || '请求失败')
      }
      
      return data.data
    } catch (error) {
      console.error('API请求失败:', error)
      ElMessage.error(error instanceof Error ? error.message : '请求失败')
      throw error
    }
  }
}

/**
 * 文件相关API
 */
export const FileApi = {
  /**
   * 上传文件
   * @param file 文件对象
   * @returns 上传结果
   *
   * 示例请求:
   * POST /api/task
   * FormData: { file: [File] }
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": {
   *     "fileId": "file_123",
   *     "fileName": "example.zip"
   *   }
   * }
   */
  async uploadFile(file: File): Promise<any> {
    const formData = new FormData()
    formData.append('file', file)
    
    try {
      const response = await fetch(`${API_BASE_URL}/api/task`, {
        method: 'POST',
        body: formData,
      })
      
      const data = await response.json()
      
      if (data.code !== 0) {
        throw new Error(data.message || '上传失败')
      }
      
      return data.data
    } catch (error) {
      console.error('文件上传失败:', error)
      ElMessage.error(error instanceof Error ? error.message : '上传失败')
      throw error
    }
  },
  
  /**
   * 下载文件
   * @param fileId 文件ID
   *
   * 示例请求:
   * GET /api/download?file_id=file_123
   *
   * 示例响应:
   * 二进制文件流 (application/octet-stream)
   */
  async downloadFile(fileId: string): Promise<void> {
    window.open(`${API_BASE_URL}/api/download?file_id=${fileId}`, '_blank')
  },
  
  /**
   * 获取文件列表
   * @returns 文件列表
   *
   * 示例请求:
   * GET /api/files
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": [
   *     {
   *       "fileId": "file_123",
   *       "fileName": "example.zip",
   *       "size": 1024,
   *       "uploadTime": "2025-08-18T13:00:00Z"
   *     }
   *   ]
   * }
   */
  async getFiles(): Promise<any[]> {
    return ApiClient.request<any[]>('/api/files')
  }
}

/**
 * 任务相关API
 */
export const TaskApi = {
  /**
   * 创建任务
   * @param taskData 任务数据
   * @returns 任务ID
   *
   * 示例请求:
   * POST /api/tasks
   * {
   *   "name": "文件处理任务",
   *   "type": "file_processing",
   *   "fileId": "file_123"
   * }
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": "task_456"
   * }
   */
  async createTask(taskData: any): Promise<string> {
    return ApiClient.request<string>('/api/tasks', {
      method: 'POST',
      body: JSON.stringify({ name: taskData.name, commandParams: { type: taskData.type, fileId: taskData.fileId } }),
    })
  },
  
  /**
   * 获取任务列表
   * @param page 页码
   * @param pageSize 每页大小
   * @returns 任务列表
   *
   * 示例请求:
   * GET /api/tasks
   * {
   *   "page": 1,
   *   "pageSize": 10
   * }
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": [
   *     {
   *       "taskId": "task_456",
   *       "name": "文件处理任务",
   *       "status": "processing",
   *       "createTime": "2025-08-18T13:00:00Z"
   *     }
   *   ]
   * }
   */
  async getTasks(page: number = 1, pageSize: number = 10): Promise<any[]> {
    return ApiClient.request<any[]>(`/api/tasks?page=${page}&pageSize=${pageSize}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
  },
  
  /**
   * 更新任务
   * @param taskId 任务ID
   * @param taskData 任务数据
   *
   * 示例请求:
   * POST /api/tasks
   * {
   *   "id": "task_456",
   *   "status": "completed"
   * }
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": null
   * }
   */
  async updateTask(taskId: string, taskData: any): Promise<void> {
    return ApiClient.request<void>('/api/tasks', {
      method: 'POST',
      body: JSON.stringify({ id: taskId, ...taskData }),
    })
  },
  
  /**
   * 删除任务
   * @param taskId 任务ID
   *
   * 示例请求:
   * POST /api/tasks
   * {
   *   "id": "task_456"
   * }
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": null
   * }
   */
  async deleteTask(taskId: string): Promise<void> {
    return ApiClient.request<void>('/api/tasks', {
      method: 'POST',
      body: JSON.stringify({ id: taskId }),
    })
  },
  
  /**
   * 取消任务
   * @param taskId 任务ID
   *
   * 示例请求:
   * POST /api/tasks/cancel
   * {
   *   "id": "task_456"
   * }
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": null
   * }
   */
  async cancelTask(taskId: string): Promise<void> {
    return ApiClient.request<void>('/api/tasks/cancel', {
      method: 'POST',
      body: JSON.stringify({ id: taskId }),
    })
  },
  
  /**
   * 获取执行器状态
   * @returns 执行器状态
   *
   * 示例请求:
   * GET /api/status
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": {
   *     "status": "online",
   *     "activeTasks": 2,
   *     "memoryUsage": "45%"
   *   }
   * }
   */
  async getExecutorStatus(): Promise<any> {
    return ApiClient.request<any>('/api/status')
  }
}

/**
 * 日志相关API
 */
export const LogApi = {
  /**
   * 获取日志
   * @param logType 日志类型
   * @param since 时间范围
   * @returns 日志列表
   *
   * 示例请求:
   * GET /api/logs
   * {
   *   "logType": "file",
   *   "since": "2025-08-18T12:00:00Z"
   * }
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": [
   *     {
   *       "logId": "log_789",
   *       "type": "file",
   *       "message": "文件上传成功",
   *       "timestamp": "2025-08-18T13:00:00Z"
   *     }
   *   ]
   * }
   */
  async getLogs(logType: string, since: string): Promise<any[]> {
    return ApiClient.request<any[]>(`/api/logs?logType=${encodeURIComponent(logType)}&since=${encodeURIComponent(since)}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }
}

/**
 * 配置相关API
 */
export const ConfigApi = {
  /**
   * 获取配置列表
   * @returns 配置列表
   *
   * 示例请求:
   * GET /api/config
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": [
   *     {
   *       "key": "maxFileSize",
   *       "value": "100MB",
   *       "description": "最大文件大小限制"
   *     }
   *   ]
   * }
   */
  async getConfigList(): Promise<any[]> {
    return ApiClient.request<any[]>('/api/config')
  },
  
  /**
   * 更新配置
   * @param key 配置键
   * @param value 配置值
   *
   * 示例请求:
   * POST /api/config
   * {
   *   "key": "maxFileSize",
   *   "value": "200MB"
   * }
   *
   * 示例响应:
   * {
   *   "code": 0,
   *   "message": "success",
   *   "data": null
   * }
   */
  async updateConfig(key: string, value: string): Promise<void> {
    return ApiClient.request<void>('/api/config', {
      method: 'POST',
      body: JSON.stringify({ key, value }),
    })
  }
}

/**
 * 认证相关API
 */
export const AuthApi = {
  /**
   * 登录
   * @param username 用户名
   * @param password 密码
   * @returns 登录结果
   *
   * 示例请求:
   * POST /api/login
   * Content-Type: application/x-www-form-urlencoded
   * username=admin&password=123456
   *
   * 示例响应:
   * {
   *   "status": "success",
   *   "message": "登录成功",
   *   "data": {
   *     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
   *     "user": {
   *       "id": 1,
   *       "name": "admin"
   *     }
   *   }
   * }
   */
  async login(username: string, password: string): Promise<any> {
    try {
      const response = await fetch(`${API_BASE_URL}/api/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      })
      
      const data = await response.json();
      console.log('登录API响应数据:', data); // 调试日志
      
      // 添加调试信息：检查返回的字段
      console.log('登录响应字段检查:', {
        status: data.status,
        code: data.code,
        message: data.message,
        data: data.data
      });
      
      localStorage.setItem('token', data.data?.token || data.token);
      
      // 修正：根据实际后端返回格式判断登录状态
      // 后端可能返回的是code字段而不是status字段
      const isSuccess = (data.status === 'success') || (data.code === 0);
      console.log('登录成功判断结果:', isSuccess, 'status:', data.status, 'code:', data.code);
      
      if (!isSuccess) {
        throw new Error(data.message || '登录失败')
      }
      
      return data
    } catch (error) {
      console.error('登录失败:', error)
      ElMessage.error(error instanceof Error ? error.message : '登录失败')
      throw error
    }
  }
}

/**
 * 统一的API接口对象
 */
export const FileFlowApi = {
  File: FileApi,
  Task: TaskApi,
  Log: LogApi,
  Config: ConfigApi,
  Auth: AuthApi
}