<template>
  <!-- 该组件不渲染任何 UI，仅作为事件监听器 -->
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuth } from '@/composables/useAuth'
import { onAuthEvent, offAuthEvent, type AuthEventType, type AuthEventPayload } from '@/utils/authFeedbackBus'

const router = useRouter()
const { handleLogout } = useAuth()

// 防抖标志，避免多个401请求重复重定向
let isRedirecting = false

// 防抖标志，避免多次跳转错误页
let isNavigatingToError = false

// 认证 toast 时间窗配置（毫秒）
const AUTH_TOAST_WINDOW_MS = 3000

// 记录每类认证提示的最近一次展示时间戳
const lastAuthToastTime: Record<string, number> = {}

/**
 * 判断是否应该展示认证 toast（基于时间窗防抖）
 * @param type 认证事件类型
 * @param message 提示消息
 * @returns 是否应该展示 toast
 */
const shouldShowAuthToast = (type: AuthEventType, message: string): boolean => {
  const key = `${type}:${message}`
  const now = Date.now()
  const last = lastAuthToastTime[key]

  // 如果不存在记录或已超过时间窗，则允许展示
  if (!last || now - last >= AUTH_TOAST_WINDOW_MS) {
    lastAuthToastTime[key] = now
    return true
  }

  // 在时间窗内，不重复展示
  return false
}

/**
 * 统一的认证事件处理器
 */
const handleAuthEvent = (type: AuthEventType, payload: AuthEventPayload) => {
  switch (type) {
    case 'auth:login-required':
    case 'auth:login-expired':
      handleLoginRequired(payload)
      break
    case 'auth:forbidden':
      handleForbidden(payload)
      break
    default:
      console.warn('[AuthFeedbackProvider] Unknown auth event type:', type)
  }
}

/**
 * 处理登录失效/未登录事件
 */
const handleLoginRequired = async (payload: AuthEventPayload) => {
  // 计算最终提示消息
  const message = payload.message || '登录已过期，请重新登录'

  // 时间窗防抖：判断是否应该展示 toast
  const shouldShowToast = shouldShowAuthToast('auth:login-expired', message)

  // 导航防抖：避免多次重定向
  if (isRedirecting) {
    return
  }

  isRedirecting = true

  // 调用登出接口并清理登录态（包含停止令牌刷新、调用后端登出接口、清除前端状态）
  await handleLogout()

  // 仅在时间窗允许时展示提示
  if (shouldShowToast) {
    ElMessage.warning(message)
  }

  // 跳转到登录页，保留重定向路径
  router.push({
    path: '/login',
    query: { redirect: payload.redirectPath || router.currentRoute.value.fullPath }
  }).finally(() => {
    setTimeout(() => {
      isRedirecting = false
    }, 1000)
  })
}

/**
 * 处理权限不足事件
 */
const handleForbidden = (payload: AuthEventPayload) => {
  // 计算最终提示消息
  const message = payload.message || '您没有权限执行此操作'

  // 时间窗防抖：判断是否应该展示 toast
  const shouldShowToast = shouldShowAuthToast('auth:forbidden', message)

  // 仅在时间窗允许时展示提示
  if (shouldShowToast) {
    ElMessage.error(message)
  }

  // 如果有错误码（如 403），跳转到错误页
  if (payload.errorCode) {
    if (isNavigatingToError) {
      return
    }

    isNavigatingToError = true

    router.push({
      name: 'Error',
      query: {
        error: payload.errorCode,
        error_description: payload.errorDescription || payload.message
      }
    }).finally(() => {
      setTimeout(() => {
        isNavigatingToError = false
      }, 1000)
    })
  } else if (payload.redirectPath) {
    // 如果指定了安全跳转路径（如 /home 或 /profile），则跳转
    router.push(payload.redirectPath)
  }
}

// 组件挂载时订阅事件
onMounted(() => {
  onAuthEvent(handleAuthEvent)
})

// 组件卸载时取消订阅
onUnmounted(() => {
  offAuthEvent(handleAuthEvent)
})
</script>

