import axios from 'axios'
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'

// 防抖标志，避免多个401请求重复重定向
let isRedirecting = false

const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:9000',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 可以在这里添加 token
    const token = localStorage.getItem('token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: AxiosError) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    
    // 处理成功响应
    if (res.success === true || response.status === 200) {
      return res
    }
    
    // 处理业务错误
    ElMessage.error(res.message || '请求失败')
    return Promise.reject(new Error(res.message || '请求失败'))
  },
  (error: AxiosError) => {
    console.error('Response error:', error)
    
    // 处理 HTTP 错误
    if (error.response) {
      const { status, data } = error.response as any
      
      // 优先使用后端返回的错误消息
      let errorMessage = data?.message || ''
      
      switch (status) {
        case 400:
          errorMessage = errorMessage || '请求参数错误'
          break
        case 401:
          errorMessage = errorMessage || '未授权，请登录'
          
          // 防止多次重定向
          if (!isRedirecting) {
            isRedirecting = true
            
            // 1. 直接清理localStorage（主要机制）
            localStorage.removeItem('token')
            localStorage.removeItem('refreshToken')
            localStorage.removeItem('expiresIn')
            localStorage.removeItem('user')
            
            // 2. 尝试清理store状态
            try {
              const authStore = useAuthStore()
              authStore.logout()
            } catch (e) {
              console.warn('清理认证状态失败:', e)
            }
            
            // 3. 使用replace强制重定向到登录页
            router.replace('/login').finally(() => {
              // 重置防抖标志
              setTimeout(() => {
                isRedirecting = false
              }, 1000)
            })
          }
          break
        case 403:
          errorMessage = errorMessage || '拒绝访问'
          break
        case 404:
          errorMessage = errorMessage || '请求的资源不存在'
          break
        case 500:
          errorMessage = errorMessage || '服务器内部错误'
          break
        default:
          errorMessage = errorMessage || '网络错误'
      }
      
      ElMessage.error(errorMessage)
      
      // 将错误消息附加到 error 对象上，方便业务层使用
      const errorWithMessage = new Error(errorMessage)
      Object.assign(errorWithMessage, { response: error.response, message: errorMessage })
      return Promise.reject(errorWithMessage)
    } else if (error.request) {
      ElMessage.error('网络连接失败，请检查网络')
    } else {
      ElMessage.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

export default request

