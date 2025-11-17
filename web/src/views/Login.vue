<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h2>用户登录</h2>
          <p>欢迎回来</p>
        </div>
      </template>

      <el-form ref="loginFormRef" :model="loginForm" :rules="rules" label-width="80px" size="large">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" show-password clearable />
        </el-form-item>

        <el-form-item label-width="0" class="button-form-item">
          <el-button type="primary" :loading="loading" @click="handleLogin" class="login-button">
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>

        <div class="register-link">
          还没有账户？
          <el-link type="primary" @click="goToRegister">立即注册</el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { useAuth } from '@/composables/useAuth'
import { usernameRules, passwordRules } from '@/utils/validators'

const router = useRouter()
const route = useRoute()
const loginFormRef = ref<FormInstance>()

// 获取重定向地址
const redirect = computed(() => route.query.redirect as string | undefined)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = reactive<FormRules>({
  username: usernameRules,
  password: passwordRules
})

// 使用 composable 管理认证逻辑
const { loading, handleLogin: login } = useAuth()

const handleLogin = async () => {
  await login(loginFormRef.value, loginForm, redirect.value)
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<style scoped>
.login-container {
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

.login-card {
  width: 100%;
  max-width: 450px;
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

.login-button {
  width: 200px;
  margin-top: 12px;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
}

.register-link {
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
