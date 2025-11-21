<template>
  <!-- 该组件不渲染任何 UI，仅作为事件监听器 -->
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/useAuthStore'
import { onAuthEvent, offAuthEvent, type AuthEventType, type AuthEventPayload } from '@/utils/authFeedbackBus'

const router = useRouter()
const authStore = useAuthStore()

// 防抖标志，避免多个401请求重复重定向
let isRedirecting = false

// 防抖标志，避免多次跳转错误页
let isNavigatingToError = false

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
const handleLoginRequired = (payload: AuthEventPayload) => {
  if (isRedirecting) {
    return
  }

  isRedirecting = true

  // 清理登录态
  authStore.logout()

  // 显示提示
  ElMessage.warning(payload.message || '登录已过期，请重新登录')

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
  // 显示提示
  ElMessage.error(payload.message || '您没有权限执行此操作')

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

