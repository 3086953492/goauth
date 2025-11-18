import request from './request'
import type { LoginRequest, LoginResponse, AccessTokenResponse } from '@/types/auth'
import type { ApiResponse } from '@/types/common'

/**
 * 用户登录
 */
export const login = (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
  return request({
    url: '/api/v1/auth/login',
    method: 'post',
    data
  })
}

/**
 * 刷新访问令牌
 * 通过 HttpOnly Cookie 中的 refresh_token 刷新 access_token
 */
export const refreshToken = (): Promise<ApiResponse<AccessTokenResponse>> => {
  return request({
    url: '/api/v1/auth/refresh_token',
    method: 'post',
    withCredentials: true  // 携带 Cookie
  })
}

