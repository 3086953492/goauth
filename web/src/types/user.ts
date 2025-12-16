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
 * 更新用户表单数据（用于 multipart/form-data 提交）
 * 字段名与后端 UpdateUserForm 一致
 */
export interface UpdateUserFormValues {
  nickname?: string
  password?: string
  confirm_password?: string
  /** 状态：0=禁用，1=启用 */
  status?: number
  role?: string
  /** 新头像文件（可选） */
  avatar?: File | null
}

export interface UserListResponse {
  id: number
  nickname: string
  avatar: string
  status: number
  role: string
}

