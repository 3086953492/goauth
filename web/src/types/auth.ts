import type { User } from './user'

export interface LoginRequest {
  username: string
  password: string
}

export interface TokenResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface LoginResponse {
  user: User
  token: TokenResponse
}
