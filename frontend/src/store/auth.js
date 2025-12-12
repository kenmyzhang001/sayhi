// 简单的状态管理（可以使用 Pinia 或 Vuex，这里用简单的响应式对象）
import { reactive } from 'vue'

const authState = reactive({
  token: localStorage.getItem('token') || '',
  username: localStorage.getItem('username') || '',
  isAuthenticated: !!localStorage.getItem('token')
})

export const useAuth = () => {
  const setAuth = (token, username) => {
    authState.token = token
    authState.username = username
    authState.isAuthenticated = true
    localStorage.setItem('token', token)
    localStorage.setItem('username', username)
  }

  const clearAuth = () => {
    authState.token = ''
    authState.username = ''
    authState.isAuthenticated = false
    localStorage.removeItem('token')
    localStorage.removeItem('username')
  }

  return {
    authState,
    setAuth,
    clearAuth
  }
}

