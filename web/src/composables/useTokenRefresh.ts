import { ref } from 'vue'
import { refreshToken as refreshTokenApi } from '@/api/auth'
import { useAuthStore } from '@/stores/useAuthStore'
import { ElNotification } from 'element-plus'

/**
 * 令牌自动刷新调度器
 * 在访问令牌即将过期前静默刷新，刷新失败时提示用户
 */

/** 刷新提前量（毫秒），在 access token 过期前多久触发刷新 */
const REFRESH_LEAD_TIME_MS = 60 * 1000 // 60 秒

/** 最小刷新间隔（毫秒），防止过于频繁的刷新 */
const MIN_REFRESH_INTERVAL_MS = 10 * 1000 // 10 秒

/** 定时器 ID */
let refreshTimerId: ReturnType<typeof setTimeout> | null = null

/** 是否已经提示过刷新失败（防止重复提示） */
const hasNotifiedFailure = ref(false)

/**
 * 格式化过期时间为本地时间字符串
 */
const formatExpireTime = (timestamp: number): string => {
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

/**
 * 清理刷新定时器
 */
const clearRefreshTimer = () => {
  if (refreshTimerId !== null) {
    clearTimeout(refreshTimerId)
    refreshTimerId = null
  }
}

/**
 * 执行令牌刷新
 */
const doRefresh = async (): Promise<boolean> => {
  const authStore = useAuthStore()
  
  try {
    const response = await refreshTokenApi()
    
    if (response.data?.access_token_expire_at) {
      // 更新 store 中的过期时间
      authStore.updateAccessTokenExpireAt(response.data.access_token_expire_at)
      return true
    }
    return false
  } catch (error) {
    console.error('[TokenRefresh] 刷新令牌失败:', error)
    return false
  }
}

/**
 * 显示刷新失败提示
 */
const notifyRefreshFailure = () => {
  if (hasNotifiedFailure.value) return
  
  const authStore = useAuthStore()
  const expireAt = authStore.accessTokenExpireAt
  
  let expireTimeStr = '未知时间'
  if (expireAt) {
    expireTimeStr = formatExpireTime(expireAt)
  }
  
  ElNotification({
    title: '登录状态刷新失败',
    message: `当前访问令牌将于 ${expireTimeStr} 失效，请注意保存未提交的数据。`,
    type: 'warning',
    duration: 10000, // 10 秒后自动关闭
    position: 'top-right'
  })
  
  hasNotifiedFailure.value = true
}

/**
 * 调度下一次刷新
 * 根据当前 access token 过期时间计算下次刷新时机
 */
const scheduleNextRefresh = () => {
  const authStore = useAuthStore()
  
  // 清理已有定时器
  clearRefreshTimer()
  
  const accessExpireAt = authStore.accessTokenExpireAt
  const refreshExpireAt = authStore.refreshTokenExpireAt
  const now = Date.now()
  
  // 没有过期时间信息，无法调度
  if (!accessExpireAt || !refreshExpireAt) {
    return
  }
  
  // refresh token 已过期，不再尝试刷新
  if (refreshExpireAt <= now) {
    return
  }
  
  // 计算下次刷新时间点
  const nextRefreshTime = accessExpireAt - REFRESH_LEAD_TIME_MS
  
  // 计算延迟时间
  let delay = nextRefreshTime - now
  
  // 如果已经过了预定刷新时间但 access token 还没过期，立即刷新
  if (delay <= 0 && accessExpireAt > now) {
    delay = 0
  } else if (delay <= 0) {
    // access token 已过期，不再刷新
    return
  }
  
  // 确保不会过于频繁刷新
  if (delay < MIN_REFRESH_INTERVAL_MS && delay > 0) {
    delay = MIN_REFRESH_INTERVAL_MS
  }
  
  refreshTimerId = setTimeout(async () => {
    const success = await doRefresh()
    
    if (success) {
      // 刷新成功，重置失败提示标志，调度下一次刷新
      hasNotifiedFailure.value = false
      scheduleNextRefresh()
    } else {
      // 刷新失败，提示用户，不再继续调度
      notifyRefreshFailure()
    }
  }, delay)
}

/**
 * 启动令牌自动刷新
 * 在登录成功或应用初始化恢复登录态后调用
 */
export const startTokenRefresh = () => {
  hasNotifiedFailure.value = false
  scheduleNextRefresh()
}

/**
 * 停止令牌自动刷新
 * 在登出时调用
 */
export const stopTokenRefresh = () => {
  clearRefreshTimer()
  hasNotifiedFailure.value = false
}

/**
 * Token 刷新相关的组合式函数
 */
export function useTokenRefresh() {
  return {
    startTokenRefresh,
    stopTokenRefresh
  }
}

