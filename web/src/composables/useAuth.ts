import { ref } from 'vue'
import { login as loginApi } from '@/api/auth'
import { register as registerApi } from '@/api/user'
import { useAuthStore } from '@/stores/useAuthStore'
import type { RegisterRequest } from '@/types/user'

/**
 * 认证操作的返回结果类型
 */
export interface AuthActionResult {
  success: boolean
  message?: string
  errorMessage?: string
}

/**
 * 认证相关的组合式函数
 * 只负责业务逻辑，不涉及 UI 提示和路由跳转
 */
export function useAuth() {
  const authStore = useAuthStore()
  const loading = ref(false)

  /**
   * 处理登录
   * @param loginForm 登录表单数据
   * @returns 登录结果
   */
  const handleLogin = async (loginForm: {
    username: string
    password: string
  }): Promise<AuthActionResult> => {
    loading.value = true
    try {
      const response = await loginApi(loginForm)

      // 后端返回 { user, access_token: { access_token, expires_in } }
      // refresh_token 在 HttpOnly Cookie 中，前端不保存
      authStore.loginSuccess({
        access_token: response.data.access_token,
        user: response.data.user
      })

      return {
        success: true,
        message: '登录成功'
      }
    } catch (error: any) {
      console.error('登录失败:', error)
      return {
        success: false,
        errorMessage: error.message || '登录失败'
      }
    } finally {
      loading.value = false
    }
  }

  /**
   * 处理注册
   * @param registerForm 注册表单数据
   * @returns 注册结果
   */
  const handleRegister = async (registerForm: {
    username: string
    password: string
    confirmPassword: string
    nickname: string
    avatar: string
  }): Promise<AuthActionResult> => {
    loading.value = true
    try {
      const requestData: RegisterRequest = {
        username: registerForm.username,
        password: registerForm.password,
        confirm_password: registerForm.confirmPassword,
        nickname: registerForm.nickname,
        avatar: registerForm.avatar || undefined
      }

      const response = await registerApi(requestData)

      return {
        success: true,
        message: response.message || '注册成功！'
      }
    } catch (error: any) {
      console.error('注册失败:', error)
      return {
        success: false,
        errorMessage: error.message || '注册失败'
      }
    } finally {
      loading.value = false
    }
  }

  /**
   * 处理登出
   * @returns 登出结果
   */
  const handleLogout = (): AuthActionResult => {
    authStore.logout()
    return {
      success: true,
      message: '已退出登录'
    }
  }

  return {
    loading,
    handleLogin,
    handleRegister,
    handleLogout
  }
}

