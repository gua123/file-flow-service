<template>
  <div class="home">
    <!-- 顶部导航栏 -->
    <el-header class="navbar">
      <div class="navbar-brand">文件执行平台</div>
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
      <div class="welcome-section">
        <h1>欢迎使用文件执行平台</h1>
        <p>这是一个安全的文件执行环境，支持多种编程语言的文件执行</p>
      </div>

      <!-- 功能卡片 -->
      <div class="features-grid">
        <el-card class="feature-card" @click="navigateTo('/files')">
          <template #header>
            <div class="card-header">
              <i class="el-icon-upload2"></i>
              <span>文件管理</span>
            </div>
          </template>
          <p>上传、下载和管理执行文件</p>
        </el-card>

        <el-card class="feature-card" @click="navigateTo('/tasks')">
          <template #header>
            <div class="card-header">
              <i class="el-icon-timer"></i>
              <span>任务管理</span>
            </div>
          </template>
          <p>创建、执行和监控任务</p>
        </el-card>

        <el-card class="feature-card" @click="navigateTo('/logs')">
          <template #header>
            <div class="card-header">
              <i class="el-icon-document"></i>
              <span>日志查看</span>
            </div>
          </template>
          <p>查看系统运行日志和执行日志</p>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '../services/authService'
import { ElHeader, ElButton, ElCard } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

// 处理退出登录
const handleLogout = () => {
  authStore.logout()
}

// 导航到指定页面
const navigateTo = (path: string) => {
  router.push(path)
}
</script>

<style scoped>
.home {
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

.welcome-section {
  text-align: center;
  margin-bottom: 40px;
}

.welcome-section h1 {
  color: #333;
  margin-bottom: 10px;
}

.welcome-section p {
  color: #666;
  font-size: 16px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-top: 30px;
}

.feature-card {
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #ebeef5;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: bold;
  color: #333;
}

.card-header i {
  font-size: 20px;
  color: #409eff;
}
</style>
