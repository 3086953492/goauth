import request from './request'
import type { ApiResponse } from '@/types/common'

/**
 * OAuth 授权请求参数
 */
export interface OAuthAuthorizeRequest {
  client_id: string
  redirect_uri: string
  response_type: string
  scope?: string
  state?: string
}

/**
 * OAuth 授权响应
 */
export interface OAuthAuthorizeResponse {
  redirect_uri: string
}

/**
 * 确认 OAuth 授权
 * 注意：此接口使用 Cookie 鉴权（HttpOnly），不使用 Bearer token
 */
export const confirmAuthorize = (data: OAuthAuthorizeRequest): Promise<ApiResponse<OAuthAuthorizeResponse>> => {
  return request({
    url: '/api/v1/oauth/authorize',
    method: 'post',
    data,
    withCredentials: true // 确保携带 Cookie
  })
}

