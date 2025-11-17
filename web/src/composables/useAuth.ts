import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { login as loginApi } from '@/api/auth'
import { register as registerApi } from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import type { RegisterRequest } from '@/types/user'

/**
 * 认证相关的组合式函数
 */
export function useAuth() {
  const router = useRouter()
  const authStore = useAuthStore()
  const loading = ref(false)

  /**
   * 处理登录
   */
  const handleLogin = async (
    formRef: FormInstance | undefined,
    loginForm: { username: string; password: string },
    redirectPath?: string
  ) => {
    if (!formRef) return

    await formRef.validate(async (valid) => {
      if (valid) {
        loading.value = true
        try {
          const response = await loginApi(loginForm)

          // 保存 token 和用户信息到 store
          authStore.setToken(response.data.token)
          authStore.setUser(response.data.user)

          ElMessage.success('登录成功')

          // 跳转到指定页面或默认主页
          router.push(redirectPath || '/home')
        } catch (error: any) {
          console.error('登录失败:', error)
          // 错误消息已在响应拦截器中统一显示，此处不再重复提示
        } finally {
          loading.value = false
        }
      } else {
        ElMessage.warning('请填写完整的登录信息')
      }
    })
  }

  /**
   * 处理注册
   */
  const handleRegister = async (
    formRef: FormInstance | undefined,
    registerForm: {
      username: string
      password: string
      confirmPassword: string
      nickname: string
      avatar: string
    }
  ) => {
    if (!formRef) return

    await formRef.validate(async (valid) => {
      if (valid) {
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

          ElMessage.success(response.message || '注册成功！')

          // 延迟跳转到登录页
          setTimeout(() => {
            router.push('/login')
          }, 1500)
        } catch (error: any) {
          console.error('注册失败:', error)
          // 错误已在拦截器中处理
        } finally {
          loading.value = false
        }
      } else {
        ElMessage.warning('请填写完整的表单信息')
      }
    })
  }

  /**
   * 处理登出
   */
  const handleLogout = () => {
    authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  }

  return {
    loading,
    handleLogin,
    handleRegister,
    handleLogout
  }
}

