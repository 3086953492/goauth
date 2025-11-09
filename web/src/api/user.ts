import request from './request'
import type { User, RegisterRequest, UpdateUserRequest, UserListResponse } from '@/types/user'
import type { ApiResponse, PaginationResponse } from '@/types/common'

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

/**
 * 删除用户
 */
export const deleteUser = (userId: string | number): Promise<ApiResponse> => {
  return request({
    url: `/api/v1/users/${userId}`,
    method: 'delete'
  })
}

