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
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getUserInfo, updateUser } from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import type { User as UserType, UpdateUserRequest } from '@/types/user'
import { createPasswordValidator, createConfirmPasswordValidator } from '@/utils/validators'
import Navbar from '@/components/Navbar.vue'
import UserInfoForm from '@/components/profile/UserInfoForm.vue'
import PasswordForm from '@/components/profile/PasswordForm.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const profileFormRef = ref<FormInstance>()
const pageLoading = ref(true)
const submitLoading = ref(false)
const targetUser = ref<UserType>({} as UserType)

// 目标用户ID
const targetUserId = computed(() => {
  const routeId = route.params.id as string
  return routeId || authStore.user?.id.toString() || ''
})

// 是否是管理员
const isAdmin = computed(() => authStore.user?.role === 'admin')

// 是否在编辑自己
const isEditingSelf = computed(() => {
  return targetUserId.value === authStore.user?.id.toString()
})

// 用户信息数据
const userInfo = ref({
  nickname: '',
  avatar: '',
  status: 1,
  role: 'user'
})

// 密码数据
const passwordData = ref({
  password: '',
  confirmPassword: ''
})

// 合并的表单数据（用于验证）
const formData = computed(() => ({
  ...userInfo.value,
  ...passwordData.value
}))

// 原始数据（用于对比变化）
let originalData: any = {}

// 表单验证规则
const formRules = computed<FormRules>(() => ({
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 1, max: 20, message: '昵称长度在 1 到 20 个字符', trigger: 'blur' }
  ],
  avatar: [
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ],
  password: createPasswordValidator(),
  confirmPassword: [
    { validator: createConfirmPasswordValidator(() => passwordData.value.password), trigger: 'blur' }
  ]
}))

// 权限检查
const checkPermission = () => {
  if (!authStore.user) {
    ElMessage.error('未登录，请先登录')
    router.push('/login')
    return false
  }

  // 非管理员尝试编辑其他用户
  if (!isAdmin.value && !isEditingSelf.value) {
    ElMessage.error('您没有权限编辑其他用户的信息')
    router.push(`/profile`)
    return false
  }

  return true
}

// 加载用户信息
const loadUserInfo = async () => {
  if (!checkPermission()) {
    return
  }

  try {
    pageLoading.value = true
    const response = await getUserInfo(targetUserId.value)

    if (response.data) {
      targetUser.value = response.data

      // 填充用户信息表单
      userInfo.value = {
        nickname: response.data.nickname,
        avatar: response.data.avatar || '',
        status: response.data.status,
        role: response.data.role
      }

      // 保存原始数据
      originalData = {
        nickname: response.data.nickname,
        avatar: response.data.avatar || '',
        status: response.data.status,
        role: response.data.role
      }
    }
  } catch (error: any) {
    console.error('加载用户信息失败:', error)
    ElMessage.error('加载用户信息失败')
    router.push('/home')
  } finally {
    pageLoading.value = false
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!profileFormRef.value) return

  await profileFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        // 构建更新数据 - 只发送已修改的字段
        const updateData: UpdateUserRequest = {}

        if (userInfo.value.nickname !== originalData.nickname) {
          updateData.nickname = userInfo.value.nickname
        }

        if (userInfo.value.avatar !== originalData.avatar) {
          updateData.avatar = userInfo.value.avatar || undefined
        }

        // 密码修改
        if (passwordData.value.password) {
          updateData.password = passwordData.value.password
          updateData.confirm_password = passwordData.value.confirmPassword
        }

        // 管理员可修改状态和角色
        if (isAdmin.value) {
          if (userInfo.value.status !== originalData.status) {
            updateData.status = userInfo.value.status
          }
          if (userInfo.value.role !== originalData.role) {
            updateData.role = userInfo.value.role
          }
        }

        // 如果没有任何修改
        if (Object.keys(updateData).length === 0) {
          ElMessage.info('没有任何修改')
          return
        }

        const response = await updateUser(targetUserId.value, updateData)

        ElMessage.success(response.message || '更新成功！')

        // 如果编辑的是当前用户，更新 authStore
        if (isEditingSelf.value) {
          const updatedUserResponse = await getUserInfo(targetUserId.value)
          if (updatedUserResponse.data) {
            authStore.setUser(updatedUserResponse.data)
          }
        }

        // 清空密码字段
        passwordData.value.password = ''
        passwordData.value.confirmPassword = ''

        // 重新加载用户信息
        await loadUserInfo()

      } catch (error: any) {
        console.error('更新失败:', error)
        // 错误已在拦截器中处理
      } finally {
        submitLoading.value = false
      }
    } else {
      ElMessage.warning('请填写正确的表单信息')
    }
  })
}

// 取消/返回
const handleCancel = () => {
  router.back()
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
