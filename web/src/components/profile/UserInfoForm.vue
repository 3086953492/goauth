<template>
  <div class="user-info-form">
    <!-- 头像预览 -->
    <div class="avatar-preview">
      <el-avatar :size="100" :src="modelValue.avatar" :icon="Avatar" />
    </div>

    <!-- 用户信息表单 -->
    <div class="form-section">
      <h3 class="section-title">基本信息</h3>

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
  margin-bottom: 24px;
}

.avatar-preview {
  display: flex;
  justify-content: center;
  margin-bottom: 32px;
}

.avatar-preview :deep(.el-avatar) {
  border: 4px solid #f5f7fa;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.form-section {
  margin-bottom: 32px;
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

:deep(.el-radio-group) {
  display: flex;
  gap: 16px;
}

:deep(.el-radio) {
  margin-right: 0;
}
</style>
