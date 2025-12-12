import axios from 'axios'
import { useAuth } from '../store/auth'
import { ElMessage } from 'element-plus'
import router from '../router'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const { authState } = useAuth()
    if (authState.token) {
      config.headers.Authorization = `Bearer ${authState.token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    const { clearAuth } = useAuth()
    const message = error.response?.data?.error || error.message || '请求失败'
    
    // 401 未授权，清除token并跳转到登录页
    if (error.response?.status === 401) {
      clearAuth()
      ElMessage.error('登录已过期，请重新登录')
      router.push('/login')
    }
    
    return Promise.reject(new Error(message))
  }
)

// 生成短信内容
export const generateTemplate = (data) => {
  return api.post('/template/generate', data)
}

// 获取所有位置值
export const getAllPositions = () => {
  return api.get('/positions')
}

// 获取指定位置值
export const getPositionValues = (position) => {
  return api.get(`/positions/${position}`)
}

// 添加位置值
export const addPositionValue = (data) => {
  return api.post('/positions', data)
}

// 设置位置的所有值
export const setPositionValues = (position, values) => {
  return api.put(`/positions/${position}`, { values })
}

// 删除位置值
export const deletePositionValue = (position, value) => {
  return api.delete(`/positions/${position}?value=${encodeURIComponent(value)}`)
}

// 登录
export const login = (data) => {
  return api.post('/auth/login', data)
}

// 注册
export const register = (data) => {
  return api.post('/auth/register', data)
}

// 获取用户信息
export const getUserInfo = () => {
  return api.get('/auth/user')
}

