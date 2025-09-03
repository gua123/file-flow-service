<template>
  <div class="log-viewer">
    <!-- 顶部导航栏 -->
    <el-header class="navbar">
      <div class="navbar-brand">日志查看</div>
      <div class="navbar-menu">
        <el-button 
          type="primary" 
          @click="handleLogout"
          size="small"
        >
          退出登录
        </el-button>
      </div>
    </el-header>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <el-card class="card">
        <template #header>
          <div class="card-header">
            <span>日志筛选</span>
          </div>
        </template>
        
        <el-form :model="logFilter" label-width="80px" inline>
          <el-form-item label="日志类型">
            <el-select v-model="logFilter.logType" placeholder="请选择日志类型">
              <el-option label="平台日志" value="service" />
              <el-option label="执行日志" value="flow" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="时间范围">
            <el-date-picker
              v-model="logFilter.timeRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="yyyy-MM-dd"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="loadLogs">查询</el-button>
            <el-button @click="resetFilter">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <el-card class="card">
        <template #header>
          <div class="card-header">
            <span>日志内容</span>
          </div>
        </template>
        
        <div class="log-content" v-loading="loading">
          <el-scrollbar height="500">
            <pre v-if="logContent">{{ logContent }}</pre>
            <div v-else class="empty-placeholder">
              请选择日志类型和时间范围进行查询
            </div>
          </el-scrollbar>
        </div>
        
        <div class="log-actions">
          <el-button 
            type="primary" 
            @click="downloadLogs"
            :disabled="!logContent"
          >
            下载日志
          </el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../services/authService'
import { 
  ElHeader, ElButton, ElCard, ElForm, ElFormItem, ElSelect, ElOption, 
  ElDatePicker, ElScrollbar, ElMessage 
} from 'element-plus'
import { FileFlowApi } from '../api/fileFlowApi'

// 定义日志筛选条件类型
interface LogFilter {
  logType: string
  timeRange: [string, string] | null
}

// 定义响应式数据
const router = useRouter()
const authStore = useAuthStore()

const logFilter = ref<LogFilter>({
  logType: 'service',
  timeRange: null
})

const logContent = ref<string>('')
const loading = ref(false)

// 加载日志
const loadLogs = async () => {
  try {
    loading.value = true
    
    // 获取当前时间范围
    let since = ''
    if (logFilter.value.timeRange && logFilter.value.timeRange[0] && logFilter.value.timeRange[1]) {
      since = `${logFilter.value.timeRange[0]} to ${logFilter.value.timeRange[1]}`
    }
    
    // 这里应该调用API获取日志
    // 暂时使用模拟数据
    logContent.value = `日志类型: ${logFilter.value.logType}
时间范围: ${since || '全部时间'}

2025-08-01 10:30:00 [INFO] 服务启动成功
2025-08-01 10:30:01 [DEBUG] 初始化日志模块
2025-08-01 10:30:02 [INFO] 配置模块初始化完成
2025-08-01 10:30:03 [INFO] 文件模块初始化完成
2025-08-01 10:30:04 [INFO] 沙箱模块初始化完成
2025-08-01 10:30:05 [INFO] Web模块初始化完成
2025-08-01 10:30:06 [INFO] 服务监听端口: 8080

2025-08-01 14:15:00 [INFO] 接收到文件上传请求
2025-08-01 14:15:01 [DEBUG] 文件大小: 1024 bytes
2025-08-01 14:15:02 [INFO] 文件上传成功
2025-08-01 14:15:03 [INFO] 创建执行任务: test_task
2025-08-01 14:15:04 [INFO] 任务开始执行
2025-08-01 14:15:05 [INFO] 任务执行完成
2025-08-01 14:15:06 [INFO] 生成执行结果文件
`
    
  } catch (error) {
    console.error('加载日志失败:', error)
    ElMessage.error('加载日志失败')
  } finally {
    loading.value = false
  }
}

// 重置筛选条件
const resetFilter = () => {
  logFilter.value = {
    logType: 'service',
    timeRange: null
  }
  logContent.value = ''
}

// 下载日志
const downloadLogs = () => {
  if (!logContent.value) return
  
  try {
    const blob = new Blob([logContent.value], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `log_${new Date().toISOString().slice(0, 10)}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    ElMessage.success('日志下载成功')
  } catch (error) {
    console.error('下载日志失败:', error)
    ElMessage.error('下载日志失败')
  }
}

// 处理退出登录
const handleLogout = () => {
  authStore.logout()
}

// 组件挂载时初始化
onMounted(() => {
  // 可以在这里添加初始化逻辑
})
</script>

<style scoped>
.log-viewer {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  height: 60px;
}

.navbar-brand {
  font-size: 1.5rem;
  font-weight: bold;
  color: #409eff;
}

.navbar-menu {
  display: flex;
  align-items: center;
  gap: 15px;
}

.main-content {
  padding: 20px;
  margin: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.card {
  margin-bottom: 20px;
}

.card-header {
  font-weight: bold;
  color: #333;
}

.log-content {
  margin: 20px 0;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.log-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.empty-placeholder {
  text-align: center;
  color: #999;
  padding: 40px;
}

.log-actions {
  text-align: right;
  padding: 10px 0;
}
</style>
