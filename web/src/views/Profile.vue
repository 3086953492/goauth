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
        <!-- 头像预览 -->
        <div class="avatar-preview">
          <el-avatar :size="100" :src="profileForm.avatar || defaultAvatar" />
        </div>

        <!-- 基本信息表单 -->
        <el-form ref="profileFormRef" :model="profileForm" :rules="profileRules" label-width="90px" size="large">
          <div class="form-section">
            <h3 class="section-title">基本信息</h3>
            
            <el-form-item label="用户名" prop="username">
              <el-input v-model="targetUser.username" disabled :prefix-icon="User" />
            </el-form-item>

            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="profileForm.nickname" placeholder="请输入昵称" clearable :prefix-icon="Avatar" />
            </el-form-item>

            <el-form-item label="头像URL" prop="avatar">
              <el-input v-model="profileForm.avatar" placeholder="请输入头像URL" clearable :prefix-icon="Picture" />
            </el-form-item>
          </div>

          <!-- 密码修改区域 -->
          <div class="form-section">
            <div class="section-header" @click="showPasswordSection = !showPasswordSection">
              <h3 class="section-title">
                <el-icon class="collapse-icon" :class="{ 'is-active': showPasswordSection }">
                  <ArrowRight />
                </el-icon>
                修改密码
              </h3>
              <span class="section-hint">(可选，不修改请留空)</span>
            </div>

            <el-collapse-transition>
              <div v-show="showPasswordSection" class="password-fields">
                <el-form-item label="新密码" prop="password">
                  <el-input 
                    v-model="profileForm.password" 
                    type="password" 
                    placeholder="留空表示不修改密码" 
                    show-password 
                    clearable 
                    :prefix-icon="Lock" 
                  />
                </el-form-item>

                <el-form-item label="确认密码" prop="confirmPassword">
                  <el-input 
                    v-model="profileForm.confirmPassword" 
                    type="password" 
                    placeholder="请再次输入新密码" 
                    show-password 
                    clearable 
                    :prefix-icon="Lock" 
                  />
                </el-form-item>
              </div>
            </el-collapse-transition>
          </div>

          <!-- 管理员专属区域 -->
          <div v-if="isAdmin" class="form-section admin-section">
            <h3 class="section-title">管理员设置</h3>
            
            <el-form-item label="账户状态" prop="status">
              <el-radio-group v-model="profileForm.status">
                <el-radio :label="1">启用</el-radio>
                <el-radio :label="0">禁用</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item label="用户角色" prop="role">
              <el-radio-group v-model="profileForm.role">
                <el-radio label="user">普通用户</el-radio>
                <el-radio label="admin">管理员</el-radio>
              </el-radio-group>
            </el-form-item>
          </div>

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
import { reactive, ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, Avatar, Picture, ArrowRight } from '@element-plus/icons-vue'
import { getUserInfo, updateUser } from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import type { User as UserType, UpdateUserRequest } from '@/types/auth'
import Navbar from '@/components/Navbar.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const profileFormRef = ref<FormInstance>()
const pageLoading = ref(true)
const submitLoading = ref(false)
const showPasswordSection = ref(false)
const targetUser = ref<UserType>({} as UserType)
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

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

// 表单数据
const profileForm = reactive({
  nickname: '',
  avatar: '',
  password: '',
  confirmPassword: '',
  status: 1,
  role: 'user'
})

// 原始数据（用于对比变化）
let originalData: any = {}

// 自定义验证：确认密码
const validateConfirmPassword = (_rule: any, value: any, callback: any) => {
  if (profileForm.password && !value) {
    callback(new Error('请再次输入新密码'))
  } else if (value && value !== profileForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

// 表单验证规则
const profileRules = reactive<FormRules>({
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 1, max: 20, message: '昵称长度在 1 到 20 个字符', trigger: 'blur' }
  ],
  avatar: [
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' }
  ],
  password: [
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
})

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
      
      // 填充表单
      profileForm.nickname = response.data.nickname
      profileForm.avatar = response.data.avatar || ''
      profileForm.status = response.data.status
      profileForm.role = response.data.role
      
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
        
        if (profileForm.nickname !== originalData.nickname) {
          updateData.nickname = profileForm.nickname
        }
        
        if (profileForm.avatar !== originalData.avatar) {
          updateData.avatar = profileForm.avatar || undefined
        }
        
        // 密码修改
        if (profileForm.password) {
          updateData.password = profileForm.password
          updateData.confirm_password = profileForm.confirmPassword
        }
        
        // 管理员可修改状态和角色
        if (isAdmin.value) {
          if (profileForm.status !== originalData.status) {
            updateData.status = profileForm.status
          }
          if (profileForm.role !== originalData.role) {
            updateData.role = profileForm.role
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
        profileForm.password = ''
        profileForm.confirmPassword = ''
        showPasswordSection.value = false
        
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
  padding-bottom: 24px;
  border-bottom: 1px solid #ebeef5;
}

.form-section:last-of-type {
  border-bottom: none;
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-header {
  cursor: pointer;
  user-select: none;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-header:hover .section-title {
  color: #409eff;
}

.section-hint {
  font-size: 13px;
  color: #909399;
  font-weight: 400;
}

.collapse-icon {
  transition: transform 0.3s ease;
  font-size: 16px;
}

.collapse-icon.is-active {
  transform: rotate(90deg);
}

.password-fields {
  padding-top: 8px;
}

.admin-section {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 12px;
  border: 1px solid #e9ecef;
}

.admin-section .section-title {
  margin-top: 0;
  color: #e6a23c;
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

:deep(.el-radio-group) {
  display: flex;
  gap: 16px;
}

:deep(.el-radio) {
  margin-right: 0;
}
</style>

