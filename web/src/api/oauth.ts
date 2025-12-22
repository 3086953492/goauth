import { API_BASE_URL } from '@/constants'
/**
 * OAuth 授权参数
 */
export interface OAuthAuthorizationParams {
  client_id: string
  redirect_uri: string
  response_type: string
  scope?: string
  state?: string
}

/**
 * 构建 OAuth 授权 URL
 * 用于直接跳转到后端授权接口，由后端完成重定向
 */
export const buildAuthorizationUrl = (params: OAuthAuthorizationParams): string => {
  const url = new URL(`${API_BASE_URL}/api/v1/oauth/authorize`, window.location.origin)
  
  url.searchParams.set('response_type', params.response_type)
  url.searchParams.set('client_id', params.client_id)
  url.searchParams.set('redirect_uri', params.redirect_uri)
  
  if (params.scope) {
    url.searchParams.set('scope', params.scope)
  }
  
  if (params.state) {
    url.searchParams.set('state', params.state)
  }
  
  return url.toString()
}

