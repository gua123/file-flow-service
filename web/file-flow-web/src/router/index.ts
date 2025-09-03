import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import App from '../App.vue'
import Login from '../views/Login.vue'
import Home from '../views/Home.vue'
import FileManager from '../views/FileManager.vue'
import TaskManager from '../views/TaskManager.vue'
import LogViewer from '../views/LogViewer.vue'
import { useAuthStore } from '../services/authService'
const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/home',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/files',
    name: 'FileManager',
    component: FileManager,
    meta: { requiresAuth: true }
  },
  {
    path: '/tasks',
    name: 'TaskManager',
    component: TaskManager,
    meta: { requiresAuth: true }
  },
  {
    path: '/logs',
    name: 'LogViewer',
    component: LogViewer,
    meta: { requiresAuth: true }
  }
]
const router = createRouter({
  history: createWebHistory('/fileflow/'),
  routes
})

// 全局路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // 如果访问需要认证的页面但未登录，则跳转到登录页
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
    return
  }
  
  // 如果访问登录页但已登录，则跳转到主页
  if (to.path === '/login' && authStore.isLoggedIn) {
    next('/home')
    return
  }
  
  next()
})

export default router