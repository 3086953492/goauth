import type { User } from './user'

export interface LoginRequest {
  username: string
  password: string
}

export interface AccessTokenResponse {
  access_token: string
  expires_in: number
}

export interface RefreshTokenResponse {
  refresh_token: string
  expires_in: number
}

export interface LoginResponse {
  user: User
  access_token: AccessTokenResponse
}
