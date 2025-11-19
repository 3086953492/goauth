<template>
  <div class="profile-page">
    <Navbar />
    <div class="profile-page__container">
      <el-card class="profile-page__card" v-loading="pageLoading">
        <template #header>
          <div class="profile-page__header">
            <h2>{{ isEditingSelf ? '个人信息' : '编辑用户信息' }}</h2>
            <p>{{ isEditingSelf ? '管理您的账户信息' : `编辑用户 @${targetUser?.username}` }}</p>
          </div>
        </template>

        <div class="profile-page__content">
          <el-form ref="profileFormRef" :model="formData" :rules="formRules" label-width="90px" size="large">
            <UserInfoForm v-model="userInfo" :username="targetUser.username"
              :can-edit-permission="isAdmin && !isEditingSelf" />

            <PasswordForm v-model="passwordData" />

            <el-form-item label-width="0" class="profile-page__actions">
              <el-button type="primary" :loading="submitLoading" @click="handleSubmit" class="profile-page__button--submit">
                {{ submitLoading ? '保存中...' : '保存修改' }}
              </el-button>
              <el-button @click="handleCancel" class="profile-page__button--cancel">
                返回
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { createPasswordValidator, createConfirmPasswordValidator, nicknameRules, avatarRules } from '@/utils/validators'
import { useUserProfile } from '@/composables/useUserProfile'
import Navbar from '@/components/Navbar.vue'
import UserInfoForm from '@/components/profile/UserInfoForm.vue'
import PasswordForm from '@/components/profile/PasswordForm.vue'

const profileFormRef = ref<FormInstance>()

// 使用 composable 管理业务逻辑
const {
  pageLoading,
  submitLoading,
  targetUser,
  userInfo,
  passwordData,
  formData,
  isEditingSelf,
  isAdmin,
  loadUserInfo,
  submitForm,
  cancel
} = useUserProfile()

// 表单验证规则
const formRules = computed<FormRules>(() => ({
  nickname: nicknameRules,
  avatar: avatarRules,
  password: createPasswordValidator(),
  confirmPassword: [createConfirmPasswordValidator(() => passwordData.value.password)]
}))

// 提交表单
const handleSubmit = async () => {
  await submitForm(profileFormRef.value)
}

// 取消/返回
const handleCancel = () => {
  cancel()
}

// 页面加载
onMounted(() => {
  loadUserInfo()
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: var(--color-page-background-alt);
  background-image:
    radial-gradient(circle at var(--pattern-size-dot-small) var(--pattern-size-dot-small), rgba(0, 0, 0, 0.03) 2%, transparent 0%),
    radial-gradient(circle at var(--pattern-size-dot-large) var(--pattern-size-dot-large), rgba(0, 0, 0, 0.03) 2%, transparent 0%);
  background-size: var(--pattern-size-grid) var(--pattern-size-grid);
}

.profile-page__container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: var(--page-padding-top) var(--spacing-lg) var(--spacing-lg);
}

.profile-page__card {
  width: 100%;
  max-width: var(--dialog-width-medium);
  border-radius: var(--border-radius-xlarge);
  box-shadow: var(--shadow-large);
  background: var(--color-card-background);
}

.profile-page__header {
  text-align: center;
}

.profile-page__header h2 {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: var(--font-size-heading);
  font-weight: 600;
  color: var(--color-text-primary);
}

.profile-page__header p {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-tertiary);
}

.profile-page__content {
  padding: var(--spacing-sm) 0;
}

.profile-page__actions {
  text-align: center;
  margin-top: var(--spacing-xl);
}

.profile-page__actions :deep(.el-form-item__content) {
  display: flex;
  justify-content: center;
  gap: var(--spacing-md);
}

.profile-page__button--submit,
.profile-page__button--cancel {
  height: var(--button-height-large);
  font-size: var(--font-size-base);
  font-weight: 500;
  min-width: var(--button-min-width);
}

.profile-page__card :deep(.el-form-item__label) {
  font-weight: 500;
}

.profile-page__card :deep(.el-input__inner) {
  border-radius: var(--border-radius-input);
}

.profile-page__card :deep(.el-button) {
  border-radius: var(--border-radius-button);
}

/* 响应式设计 */
/* 移动端：对应 --breakpoint-mobile (480px) */
@media (max-width: 480px) {
  .profile-page__container {
    padding: var(--page-padding-top-mobile) var(--spacing-md) var(--spacing-md);
  }

  .profile-page__card {
    border-radius: var(--border-radius-large);
  }

  .profile-page__header h2 {
    font-size: var(--font-size-title);
  }

  .profile-page__actions :deep(.el-form-item__content) {
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .profile-page__button--submit,
  .profile-page__button--cancel {
    width: 100%;
  }
}
</style>
