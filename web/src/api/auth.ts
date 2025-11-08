import request from './request'
import type { RegisterRequest, LoginRequest, ApiResponse, LoginResponse } from '@/types/auth'

/**
 * 用户注册
 */
export const register = (data: RegisterRequest): Promise<ApiResponse> => {
  return request({
    url: '/users',
    method: 'post',
    data
  })
}

/**
 * 用户登录
 */
export const login = (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

