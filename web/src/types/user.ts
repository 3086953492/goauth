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

/**
 * 注册表单 UI 层数据（驼峰命名，头像为 File）
 */
export interface RegisterFormValues {
  username: string
  password: string
  confirmPassword: string
  nickname: string
  avatar?: File | null
}

/**
 * @deprecated 旧的 JSON 请求类型，现已改用 multipart/form-data 提交
 */
export interface RegisterRequest {
  username: string
  password: string
  confirm_password: string
  nickname: string
  avatar?: string
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

