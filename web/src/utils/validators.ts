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
 * 头像URL验证规则
 */
export const avatarRules: FormItemRule[] = [
  { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
]

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

