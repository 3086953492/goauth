import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

/**
 * 权限检查相关的组合式函数
 */
export function usePermission() {
  const router = useRouter()
  const route = useRoute()
  const authStore = useAuthStore()

  // 是否是管理员
  const isAdmin = computed(() => authStore.user?.role === 'admin')

  // 是否已登录（使用 isAuthenticated 包含过期判断）
  const isLoggedIn = computed(() => authStore.isAuthenticated)

  /**
   * 检查用户是否已登录
   */
  const checkLogin = (): boolean => {
    if (!authStore.isAuthenticated) {
      ElMessage.error('未登录或登录已过期，请先登录')
      authStore.logout() // 清理可能过期的状态
      router.push({
        path: '/login',
        query: { redirect: route.fullPath }
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

    ElMessage.error('您没有权限编辑其他用户的信息')
    router.push('/profile')
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
      ElMessage.error('您没有管理员权限')
      router.push('/home')
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

