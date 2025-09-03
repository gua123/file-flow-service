<template>
  <div class="task-manager">
    <!-- 顶部导航栏 -->
    <el-header class="navbar">
      <div class="navbar-brand">任务管理</div>
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
            <span>创建新任务</span>
          </div>
        </template>
        
        <el-form :model="taskForm" :rules="rules" ref="taskFormRef" label-width="100px">
          <el-form-item label="任务名称" prop="name">
            <el-input v-model="taskForm.name" placeholder="请输入任务名称" />
          </el-form-item>
          
          <el-form-item label="描述" prop="description">
            <el-input 
              v-model="taskForm.description" 
              type="textarea" 
              placeholder="请输入任务描述"
              :rows="3"
            />
          </el-form-item>
          
          <el-form-item label="执行命令" prop="cmd">
            <el-input v-model="taskForm.cmd" placeholder="请输入执行命令" />
          </el-form-item>
          
          <el-form-item label="参数" prop="args">
            <el-input 
              v-model="taskForm.args" 
              placeholder="请输入参数，多个参数用逗号分隔"
            />
          </el-form-item>
          
          <el-form-item label="工作目录" prop="dir">
            <el-input v-model="taskForm.dir" placeholder="请输入工作目录" />
          </el-form-item>
          
          <el-form-item label="创建者" prop="creator">
            <el-input v-model="taskForm.creator" placeholder="请输入创建者" />
          </el-form-item>
          
          <el-form-item label="分配给" prop="assignedTo">
            <el-input v-model="taskForm.assignedTo" placeholder="请输入分配给" />
          </el-form-item>
          
          <el-form-item label="结果路径" prop="resultPath">
            <el-input v-model="taskForm.resultPath" placeholder="请输入结果路径" />
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="createTask">创建任务</el-button>
            <el-button @click="resetForm">重置</el-button>
          </el-form-item>
        </el-form>
      </el-card>

      <el-card class="card">
        <template #header>
          <div class="card-header">
            <span>任务列表</span>
          </div>
        </template>
        
        <el-table :data="taskList" style="width: 100%" v-loading="loading">
          <el-table-column prop="name" label="任务名称" width="150" />
          <el-table-column prop="cmd" label="执行命令" width="200" />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <el-tag :type="getStatusType(scope.row.status)">
                {{ scope.row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="progress" label="进度" width="120">
            <template #default="scope">
              <el-progress :percentage="scope.row.progress" />
            </template>
          </el-table-column>
          <el-table-column prop="creator" label="创建者" width="120" />
          <el-table-column prop="createTime" label="创建时间" width="180" />
          <el-table-column label="操作" width="250">
            <template #default="scope">
              <el-button 
                type="primary" 
                size="small" 
                @click="startTask(scope.row.id)"
                :disabled="scope.row.status === 'running' || scope.row.status === 'completed'"
              >
                开始
              </el-button>
              <el-button 
                type="warning" 
                size="small" 
                @click="cancelTask(scope.row.id)"
                :disabled="scope.row.status !== 'running'"
              >
                取消
              </el-button>
              <el-button 
                type="danger" 
                size="small" 
                @click="deleteTask(scope.row.id)"
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
import { 
  ElHeader, ElButton, ElCard, ElForm, ElFormItem, ElInput, ElTable, 
  ElTableColumn, ElPagination, ElTag, ElProgress, ElMessage 
} from 'element-plus'
import { FileFlowApi } from '../api/fileFlowApi'

// 定义任务表单类型
interface TaskForm {
  name: string
  description: string
  cmd: string
  args: string
  dir: string
  creator: string
  assignedTo: string
  resultPath: string
}

// 定义任务数据类型
interface Task {
  id: string
  name: string
  description: string
  cmd: string
  args: string[]
  dir: string
  creator: string
  assignedTo: string
  resultPath: string
  status: string
  progress: number
  createTime: string
}

// 定义响应式数据
const router = useRouter()
const authStore = useAuthStore()

const taskFormRef = ref<any>(null)
const taskForm = ref<TaskForm>({
  name: '',
  description: '',
  cmd: '',
  args: '',
  dir: '',
  creator: '',
  assignedTo: '',
  resultPath: ''
})

const taskList = ref<Task[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' }
  ],
  cmd: [
    { required: true, message: '请输入执行命令', trigger: 'blur' }
  ],
  creator: [
    { required: true, message: '请输入创建者', trigger: 'blur' }
  ]
}

// 获取状态类型
const getStatusType = (status: string): string => {
  switch (status) {
    case 'pending':
      return 'info'
    case 'running':
      return 'primary'
    case 'completed':
      return 'success'
    case 'error':
      return 'danger'
    default:
      return 'info'
  }
}

// 创建任务
const createTask = async () => {
  if (!taskFormRef.value) return
  
  try {
    await taskFormRef.value.validate()
    
    // 处理参数
    const args = taskForm.value.args ? taskForm.value.args.split(',').map(arg => arg.trim()) : []
    
    const taskData = {
      name: taskForm.value.name,
      description: taskForm.value.description,
      cmd: taskForm.value.cmd,
      args: args,
      dir: taskForm.value.dir,
      creator: taskForm.value.creator,
      assignedTo: taskForm.value.assignedTo,
      resultPath: taskForm.value.resultPath,
      status: 'pending',
      progress: 0
    }
    
    await FileFlowApi.Task.createTask(taskData)
    ElMessage.success('任务创建成功')
    resetForm()
    await loadTasks()
  } catch (error) {
    console.error('创建任务失败:', error)
    ElMessage.error('任务创建失败')
  }
}

// 重置表单
const resetForm = () => {
  taskForm.value = {
    name: '',
    description: '',
    cmd: '',
    args: '',
    dir: '',
    creator: '',
    assignedTo: '',
    resultPath: ''
  }
  taskFormRef.value?.resetFields()
}

// 开始任务
const startTask = async (taskId: string) => {
  try {
    // 这里应该实现开始任务的逻辑
    ElMessage.success('任务开始执行')
    await loadTasks()
  } catch (error) {
    console.error('开始任务失败:', error)
    ElMessage.error('任务开始失败')
  }
}

// 取消任务
const cancelTask = async (taskId: string) => {
  try {
    await FileFlowApi.Task.cancelTask(taskId)
    ElMessage.success('任务已取消')
    await loadTasks()
  } catch (error) {
    console.error('取消任务失败:', error)
    ElMessage.error('任务取消失败')
  }
}

// 删除任务
const deleteTask = async (taskId: string) => {
  try {
    await FileFlowApi.Task.deleteTask(taskId)
    ElMessage.success('任务已删除')
    await loadTasks()
  } catch (error) {
    console.error('删除任务失败:', error)
    ElMessage.error('任务删除失败')
  }
}

// 加载任务列表
const loadTasks = async () => {
  try {
    loading.value = true
    // 这里应该调用API获取任务列表
    // 暂时使用模拟数据
    taskList.value = [
      {
        id: '1',
        name: '测试任务1',
        description: '测试任务描述1',
        cmd: 'python test.py',
        args: ['arg1', 'arg2'],
        dir: '/home/user',
        creator: 'admin',
        assignedTo: 'editor',
        resultPath: '/home/user/result',
        status: 'pending',
        progress: 0,
        createTime: '2025-08-01 10:30:00'
      },
      {
        id: '2',
        name: '测试任务2',
        description: '测试任务描述2',
        cmd: 'node test.js',
        args: [],
        dir: '/home/user',
        creator: 'admin',
        assignedTo: 'viewer',
        resultPath: '/home/user/result',
        status: 'running',
        progress: 50,
        createTime: '2025-08-02 14:15:00'
      }
    ]
    total.value = 2
  } catch (error) {
    console.error('加载任务列表失败:', error)
    ElMessage.error('加载任务列表失败')
  } finally {
    loading.value = false
  }
}

// 处理分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadTasks()
}

// 处理退出登录
const handleLogout = () => {
  authStore.logout()
}

// 组件挂载时加载任务列表
onMounted(() => {
  loadTasks()
})
</script>

<style scoped>
.task-manager {
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

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
