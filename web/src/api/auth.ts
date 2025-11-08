import request from './request'
import type { LoginRequest, ApiResponse, LoginResponse } from '@/types/auth'

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

