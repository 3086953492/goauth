import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { emitAuthEvent } from '@/utils/authFeedbackBus'

/**
 * 权限检查相关的组合式函数
 */
export function usePermission() {
  const route = useRoute()
  const authStore = useAuthStore()

  // 是否是管理员
  const isAdmin = computed(() => authStore.user?.role === 'admin')

  // 是否已登录
  const isLoggedIn = computed(() => authStore.isAuthenticated)

  /**
   * 检查用户是否已登录
   */
  const checkLogin = (): boolean => {
    if (!authStore.isAuthenticated) {
      // 清理可能存在的无效状态
      authStore.logout()
      // 通过事件总线发出登录要求事件，由 AuthFeedbackProvider 统一处理提示和跳转
      emitAuthEvent('auth:login-required', {
        message: '未登录或登录已过期，请先登录',
        redirectPath: route.fullPath
      })
      return false
    }
    return true
  }

  /**
   * 检查用户是否有权限编辑指定用户
   * @param targetUserId 目标用户ID
   * @param currentUserId 当前用户ID
   */
  const checkEditPermission = (targetUserId: string | number, currentUserId?: string | number): boolean => {
    if (!checkLogin()) {
      return false
    }

    const currentId = currentUserId || authStore.user?.id.toString() || ''
    const targetId = targetUserId.toString()

    // 管理员可以编辑任何人
    if (isAdmin.value) {
      return true
    }

    // 普通用户只能编辑自己
    if (targetId === currentId) {
      return true
    }

    // 通过事件总线发出权限不足事件
    emitAuthEvent('auth:forbidden', {
      message: '您没有权限编辑其他用户的信息',
      redirectPath: '/profile'
    })
    return false
  }

  /**
   * 检查用户是否有管理员权限
   */
  const checkAdminPermission = (): boolean => {
    if (!checkLogin()) {
      return false
    }

    if (!isAdmin.value) {
      // 通过事件总线发出权限不足事件
      emitAuthEvent('auth:forbidden', {
        message: '您没有管理员权限',
        redirectPath: '/home'
      })
      return false
    }

    return true
  }

  return {
    isAdmin,
    isLoggedIn,
    checkLogin,
    checkEditPermission,
    checkAdminPermission
  }
}

