import request from './request'
import type { LoginRequest } from '@/types/auth'
import type { User } from '@/types/user'
import type { ApiResponse } from '@/types/common'

/**
 * 用户登录
 */
export const login = (data: LoginRequest): Promise<ApiResponse<User>> => {
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
