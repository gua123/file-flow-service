<template>
  <div class="file-manager">
    <!-- 顶部导航栏 -->
    <el-header class="navbar">
      <div class="navbar-brand">文件管理</div>
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
            <span>文件上传</span>
          </div>
        </template>
        
        <div class="upload-section">
          <el-upload
            :auto-upload="false"
            :on-change="handleFileChange"
            :limit="1"
            :on-exceed="handleExceed"
            multiple
          >
            <el-button type="primary">选择文件</el-button>
          </el-upload>
          
          <el-button 
            type="success" 
            @click="handleUpload"
            :loading="uploading"
            style="margin-left: 10px;"
          >
            上传文件
          </el-button>
        </div>
      </el-card>

      <el-card class="card">
        <template #header>
          <div class="card-header">
            <span>文件列表</span>
          </div>
        </template>
        
        <el-table :data="fileList" style="width: 100%" v-loading="loading">
          <el-table-column prop="name" label="文件名" width="300" />
          <el-table-column prop="size" label="大小" width="150">
            <template #default="scope">
              {{ formatFileSize(scope.row.size) }}
            </template>
          </el-table-column>
          <el-table-column prop="uploadTime" label="上传时间" width="200" />
          <el-table-column label="操作" width="200">
            <template #default="scope">
              <el-button 
                type="primary" 
                size="small" 
                @click="downloadFile(scope.row.id)"
              >
                下载
              </el-button>
              <el-button 
                type="danger" 
                size="small" 
                @click="deleteFile(scope.row.id)"
              >
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <div class="pagination" v-if="total > 0">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :total="total"
            @current-change="handlePageChange"
            layout="prev, pager, next, jumper"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../services/authService'
import { ElHeader, ElButton, ElCard, ElUpload, ElTable, ElTableColumn, ElPagination, ElMessage } from 'element-plus'
import { FileFlowApi } from '../api/fileFlowApi'

// 定义响应式数据
const router = useRouter()
const authStore = useAuthStore()

const fileList = ref<any[]>([])
const loading = ref(false)
const uploading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const selectedFiles = ref<File[]>([])

// 处理文件选择
const handleFileChange = (file: any, fileList: any[]) => {
  selectedFiles.value = fileList.map(f => f.raw)
}

// 处理文件超出限制
const handleExceed = () => {
  ElMessage.warning('只能选择一个文件')
}

// 处理上传
const handleUpload = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请选择文件')
    return
  }

  try {
    uploading.value = true
    const file = selectedFiles.value[0]
    await FileFlowApi.File.uploadFile(file)
    ElMessage.success('文件上传成功')
    // 重新加载文件列表
    await loadFiles()
  } catch (error) {
    console.error('文件上传失败:', error)
    ElMessage.error('文件上传失败')
  } finally {
    uploading.value = false
  }
}

// 下载文件
const downloadFile = async (fileId: string) => {
  try {
    await FileFlowApi.File.downloadFile(fileId)
  } catch (error) {
    console.error('文件下载失败:', error)
    ElMessage.error('文件下载失败')
  }
}

// 删除文件
const deleteFile = async (fileId: string) => {
  try {
    // 这里应该实现删除文件的逻辑
    ElMessage.success('文件删除成功')
    // 重新加载文件列表
    await loadFiles()
  } catch (error) {
    console.error('文件删除失败:', error)
    ElMessage.error('文件删除失败')
  }
}

// 格式化文件大小
const formatFileSize = (size: number): string => {
  if (size < 1024) {
    return size + ' B'
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(2) + ' KB'
  } else {
    return (size / (1024 * 1024)).toFixed(2) + ' MB'
  }
}

// 加载文件列表
const loadFiles = async () => {
  try {
    loading.value = true
    // 这里应该调用API获取文件列表
    // 暂时使用模拟数据
    fileList.value = [
      {
        id: '1',
        name: 'example.py',
        size: 1024,
        uploadTime: '2025-08-01 10:30:00'
      },
      {
        id: '2',
        name: 'test.js',
        size: 2048,
        uploadTime: '2025-08-02 14:15:00'
      }
    ]
    total.value = 2
  } catch (error) {
    console.error('加载文件列表失败:', error)
    ElMessage.error('加载文件列表失败')
  } finally {
    loading.value = false
  }
}

// 处理分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadFiles()
}

// 处理退出登录
const handleLogout = () => {
  authStore.logout()
}

// 组件挂载时加载文件列表
onMounted(() => {
  loadFiles()
})
</script>

<style scoped>
.file-manager {
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

.upload-section {
  display: flex;
  align-items: center;
  padding: 20px 0;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
