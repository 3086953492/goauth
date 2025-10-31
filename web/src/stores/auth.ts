import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)
  const isLoggedIn = ref<boolean>(!!token.value)

  /**
   * 设置令牌
   */
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
    isLoggedIn.value = true
  }

  /**
   * 设置用户信息
   */
  const setUser = (userData: User) => {
    user.value = userData
  }

  /**
   * 登出
   */
  const logout = () => {
    token.value = ''
    user.value = null
    isLoggedIn.value = false
    localStorage.removeItem('token')
  }

  /**
   * 清除所有状态
   */
  const reset = () => {
    logout()
  }

  return {
    token,
    user,
    isLoggedIn,
    setToken,
    setUser,
    logout,
    reset
  }
})

