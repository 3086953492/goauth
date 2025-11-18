import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types/user'
import type { TokenResponse } from '@/types/auth'

// sessionStorage 键名常量
const AUTH_STORAGE_KEY = 'auth_state'

// 时间偏移量（提前 30 秒视为过期，用于容错）
const TOKEN_EXPIRY_SKEW = 30 * 1000

// 认证状态接口
interface AuthState {
  accessToken: string
  tokenType: string
  expiresAt: number // 毫秒时间戳
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
  const accessToken = ref<string | null>(null)
  const tokenType = ref<string>('Bearer')
  const expiresAt = ref<number | null>(null)
  const user = ref<User | null>(null)

  // ========== 计算属性 ==========
  
  /**
   * 判断 token 是否已过期
   */
  const isExpired = computed(() => {
    if (!expiresAt.value) return true
    return Date.now() >= expiresAt.value
  })

  /**
   * 判断用户是否已认证（有 token 且未过期）
   */
  const isAuthenticated = computed(() => {
    return !!accessToken.value && !isExpired.value
  })

  // ========== 操作方法 ==========

  /**
   * 登录成功后调用，保存 token 和用户信息
   * @param payload 包含 token 和 user 的登录响应数据
   */
  const loginSuccess = (payload: { token: TokenResponse; user: User }) => {
    const { token, user: userData } = payload
    
    // 计算过期时间戳（提前 skew 秒视为过期）
    const calculatedExpiresAt = Date.now() + token.expires_in * 1000 - TOKEN_EXPIRY_SKEW
    
    // 更新状态
    accessToken.value = token.access_token
    tokenType.value = token.token_type || 'Bearer'
    expiresAt.value = calculatedExpiresAt
    user.value = userData
    
    // 持久化到 sessionStorage
    saveToSession({
      accessToken: accessToken.value,
      tokenType: tokenType.value,
      expiresAt: expiresAt.value,
      user: user.value
    })
  }

  /**
   * 从 sessionStorage 初始化认证状态
   * 应用启动时调用，如果 token 已过期则自动清理
   */
  const initFromSession = () => {
    const stored = loadFromSession()
    if (!stored) return
    
    // 检查是否已过期
    if (Date.now() >= stored.expiresAt) {
      // 已过期，清理状态
      clearSession()
      return
    }
    
    // 恢复状态
    accessToken.value = stored.accessToken
    tokenType.value = stored.tokenType
    expiresAt.value = stored.expiresAt
    user.value = stored.user
  }

  /**
   * 登出，清空所有认证状态
   */
  const logout = () => {
    accessToken.value = null
    tokenType.value = 'Bearer'
    expiresAt.value = null
    user.value = null
    
    clearSession()
  }

  return {
    // 状态
    accessToken,
    tokenType,
    expiresAt,
    user,
    
    // 计算属性
    isExpired,
    isAuthenticated,
    
    // 操作方法
    loginSuccess,
    initFromSession,
    logout
  }
})
