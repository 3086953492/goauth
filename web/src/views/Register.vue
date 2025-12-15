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
          <div class="register-page__avatar-upload">
            <div v-if="avatarPreview" class="register-page__avatar-preview">
              <el-avatar :size="80" :src="avatarPreview" />
              <el-button class="register-page__avatar-remove" type="danger" circle size="small" @click="removeAvatar">
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
            <el-button v-else type="default" @click="triggerFileInput">
              <el-icon class="register-page__upload-icon"><Plus /></el-icon>
              选择头像（可选）
            </el-button>
            <input
              ref="fileInputRef"
              type="file"
              accept="image/png,image/jpeg,image/jpg,image/webp"
              class="register-page__file-input"
              @change="handleFileChange"
            />
          </div>
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
import { reactive, ref, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, Avatar, Plus, Close } from '@element-plus/icons-vue'
import { useAuth } from '@/composables/useAuth'
import { usernameRules, passwordRules, nicknameRules, createConfirmPasswordValidator, createAvatarFileValidator } from '@/utils/validators'

const router = useRouter()
const registerFormRef = ref<FormInstance>()
const fileInputRef = ref<HTMLInputElement | null>(null)

// 表单数据（avatar 改为 File | null）
const registerForm = reactive<{
  username: string
  password: string
  confirmPassword: string
  nickname: string
  avatar: File | null
}>({
  username: '',
  password: '',
  confirmPassword: '',
  nickname: '',
  avatar: null
})

// 头像本地预览 URL
const avatarPreviewUrl = ref<string | null>(null)
const avatarPreview = computed(() => avatarPreviewUrl.value)

// 清理预览 URL 防止内存泄漏
onUnmounted(() => {
  if (avatarPreviewUrl.value) {
    URL.revokeObjectURL(avatarPreviewUrl.value)
  }
})

// 触发原生文件选择
const triggerFileInput = () => {
  fileInputRef.value?.click()
}

// 允许的图片类型与最大大小（4MB）
const ALLOWED_TYPES = ['image/png', 'image/jpeg', 'image/jpg', 'image/webp']
const MAX_SIZE = 4 * 1024 * 1024

// 处理文件选择
const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  // 验证类型
  if (!ALLOWED_TYPES.includes(file.type)) {
    ElMessage.warning('仅支持 PNG、JPG、JPEG、WebP 格式的图片')
    target.value = ''
    return
  }
  // 验证大小
  if (file.size > MAX_SIZE) {
    ElMessage.warning('头像文件不能超过 4MB')
    target.value = ''
    return
  }

  // 释放旧的预览 URL
  if (avatarPreviewUrl.value) {
    URL.revokeObjectURL(avatarPreviewUrl.value)
  }
  registerForm.avatar = file
  avatarPreviewUrl.value = URL.createObjectURL(file)
}

// 移除已选头像
const removeAvatar = () => {
  if (avatarPreviewUrl.value) {
    URL.revokeObjectURL(avatarPreviewUrl.value)
  }
  registerForm.avatar = null
  avatarPreviewUrl.value = null
  if (fileInputRef.value) {
    fileInputRef.value.value = ''
  }
}

// 表单验证规则
const rules = reactive<FormRules>({
  username: usernameRules,
  password: passwordRules,
  confirmPassword: [createConfirmPasswordValidator(() => registerForm.password, true)],
  nickname: nicknameRules,
  avatar: [createAvatarFileValidator()]
})

// 使用 composable 管理认证逻辑
const { loading, handleRegister: register } = useAuth()

// 注册处理
const handleRegister = async () => {
  if (!registerFormRef.value) return

  // 执行表单校验
  const valid = await registerFormRef.value.validate().catch(() => false)
  
  if (!valid) {
    ElMessage.warning('请填写完整的表单信息')
    return
  }

  // 调用注册逻辑
  const result = await register(registerForm)

  if (result.success) {
    // 注册成功：显示提示并延迟跳转到登录页
    ElMessage.success(result.message || '注册成功！')
    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } else {
    // 注册失败：错误已在拦截器中处理，这里可选择性提示
    // 如果需要额外提示可以取消注释下面这行
    // ElMessage.error(result.errorMessage || '注册失败')
  }
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

/* 头像上传 */
.register-page__avatar-upload {
  display: flex;
  align-items: center;
}

.register-page__avatar-preview {
  position: relative;
  display: inline-block;
}

.register-page__avatar-remove {
  position: absolute;
  top: -8px;
  right: -8px;
}

.register-page__upload-icon {
  margin-right: var(--spacing-xs);
}

.register-page__file-input {
  display: none;
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
