<template>
  <div>
    <!-- 如果用户有权限，则显示内容 -->
    <slot v-if="hasPermission" />
    
    <!-- 如果用户没有权限，则显示无权限提示 -->
    <div v-else class="permission-denied">
      <el-alert
        title="权限不足"
        type="error"
        description="您没有访问此功能的权限"
        show-icon
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '../services/authService'
import { ElAlert } from 'element-plus'

// 定义props
const props = defineProps<{
  permission: string
}>()

// 获取认证store
const authStore = useAuthStore()

// 计算用户是否有权限
const hasPermission = computed(() => {
  return authStore.hasPermission(props.permission)
})
</script>

<style scoped>
.permission-denied {
  padding: 20px;
  text-align: center;
}
</style>
