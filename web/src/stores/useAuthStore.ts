import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/user'

// sessionStorage 键名常量
const AUTH_STORAGE_KEY = 'auth_state'

// 认证状态接口
interface AuthState {
  user: User | null
}

/**
 * 从 sessionStorage 读取认证状态
 */
const loadFromSession = (): AuthState | null => {
  try {
    const stored = sessionStorage.getItem(AUTH_STORAGE_KEY)
    if (!stored) return null
    return JSON.parse(stored) as AuthState
  } catch {
    return null
  }
}

/**
 * 保存认证状态到 sessionStorage
 */
const saveToSession = (state: AuthState): void => {
  try {
    sessionStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(state))
  } catch (error) {
    console.error('保存认证状态失败:', error)
  }
}

/**
 * 清除 sessionStorage 中的认证状态
 */
const clearSession = (): void => {
  sessionStorage.removeItem(AUTH_STORAGE_KEY)
}

export const useAuthStore = defineStore('auth', () => {
  // ========== 状态定义 ==========
  const user = ref<User | null>(null)

  // ========== 计算属性 ==========
  
  /**
   * 判断用户是否已认证（有用户信息）
   */
  const isAuthenticated = computed(() => {
    return !!user.value
  })

  // ========== 操作方法 ==========

  /**
   * 登录成功后调用，保存用户信息
   * @param userData 用户信息
   */
  const loginSuccess = (userData: User) => {
    // 更新状态
    user.value = userData
    
    // 持久化到 sessionStorage
    saveToSession({
      user: user.value
    })
  }

  /**
   * 清除认证状态（用于登出或刷新失败）
   */
  const clearAuth = () => {
    logout()
  }

  /**
   * 从 sessionStorage 初始化认证状态
   * 应用启动时调用
   */
  const initFromSession = () => {
    const stored = loadFromSession()
    if (!stored) return
    
    // 恢复状态
    user.value = stored.user
  }

  /**
   * 登出，清空所有认证状态
   */
  const logout = () => {
    user.value = null
    clearSession()
  }

  return {
    // 状态
    user,
    
    // 计算属性
    isAuthenticated,
    
    // 操作方法
    loginSuccess,
    clearAuth,
    initFromSession,
    logout
  }
})
