<template>
  <div class="login-page">
    <el-card class="login-page__card">
      <template #header>
        <div class="login-page__header">
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

        <el-form-item label-width="0" class="login-page__button-form-item">
          <el-button type="primary" :loading="loading" @click="handleLogin" class="login-page__button">
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>

        <div class="login-page__register-link">
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
import { ElMessage } from 'element-plus'
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
  if (!loginFormRef.value) return

  // 执行表单校验
  const valid = await loginFormRef.value.validate().catch(() => false)
  
  if (!valid) {
    ElMessage.warning('请填写完整的登录信息')
    return
  }

  // 调用登录逻辑
  const result = await login(loginForm)

  if (result.success) {
    // 登录成功：显示提示并跳转
    ElMessage.success(result.message || '登录成功')
    router.push(redirect.value || '/home')
  } else {
    // 登录失败：错误已在拦截器中处理，这里可选择性提示
    // 如果需要额外提示可以取消注释下面这行
    // ElMessage.error(result.errorMessage || '登录失败')
  }
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<style scoped>
.login-page {
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

.login-page__card {
  width: 100%;
  max-width: var(--container-max-width-auth);
  border-radius: var(--border-radius-xlarge);
  box-shadow: var(--shadow-auth-card);
  background: var(--color-card-background);
}

.login-page__header {
  text-align: center;
}

.login-page__header h2 {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: var(--font-size-display);
  font-weight: 600;
  color: var(--color-text-primary);
}

.login-page__header p {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-tertiary);
}

.login-page__button-form-item {
  text-align: center;
}

.login-page__button-form-item :deep(.el-form-item__content) {
  display: flex;
  justify-content: center;
}

.login-page__button {
  width: var(--button-width-auth);
  margin-top: var(--spacing-sm);
  height: var(--button-height-large);
  font-size: var(--font-size-base);
  font-weight: 500;
}

.login-page__register-link {
  text-align: center;
  margin-top: var(--spacing-md);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-input__inner) {
  border-radius: var(--border-radius-input);
}

:deep(.el-button) {
  border-radius: var(--border-radius-button);
}

/* 响应式设计 */
/* 小屏幕：对应 --breakpoint-mobile (480px) */
@media (max-width: 480px) {
  .login-page {
    padding: var(--spacing-md);
  }

  .login-page__card {
    border-radius: var(--border-radius-large);
  }

  .login-page__header h2 {
    font-size: var(--font-size-title);
  }

  .login-page__button {
    width: 100%;
  }
}
</style>
