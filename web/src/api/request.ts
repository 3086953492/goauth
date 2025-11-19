import axios from 'axios'
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useAuthStore } from '@/stores/useAuthStore'
import { API_BASE_URL } from '@/constants'

// 防抖标志，避免多个401请求重复重定向
let isRedirecting = false

// 刷新 token 的锁，避免并发刷新
let isRefreshing = false
// 待重试的请求队列
let refreshSubscribers: Array<(token: string) => void> = []

/**
 * 添加请求到刷新队列
 */
const subscribeTokenRefresh = (callback: (token: string) => void) => {
  refreshSubscribers.push(callback)
}

/**
 * 刷新成功后，通知所有等待的请求
 */
const onRefreshed = (token: string) => {
  refreshSubscribers.forEach((callback) => callback(token))
  refreshSubscribers = []
}

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
    
    // 判断是否为 OAuth 授权接口
    if (isOAuthAuthorizationRequest(config.url)) {
      // OAuth 授权接口：使用 Cookie 鉴权，不带 Bearer 头
      config.withCredentials = true
      // 不设置 Authorization 头
    } else if (config.url !== '/api/v1/auth/login' && config.url !== '/api/v1/auth/refresh_token') {
      // 普通业务接口：使用 Bearer token（排除登录和刷新接口）
      const token = authStore.getAccessToken()
      if (token && config.headers) {
        config.headers.Authorization = `${authStore.tokenType} ${token}`
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
  async (error: AxiosError) => {
    console.error('Response error:', error)
    
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }
    
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
          
          // 如果是登录或刷新接口本身返回 401，直接跳转登录
          if (originalRequest.url === '/api/v1/auth/login' || originalRequest.url === '/api/v1/auth/refresh_token') {
            // 刷新 token 失败，清理状态并跳转登录
            if (originalRequest.url === '/api/v1/auth/refresh_token') {
              const authStore = useAuthStore()
              authStore.clearAuth()
              
              if (!isRedirecting) {
                isRedirecting = true
                router.replace('/login').finally(() => {
                  setTimeout(() => {
                    isRedirecting = false
                  }, 1000)
                })
              }
            }
            break
          }
          
          // 如果请求未标记为重试过，尝试刷新 token
          if (!originalRequest._retry) {
            // 如果正在刷新 token，将请求加入队列
            if (isRefreshing) {
              return new Promise((resolve) => {
                subscribeTokenRefresh((token: string) => {
                  if (originalRequest.headers) {
                    originalRequest.headers.Authorization = `Bearer ${token}`
                  }
                  resolve(request(originalRequest))
                })
              })
            }
            
            // 标记为重试，避免死循环
            originalRequest._retry = true
            isRefreshing = true
            
            try {
              // 调用刷新 token 接口
              const authStore = useAuthStore()
              const refreshResponse = await axios.post(
                `${API_BASE_URL}/api/v1/auth/refresh_token`,
                {},
                { withCredentials: true }
              )
              
              // 刷新成功，更新 token
              if (refreshResponse.data.success && refreshResponse.data.data) {
                const newAccessToken = refreshResponse.data.data.access_token
                const expiresIn = refreshResponse.data.data.expires_in
                
                // 更新 store 中的 token
                authStore.setAccessToken(newAccessToken, expiresIn)
                
                // 通知所有等待的请求
                onRefreshed(newAccessToken)
                
                // 更新原请求的 Authorization 头
                if (originalRequest.headers) {
                  originalRequest.headers.Authorization = `Bearer ${newAccessToken}`
                }
                
                // 重试原请求
                return request(originalRequest)
              } else {
                throw new Error('刷新令牌失败')
              }
            } catch (refreshError) {
              // 刷新失败，清理状态并跳转登录
              console.error('刷新令牌失败:', refreshError)
              const authStore = useAuthStore()
              authStore.clearAuth()
              
              if (!isRedirecting) {
                isRedirecting = true
                ElMessage.warning('登录已过期，请重新登录')
                router.replace('/login').finally(() => {
                  setTimeout(() => {
                    isRedirecting = false
                  }, 1000)
                })
              }
              
              return Promise.reject(refreshError)
            } finally {
              isRefreshing = false
            }
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
