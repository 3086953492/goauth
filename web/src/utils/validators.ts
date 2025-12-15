import type { FormItemRule } from 'element-plus'

/**
 * 用户名验证规则
 */
export const usernameRules: FormItemRule[] = [
  { required: true, message: '请输入用户名', trigger: 'blur' },
  { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' },
  { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
]

/**
 * 昵称验证规则
 */
export const nicknameRules: FormItemRule[] = [
  { required: true, message: '请输入昵称', trigger: 'blur' },
  { min: 1, max: 100, message: '昵称长度在 1 到 100 个字符', trigger: 'blur' }
]

/**
 * 密码验证规则（用于登录，必填）
 */
export const passwordRules: FormItemRule[] = [
  { required: true, message: '请输入密码', trigger: 'blur' },
  { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
]

/**
 * 密码验证规则（可选，用于修改密码场景）
 */
export const createPasswordValidator = (): FormItemRule[] => [
  { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
]

/**
 * 头像URL验证规则（用于 URL 输入场景，如个人资料编辑）
 */
export const avatarRules: FormItemRule[] = [
  { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
]

// 头像文件允许类型与最大大小（4MB）
const AVATAR_ALLOWED_TYPES = ['image/png', 'image/jpeg', 'image/jpg', 'image/webp']
const AVATAR_MAX_SIZE = 4 * 1024 * 1024

/**
 * 头像文件验证规则工厂（用于文件上传场景）
 * 头像可选，只有选择文件时才校验类型/大小
 */
export const createAvatarFileValidator = (): FormItemRule => {
  return {
    validator: (_rule: any, value: any, callback: any) => {
      // 头像可选，未选择文件时跳过校验
      if (!value) {
        callback()
        return
      }
      // 类型检查
      if (!AVATAR_ALLOWED_TYPES.includes(value.type)) {
        callback(new Error('仅支持 PNG、JPG、JPEG、WebP 格式的图片'))
        return
      }
      // 大小检查
      if (value.size > AVATAR_MAX_SIZE) {
        callback(new Error('头像文件不能超过 4MB'))
        return
      }
      callback()
    },
    trigger: 'change'
  }
}

/**
 * 确认密码验证规则工厂
 * @param getPassword 获取原密码的函数
 * @param isRequired 是否必填（默认 false）
 */
export const createConfirmPasswordValidator = (
  getPassword: () => string,
  isRequired: boolean = false
): FormItemRule => {
  return {
    validator: (_rule: any, value: any, callback: any) => {
      const password = getPassword()
      
      // 如果密码为空且不是必填，则跳过验证
      if (!password && !isRequired) {
        callback()
        return
      }
      
      // 如果密码不为空，但确认密码为空
      if (password && !value) {
        callback(new Error(isRequired ? '请输入密码' : '请再次输入新密码'))
        return
      }
      
      // 验证两次密码是否一致
      if (value && value !== password) {
        callback(new Error('两次输入的密码不一致'))
        return
      }
      
      callback()
    },
    trigger: 'blur'
  }
}

