import type { User } from './user'

export interface LoginRequest {
  username: string
  password: string
}

/**
 * 登录响应（后端返回的用户信息 + 令牌过期时间）
 */
export interface LoginResponse {
  user: User
  /** 访问令牌过期时间（ISO 8601 格式字符串） */
  access_token_expire_at: string
  /** 刷新令牌过期时间（ISO 8601 格式字符串） */
  refresh_token_expire_at: string
}

/**
 * 刷新令牌响应
 */
export interface RefreshTokenResponse {
  /** 新的访问令牌过期时间（ISO 8601 格式字符串） */
  access_token_expire_at: string
}

