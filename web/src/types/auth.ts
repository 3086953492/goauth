export interface RegisterRequest {
  username: string
  password: string
  confirm_password: string
  nickname: string
  avatar?: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface ApiResponse<T = any> {
  success: boolean
  message: string
  data: T
}

export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
  status: number
  role: string
  created_at: string
  updated_at: string
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

export interface UpdateUserRequest {
  nickname?: string
  avatar?: string
  password?: string
  confirm_password?: string
  status?: number
  role?: string
}

export interface UserListResponse {
  id: number
  nickname: string
  avatar: string
  status: number
  role: string
}

export interface PaginationResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

