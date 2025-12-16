<template>
  <div class="user-info-form">
    <!-- 头像展示区（置顶） -->
    <el-form-item label-width="0" prop="avatarFile" class="user-info-form__avatar-item">
      <div class="user-info-form__avatar-card">
        <AvatarUploadCard
          :model-value="avatarFile"
          :preview-url="currentAvatarUrl"
          :max-size="MAX_SIZE"
          :allowed-types="ALLOWED_TYPES"
          @update:model-value="onAvatarFileChange"
          @crop="openCropperDialog"
        />
      </div>
    </el-form-item>

    <!-- 用户信息表单 -->
    <div class="user-info-form__section">
      <h3 class="user-info-form__title">基本信息</h3>

      <el-form-item label="用户名" prop="username">
        <el-input :model-value="username || '加载中...'" disabled :prefix-icon="User" />
      </el-form-item>

      <el-form-item label="昵称" prop="nickname">
        <el-input
          :model-value="modelValue.nickname"
          @input="updateField('nickname', $event)"
          placeholder="请输入昵称"
          clearable
          :prefix-icon="AvatarIcon"
        />
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

    <!-- 裁剪弹窗 -->
    <AvatarCropperDialog
      v-model="cropperDialogVisible"
      :file="cropperSourceFile"
      @confirm="onCropConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { User, Avatar as AvatarIcon } from '@element-plus/icons-vue'
import { USER_ROLES, USER_STATUS } from '@/constants'
import { AvatarUploadCard, AvatarCropperDialog } from '@/components/base/avatar'

interface UserInfo {
  nickname: string
  avatar: string
  status: number
  role: string
}

interface Props {
  modelValue: UserInfo
  /** 新头像文件 (v-model:avatarFile) */
  avatarFile?: File | null
  username?: string
  canEditPermission: boolean
}

const props = withDefaults(defineProps<Props>(), {
  avatarFile: null
})

const emit = defineEmits<{
  'update:modelValue': [value: UserInfo]
  'update:avatarFile': [value: File | null]
}>()

const userRoleOptions = USER_ROLES
const userStatusOptions = USER_STATUS

// 允许的图片类型与最大大小（4MB）
const ALLOWED_TYPES = ['image/png', 'image/jpeg', 'image/jpg', 'image/webp']
const MAX_SIZE = 4 * 1024 * 1024

// 当前头像 URL（旧头像），当没有选择新文件时显示
const currentAvatarUrl = computed(() => props.modelValue.avatar || null)

/**
 * 响应 AvatarUploadCard 的 v-model 更新
 */
const onAvatarFileChange = (file: File | null) => {
  emit('update:avatarFile', file)
}

// ========== 裁剪弹窗相关 ==========
const cropperDialogVisible = ref(false)
const cropperSourceFile = ref<File | null>(null)

/** 打开裁剪弹窗 */
const openCropperDialog = () => {
  if (!props.avatarFile) return
  cropperSourceFile.value = props.avatarFile
  cropperDialogVisible.value = true
}

/** 裁剪确认回调 */
const onCropConfirm = (croppedFile: File) => {
  emit('update:avatarFile', croppedFile)
}

const updateField = (field: keyof UserInfo, value: unknown) => {
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

/* 头像展示区：居中，与下方内容留出间距 */
.user-info-form__avatar-item {
  margin-bottom: var(--spacing-xl);
}

.user-info-form__avatar-item :deep(.el-form-item__content) {
  justify-content: center;
}

.user-info-form__avatar-card {
  display: flex;
  justify-content: center;
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
