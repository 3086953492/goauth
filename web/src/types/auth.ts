export interface RegisterRequest {
  username: string
  password: string
  confirm_password: string
  nickname: string
  avatar?: string
}

export interface ApiResponse<T = any> {
  code: number
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

