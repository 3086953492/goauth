// 密码验证规则
export const createPasswordValidator = () => [
  { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
]

// 确认密码验证规则工厂
export const createConfirmPasswordValidator = (getPassword: () => string) => {
  return (_rule: any, value: any, callback: any) => {
    const password = getPassword()
    if (password && !value) {
      callback(new Error('请再次输入新密码'))
    } else if (value && value !== password) {
      callback(new Error('两次输入的密码不一致'))
    } else {
      callback()
    }
  }
}

