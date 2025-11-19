<template>
  <div class="register-page">
    <el-card class="register-page__card">
      <template #header>
        <div class="register-page__header">
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

        <el-form-item label-width="0" class="register-page__button-form-item">
          <el-button type="primary" :loading="loading" @click="handleRegister" class="register-page__button">
            {{ loading ? '注册中...' : '注册' }}
          </el-button>
        </el-form-item>

        <div class="register-page__login-link">
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
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, Avatar, Picture } from '@element-plus/icons-vue'
import { useAuth } from '@/composables/useAuth'
import { usernameRules, passwordRules, nicknameRules, avatarRules, createConfirmPasswordValidator } from '@/utils/validators'

const router = useRouter()
const registerFormRef = ref<FormInstance>()

// 表单数据
const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  nickname: '',
  avatar: ''
})

// 表单验证规则
const rules = reactive<FormRules>({
  username: usernameRules,
  password: passwordRules,
  confirmPassword: [createConfirmPasswordValidator(() => registerForm.password, true)],
  nickname: nicknameRules,
  avatar: avatarRules
})

// 使用 composable 管理认证逻辑
const { loading, handleRegister: register } = useAuth()

// 注册处理
const handleRegister = async () => {
  await register(registerFormRef.value, registerForm)
}

// 跳转到登录页
const goToLogin = () => {
  router.push('/login')
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--color-page-background-alt);
  background-image: 
    radial-gradient(circle at var(--pattern-size-dot-small) var(--pattern-size-dot-small), rgba(0, 0, 0, 0.03) 2%, transparent 0%),
    radial-gradient(circle at var(--pattern-size-dot-large) var(--pattern-size-dot-large), rgba(0, 0, 0, 0.03) 2%, transparent 0%);
  background-size: var(--pattern-size-grid) var(--pattern-size-grid);
  padding: var(--spacing-lg);
}

.register-page__card {
  width: 100%;
  max-width: var(--dialog-width-small);
  border-radius: var(--border-radius-xlarge);
  box-shadow: var(--shadow-auth-card);
  background: var(--color-card-background);
}

.register-page__header {
  text-align: center;
}

.register-page__header h2 {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: var(--font-size-heading);
  font-weight: 600;
  color: var(--color-text-primary);
}

.register-page__header p {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-tertiary);
}

.register-page__button-form-item {
  text-align: center;
}

.register-page__button-form-item :deep(.el-form-item__content) {
  display: flex;
  justify-content: center;
}

.register-page__button {
  width: var(--button-width-auth);
  margin-top: var(--spacing-sm-lg);
  height: var(--button-height-large);
  font-size: var(--font-size-base);
  font-weight: 500;
}

.register-page__login-link {
  text-align: center;
  margin-top: var(--spacing-md);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.register-page__card :deep(.el-form-item__label) {
  font-weight: 500;
}

.register-page__card :deep(.el-input__inner) {
  border-radius: var(--border-radius-input);
}

.register-page__card :deep(.el-button) {
  border-radius: var(--border-radius-button);
}

/* 响应式设计 */
/* 小屏幕：对应 --breakpoint-mobile (480px) */
@media (max-width: 480px) {
  .register-page {
    padding: var(--spacing-md);
  }

  .register-page__card {
    border-radius: var(--border-radius-large);
  }

  .register-page__header h2 {
    font-size: var(--font-size-title);
  }

  .register-page__button {
    width: 100%;
  }
}
</style>
