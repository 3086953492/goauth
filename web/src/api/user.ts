import request from './request'
import type { User, RegisterFormValues, UpdateUserRequest, UserListResponse } from '@/types/user'
import type { ApiResponse, PaginationResponse } from '@/types/common'

/**
 * 用户注册（multipart/form-data 提交）
 * @param values UI 表单数据（驼峰命名），包含可选的 avatar 文件
 */
export const register = (values: RegisterFormValues): Promise<ApiResponse> => {
  const formData = new FormData()
  formData.append('username', values.username)
  formData.append('password', values.password)
  formData.append('confirm_password', values.confirmPassword)
  formData.append('nickname', values.nickname)
  if (values.avatar) {
    formData.append('avatar', values.avatar)
  }
  return request({
    url: '/api/v1/users',
    method: 'post',
    data: formData
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
  nickname?: string
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

