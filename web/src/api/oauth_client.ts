import request from './request'
import type { OAuthClientListResponse, CreateOAuthClientRequest } from '@/types/oauth_client'
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

