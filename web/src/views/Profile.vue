<template>
  <div class="profile-wrapper">
    <Navbar />
    <div class="profile-container">
      <el-card class="profile-card" v-loading="pageLoading">
        <template #header>
          <div class="card-header">
            <h2>{{ isEditingSelf ? '个人信息' : '编辑用户信息' }}</h2>
            <p>{{ isEditingSelf ? '管理您的账户信息' : `编辑用户 @${targetUser?.username}` }}</p>
          </div>
        </template>

        <div class="form-content">
          <el-form ref="profileFormRef" :model="formData" :rules="formRules" label-width="90px" size="large">
            <UserInfoForm v-model="userInfo" :username="targetUser.username"
              :can-edit-permission="isAdmin && !isEditingSelf" />

            <PasswordForm v-model="passwordData" />

            <el-form-item label-width="0" class="button-form-item">
              <el-button type="primary" :loading="submitLoading" @click="handleSubmit" class="submit-button">
                {{ submitLoading ? '保存中...' : '保存修改' }}
              </el-button>
              <el-button @click="handleCancel" class="cancel-button">
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
.profile-wrapper {
  min-height: 100vh;
  background: #f0f2f5;
  background-image:
    radial-gradient(circle at 25px 25px, rgba(0, 0, 0, 0.03) 2%, transparent 0%),
    radial-gradient(circle at 75px 75px, rgba(0, 0, 0, 0.03) 2%, transparent 0%);
  background-size: 100px 100px;
}

.profile-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 84px 20px 20px;
}

.profile-card {
  width: 100%;
  max-width: 700px;
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

.form-content {
  padding: 8px 0;
}

.button-form-item {
  text-align: center;
  margin-top: 32px;
}

.button-form-item :deep(.el-form-item__content) {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.submit-button,
.cancel-button {
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  min-width: 140px;
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
