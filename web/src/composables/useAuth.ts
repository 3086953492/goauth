import { ref } from 'vue'
import { login as loginApi, logout as logoutApi } from '@/api/auth'
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

      // 后端返回 ApiResponse<User>，data 即为 User 对象
      if (response.data) {
        authStore.loginSuccess(response.data)
      } else {
         throw new Error('登录返回数据异常')
      }

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
  const handleLogout = async (): Promise<AuthActionResult> => {
    try {
      // 调用后端登出接口清理 Cookie
      await logoutApi()
    } catch (error) {
      console.error('登出接口调用失败:', error)
      // 即使后端接口调用失败，前端也要执行登出操作
    } finally {
      authStore.logout()
    }
    
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
