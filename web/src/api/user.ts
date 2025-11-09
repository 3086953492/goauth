import request from './request'
import type { User, ApiResponse, RegisterRequest, UpdateUserRequest, UserListResponse, PaginationResponse } from '@/types/auth'

/**
 * 用户注册
 */
export const register = (data: RegisterRequest): Promise<ApiResponse> => {
  return request({
    url: '/api/v1/users',
    method: 'post',
    data
  })
}

/**
 * 获取用户信息
 */
export const getUserInfo = (userId: string | number): Promise<ApiResponse<User>> => {
  return request({
    url: `/api/v1/users/${userId}`,
    method: 'get'
  })
}

/**
 * 更新用户信息
 */
export const updateUser = (userId: string | number, data: UpdateUserRequest): Promise<ApiResponse> => {
  return request({
    url: `/api/v1/users/${userId}`,
    method: 'patch',
    data
  })
}

/**
 * 获取用户列表
 */
export const listUsers = (params?: {
  page?: number
  page_size?: number
  status?: number | string
  role?: string
}): Promise<ApiResponse<PaginationResponse<UserListResponse>>> => {
  return request({
    url: '/api/v1/users',
    method: 'get',
    params
  })
}

