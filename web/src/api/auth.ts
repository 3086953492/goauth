import request from './request'
import type { RegisterRequest, ApiResponse } from '@/types/auth'

/**
 * 用户注册
 */
export const register = (data: RegisterRequest): Promise<ApiResponse> => {
  return request({
    url: '/auth/register',
    method: 'post',
    data
  })
}

/**
 * 用户登录 (预留接口)
 */
export const login = (data: { username: string; password: string }): Promise<ApiResponse> => {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

