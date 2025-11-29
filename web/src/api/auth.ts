import request from './request'
import type { LoginRequest, LoginResponse, RefreshTokenResponse } from '@/types/auth'
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
 * 退出登录
 */
export const logout = (): Promise<ApiResponse> => {
  return request({
    url: '/api/v1/auth/logout',
    method: 'post'
  })
}

/**
 * 刷新访问令牌
 * 后端通过 HTTP-only Cookie 更新令牌，前端只需关心过期时间
 */
export const refreshToken = (): Promise<ApiResponse<RefreshTokenResponse>> => {
  return request({
    url: '/api/v1/auth/refresh_token',
    method: 'post'
  })
}
