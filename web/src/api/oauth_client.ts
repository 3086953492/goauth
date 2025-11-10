import request from './request'
import type { OAuthClientListResponse, CreateOAuthClientRequest, UpdateOAuthClientRequest, OAuthClientDetailResponse } from '@/types/oauth_client'
import type { ApiResponse, PaginationResponse } from '@/types/common'

/**
 * 获取 OAuth 客户端列表
 */
export const listOAuthClients = (params?: {
  page?: number
  page_size?: number
  name?: string
  status?: number | string
}): Promise<ApiResponse<PaginationResponse<OAuthClientListResponse>>> => {
  return request({
    url: '/api/v1/oauth_clients',
    method: 'get',
    params
  })
}

/**
 * 创建 OAuth 客户端
 */
export const createOAuthClient = (data: CreateOAuthClientRequest): Promise<ApiResponse<null>> => {
  return request({
    url: '/api/v1/oauth_clients',
    method: 'post',
    data
  })
}

/**
 * 获取 OAuth 客户端详情
 */
export const getOAuthClient = (id: number): Promise<ApiResponse<OAuthClientDetailResponse>> => {
  return request({
    url: `/api/v1/oauth_clients/${id}`,
    method: 'get'
  })
}

/**
 * 更新 OAuth 客户端
 */
export const updateOAuthClient = (id: number, data: UpdateOAuthClientRequest): Promise<ApiResponse<null>> => {
  return request({
    url: `/api/v1/oauth_clients/${id}`,
    method: 'patch',
    data
  })
}

/**
 * 删除 OAuth 客户端
 */
export const deleteOAuthClient = (id: number): Promise<ApiResponse<null>> => {
  return request({
    url: `/api/v1/oauth_clients/${id}`,
    method: 'delete'
  })
}

