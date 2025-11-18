import axios from 'axios'
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'
import { API_BASE_URL } from '@/constants'

// 防抖标志，避免多个401请求重复重定向
let isRedirecting = false

const request: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

/**
 * 判断是否为 OAuth 授权相关接口（使用 Cookie 鉴权，不带 Bearer 头）
 */
const isOAuthAuthorizationRequest = (url?: string): boolean => {
  if (!url) return false
  // OAuth 授权确认接口通过 HttpOnly Cookie 鉴权
  return url.includes('/oauth/authorize') && !url.includes('/oauth/authorize?')
}

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const authStore = useAuthStore()
    
    // 检查 token 是否已过期
    if (authStore.isExpired && authStore.accessToken) {
      // token 已过期，清理状态
      authStore.logout()
      ElMessage.warning('登录已过期，请重新登录')
      
      // 拒绝本次请求
      return Promise.reject(new Error('Token expired'))
    }
    
    // 判断是否为 OAuth 授权接口
    if (isOAuthAuthorizationRequest(config.url)) {
      // OAuth 授权接口：使用 Cookie 鉴权，不带 Bearer 头
      config.withCredentials = true
      // 不设置 Authorization 头
    } else {
      // 普通业务接口：使用 Bearer token
      if (authStore.isAuthenticated && config.headers) {
        config.headers.Authorization = `${authStore.tokenType} ${authStore.accessToken}`
      }
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
            
            // 清理认证状态
            const authStore = useAuthStore()
            authStore.logout()
            
            // 使用 replace 强制重定向到登录页
            router.replace('/login').finally(() => {
              // 重置防抖标志
              setTimeout(() => {
                isRedirecting = false
              }, 1000)
            })
          }
          
          // TODO: 未来实现 refresh token 自动刷新逻辑
          // 1. 检查是否为 access token 过期（通过后端返回的特定错误码判断）
          // 2. 如果是，调用 refresh 接口（withCredentials: true，依赖 HttpOnly Cookie）
          // 3. 刷新成功后，更新 store 中的 accessToken 和 expiresAt
          // 4. 重试原请求
          // 5. 刷新失败则执行当前的 logout + 跳转登录逻辑
          
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
