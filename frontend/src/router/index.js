import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '../store/auth'
import { ElMessage } from 'element-plus'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'TemplateEditor',
    component: () => import('../views/TemplateEditor.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/config',
    name: 'PositionConfig',
    component: () => import('../views/PositionConfig.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const { authState } = useAuth()
  
  if (to.meta.requiresAuth && !authState.isAuthenticated) {
    ElMessage.warning('请先登录')
    next('/login')
  } else if (to.path === '/login' && authState.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router

