<template>
  <div class="user-info-form">
    <!-- 头像预览 -->
    <div class="user-info-form__avatar-preview">
      <el-avatar :size="avatarSize" :src="modelValue.avatar" :icon="Avatar" />
    </div>

    <!-- 用户信息表单 -->
    <div class="user-info-form__section">
      <h3 class="user-info-form__title">基本信息</h3>

      <el-form-item label="用户名" prop="username">
        <el-input :model-value="username || '加载中...'" disabled :prefix-icon="User" />
      </el-form-item>

      <el-form-item label="昵称" prop="nickname">
        <el-input :model-value="modelValue.nickname" @input="updateField('nickname', $event)" placeholder="请输入昵称"
          clearable :prefix-icon="Avatar" />
      </el-form-item>

      <el-form-item label="头像URL" prop="avatar">
        <el-input :model-value="modelValue.avatar" @input="updateField('avatar', $event)" placeholder="请输入头像URL"
          clearable :prefix-icon="Picture" />
      </el-form-item>

      <!-- 仅管理员可编辑的字段 -->
      <template v-if="canEditPermission">
        <el-form-item label="账户状态" prop="status">
          <el-radio-group :model-value="modelValue.status" @update:model-value="updateField('status', $event)">
            <el-radio v-for="status in userStatusOptions" :key="status.value" :label="status.value">
              {{ status.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="用户角色" prop="role">
          <el-radio-group :model-value="modelValue.role" @update:model-value="updateField('role', $event)">
            <el-radio v-for="role in userRoleOptions" :key="role.value" :label="role.value">
              {{ role.label }}
            </el-radio>
          </el-radio-group>
        </el-form-item>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { User, Avatar, Picture } from '@element-plus/icons-vue'
import { USER_ROLES, USER_STATUS } from '@/constants'

interface UserInfo {
  nickname: string
  avatar: string
  status: number
  role: string
}

interface Props {
  modelValue: UserInfo
  username?: string
  canEditPermission: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:modelValue', value: UserInfo): void
}>()

// 头像尺寸
const avatarSize = 100

const userRoleOptions = USER_ROLES
const userStatusOptions = USER_STATUS

const updateField = (field: keyof UserInfo, value: any) => {
  emit('update:modelValue', {
    ...props.modelValue,
    [field]: value
  })
}
</script>

<style scoped>
.user-info-form {
  margin-bottom: var(--spacing-lg);
}

.user-info-form__avatar-preview {
  display: flex;
  justify-content: center;
  margin-bottom: var(--spacing-xl);
}

.user-info-form__avatar-preview :deep(.el-avatar) {
  border: var(--border-width-thick) solid var(--color-background-lighter);
  box-shadow: var(--shadow-avatar);
}

.user-info-form__section {
  margin-bottom: var(--spacing-xl);
}

.user-info-form__title {
  margin: 0 0 var(--spacing-lg) 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text-primary);
}

.user-info-form :deep(.el-radio-group) {
  display: flex;
  gap: var(--spacing-md);
}

.user-info-form :deep(.el-radio) {
  margin-right: 0;
}
</style>
