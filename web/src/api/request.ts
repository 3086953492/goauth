import axios from 'axios'
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { useAuthStore } from '@/stores/useAuthStore'
import { API_BASE_URL } from '@/constants'

// 防抖标志，避免多个401请求重复重定向
let isRedirecting = false

// 防抖标志，避免多次跳转错误页
let isNavigatingToError = false

const request: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json'
  }
})


// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 全局使用 Cookie 鉴权，不再需要手动设置 Authorization 头
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
    
    // 严格以业务 success 字段为准
    if (res && res.success === true) {
      return res
    }
    
    // 处理业务错误（即便 HTTP 200，只要 success !== true 就视为错误）
    const message = res?.message || '请求失败'
    ElMessage.error(message)
    return Promise.reject(new Error(message))
  },
  async (error: AxiosError) => {
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
          // 401 未授权/登录过期，优先使用后端返回的 message
          errorMessage = errorMessage || '登录已过期，请重新登录'
          
          // 清理状态并跳转登录
          const authStore = useAuthStore()
          authStore.logout()
          
          if (!isRedirecting) {
            isRedirecting = true
            // 统一在此处提示一次登录失效消息
            ElMessage.warning(errorMessage)
            router.push({
              path: '/login',
              query: { redirect: router.currentRoute.value.fullPath }
            }).finally(() => {
              setTimeout(() => {
                isRedirecting = false
              }, 1000)
            })
          }
          break
        case 403:
          errorMessage = errorMessage || '拒绝访问'
          // 403 权限错误，跳转到错误页
          if (!isNavigatingToError) {
            isNavigatingToError = true
            router.push({
              name: 'Error',
              query: {
                error: '403',
                error_description: errorMessage
              }
            }).finally(() => {
              setTimeout(() => {
                isNavigatingToError = false
              }, 1000)
            })
          }
          break
        case 404:
          errorMessage = errorMessage || '请求的资源不存在'
          // 404 错误仅提示，不跳转（通常是接口不存在，不是页面不存在）
          ElMessage.error(errorMessage)
          break
        case 500:
        case 502:
        case 503:
          errorMessage = errorMessage || '服务器内部错误'
          // 5xx 服务器错误，跳转到错误页
          if (!isNavigatingToError) {
            isNavigatingToError = true
            router.push({
              name: 'Error',
              query: {
                error: String(status),
                error_description: errorMessage
              }
            }).finally(() => {
              setTimeout(() => {
                isNavigatingToError = false
              }, 1000)
            })
          }
          break
        default:
          errorMessage = errorMessage || '网络错误'
          ElMessage.error(errorMessage)
      }
      
      // 对于不跳转错误页的情况，显示消息提示
      if (status !== 403 && status !== 500 && status !== 502 && status !== 503) {
        // 404 等已经在上面处理过了，这里处理其他状态码
        if (status !== 404 && status !== 401) {
          ElMessage.error(errorMessage)
        }
      }
      
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
