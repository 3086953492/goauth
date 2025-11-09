import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types/user'
import type { TokenResponse } from '@/types/auth'

// 从 localStorage 恢复用户信息
const getUserFromStorage = (): User | null => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      return JSON.parse(userStr) as User
    } catch {
      return null
    }
  }
  return null
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const refreshToken = ref<string>(localStorage.getItem('refreshToken') || '')
  const expiresIn = ref<number>(Number(localStorage.getItem('expiresIn') || 0))
  const user = ref<User | null>(getUserFromStorage())
  const isLoggedIn = ref<boolean>(!!token.value)

  /**
   * 设置令牌信息
   */
  const setToken = (tokenData: TokenResponse) => {
    token.value = tokenData.access_token
    refreshToken.value = tokenData.refresh_token
    expiresIn.value = tokenData.expires_in
    
    localStorage.setItem('token', tokenData.access_token)
    localStorage.setItem('refreshToken', tokenData.refresh_token)
    localStorage.setItem('expiresIn', tokenData.expires_in.toString())
    
    isLoggedIn.value = true
  }

  /**
   * 设置用户信息
   */
  const setUser = (userData: User) => {
    user.value = userData
    // 持久化用户信息到 localStorage
    localStorage.setItem('user', JSON.stringify(userData))
  }

  /**
   * 登出
   */
  const logout = () => {
    token.value = ''
    refreshToken.value = ''
    expiresIn.value = 0
    user.value = null
    isLoggedIn.value = false
    
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('expiresIn')
    localStorage.removeItem('user')
  }

  /**
   * 清除所有状态
   */
  const reset = () => {
    logout()
  }

  return {
    token,
    refreshToken,
    expiresIn,
    user,
    isLoggedIn,
    setToken,
    setUser,
    logout,
    reset
  }
})

