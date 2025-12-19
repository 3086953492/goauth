import axios from 'axios'
import type { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { emitAuthEvent } from '@/utils/authFeedbackBus'
import { API_BASE_URL } from '@/constants'
import type { ApiResponse, ProblemDetails } from '@/types/common'

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
    // 如果请求体是 FormData，删除默认 Content-Type，让浏览器自动带上 multipart/form-data + boundary
    if (config.data instanceof FormData) {
      delete config.headers['Content-Type']
    }
    // 全局使用 Cookie 鉴权，不再需要手动设置 Authorization 头
    return config
  },
  (error: AxiosError) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

/**
 * 检测响应是否符合 ApiResponse 结构
 * 注意：data 字段是可选的，后端在无数据时可能省略此字段
 */
function isApiResponse(data: unknown): data is ApiResponse {
  return (
    typeof data === 'object' &&
    data !== null &&
    'code' in data &&
    'message' in data &&
    typeof (data as ApiResponse).code === 'number' &&
    typeof (data as ApiResponse).message === 'string'
  )
}

/**
 * 从 Problem JSON 或其他格式中提取错误消息
 * 按 RFC7807 优先读取 detail
 */
function extractErrorMessage(data: unknown): string {
  if (typeof data === 'object' && data !== null) {
    const problem = data as Partial<ProblemDetails>
    if (typeof problem.detail === 'string' && problem.detail) {
      return problem.detail
    }
    if (typeof problem.title === 'string' && problem.title) {
      return problem.title
    }
  }
  return ''
}

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data

    // 严格校验响应结构：必须是 { code, message, data? }
    if (!isApiResponse(res)) {
      const errorMsg = '响应格式错误'
      ElMessage.error(errorMsg)
      return Promise.reject(new Error(errorMsg))
    }

    // HTTP 2xx 且符合协议即视为成功，直接返回完整信封
    // 使用 as any 绕过 AxiosResponse 类型检查，因为我们需要返回解包后的业务数据
    return res as any
  },
  async (error: AxiosError) => {
    console.error('Response error:', error)

    // 处理 HTTP 错误
    if (error.response) {
      const { status, data } = error.response

      // 从 Problem JSON 提取错误消息
      let errorMessage = extractErrorMessage(data)

      switch (status) {
        case 400:
          errorMessage = errorMessage || '请求参数错误'
          ElMessage.error(errorMessage)
          break
        case 401:
          // 401 未授权/登录过期
          errorMessage = errorMessage || '登录已过期，请重新登录'

          // 通过事件总线发出登录失效事件，由 AuthFeedbackProvider 统一处理
          if (!isRedirecting) {
            isRedirecting = true
            emitAuthEvent('auth:login-expired', {
              message: errorMessage,
              redirectPath: router.currentRoute.value.fullPath
            })
            // 防抖重置
            setTimeout(() => {
              isRedirecting = false
            }, 1000)
          }
          break
        case 403:
          errorMessage = errorMessage || '拒绝访问'
          // 403 权限错误，通过事件总线发出权限不足事件
          if (!isNavigatingToError) {
            isNavigatingToError = true
            emitAuthEvent('auth:forbidden', {
              message: errorMessage,
              errorCode: '403',
              errorDescription: errorMessage
            })
            // 防抖重置
            setTimeout(() => {
              isNavigatingToError = false
            }, 1000)
          }
          break
        case 404:
          errorMessage = errorMessage || '请求的资源不存在'
          // 404 错误仅提示，不跳转
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
