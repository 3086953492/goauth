<template>
  <div class="register-container">
    <el-card class="register-card">
      <template #header>
        <div class="card-header">
          <h2>用户注册</h2>
          <p>创建您的账户</p>
        </div>
      </template>

      <el-form ref="registerFormRef" :model="registerForm" :rules="rules" label-width="80px" size="large">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="registerForm.username" placeholder="请输入用户名" clearable :prefix-icon="User" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入密码" show-password clearable
            :prefix-icon="Lock" />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="registerForm.confirmPassword" type="password" placeholder="请再次输入密码" show-password clearable
            :prefix-icon="Lock" />
        </el-form-item>

        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="registerForm.nickname" placeholder="请输入昵称" clearable :prefix-icon="Avatar" />
        </el-form-item>

        <el-form-item label="头像" prop="avatar">
          <el-input v-model="registerForm.avatar" placeholder="请输入头像URL（可选）" clearable :prefix-icon="Picture" />
        </el-form-item>

        <el-form-item label-width="0" class="button-form-item">
          <el-button type="primary" :loading="loading" @click="handleRegister" class="register-button">
            {{ loading ? '注册中...' : '注册' }}
          </el-button>
        </el-form-item>

        <div class="login-link">
          已有账户？
          <el-link type="primary" @click="goToLogin">立即登录</el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, Avatar, Picture } from '@element-plus/icons-vue'
import { register } from '@/api/auth'
import type { RegisterRequest } from '@/types/auth'

const router = useRouter()
const registerFormRef = ref<FormInstance>()
const loading = ref(false)

// 表单数据
const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  nickname: '',
  avatar: ''
})

// 自定义验证：确认密码
const validateConfirmPassword = (_rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

// 表单验证规则
const rules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度至少 6 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 100, message: '昵称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  avatar: [
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ]
})

// 注册处理
const handleRegister = async () => {
  if (!registerFormRef.value) return

  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const requestData: RegisterRequest = {
          username: registerForm.username,
          password: registerForm.password,
          nickname: registerForm.nickname,
          avatar: registerForm.avatar || undefined
        }

        const response = await register(requestData)
        
        ElMessage.success(response.message || '注册成功！')
        
        // 延迟跳转到登录页
        setTimeout(() => {
          router.push('/login')
        }, 1500)
      } catch (error: any) {
        console.error('注册失败:', error)
        // 错误已在拦截器中处理
      } finally {
        loading.value = false
      }
    } else {
      ElMessage.warning('请填写完整的表单信息')
    }
  })
}

// 跳转到登录页
const goToLogin = () => {
  router.push('/login')
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f0f2f5;
  background-image: 
    radial-gradient(circle at 25px 25px, rgba(0, 0, 0, 0.03) 2%, transparent 0%),
    radial-gradient(circle at 75px 75px, rgba(0, 0, 0, 0.03) 2%, transparent 0%);
  background-size: 100px 100px;
  padding: 20px;
}

.register-card {
  width: 100%;
  max-width: 500px;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08);
  background: #ffffff;
}

.card-header {
  text-align: center;
}

.card-header h2 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.card-header p {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

.button-form-item {
  text-align: center;
}

.button-form-item :deep(.el-form-item__content) {
  display: flex;
  justify-content: center;
}

.register-button {
  width: 200px;
  margin-top: 12px;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
}

.login-link {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;
  color: #606266;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-input__inner) {
  border-radius: 6px;
}

:deep(.el-button) {
  border-radius: 6px;
}
</style>
