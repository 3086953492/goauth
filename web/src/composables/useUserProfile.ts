import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { getUserInfo, updateUser } from '@/api/user'
import { useAuthStore } from '@/stores/useAuthStore'
import type { User as UserType, UpdateUserFormValues } from '@/types/user'
import { usePermission } from './usePermission'
import { detectChanges } from '@/utils/dataTransform'

/**
 * 用户信息管理相关的组合式函数
 */
export function useUserProfile() {
  const router = useRouter()
  const route = useRoute()
  const authStore = useAuthStore()
  const { isAdmin, checkEditPermission } = usePermission()

  const pageLoading = ref(true)
  const submitLoading = ref(false)
  const targetUser = ref<UserType>({} as UserType)

  // 用户信息数据（不再包含 avatar URL，avatar 只用于显示当前头像）
  const userInfo = ref({
    nickname: '',
    avatar: '',
    status: 1,
    role: 'user'
  })

  // 新头像文件（可选）
  const avatarFile = ref<File | null>(null)

  // 密码数据
  const passwordData = ref({
    password: '',
    confirmPassword: ''
  })

  // 原始数据（用于对比变化，不包含 avatar）
  let originalData: { nickname: string; status: number; role: string } = {
    nickname: '',
    status: 1,
    role: 'user'
  }

  // 目标用户ID
  const targetUserId = computed(() => {
    const routeId = route.params.id as string
    return routeId || authStore.user?.id.toString() || ''
  })

  // 是否在编辑自己
  const isEditingSelf = computed(() => {
    return targetUserId.value === authStore.user?.id.toString()
  })

  // 合并的表单数据（用于验证）
  const formData = computed(() => ({
    ...userInfo.value,
    avatarFile: avatarFile.value,
    ...passwordData.value
  }))

  /**
   * 加载用户信息
   */
  const loadUserInfo = async () => {
    if (!checkEditPermission(targetUserId.value)) {
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

        // 保存原始数据（不包含 avatar，因为 avatar 通过文件上传）
        originalData = {
          nickname: response.data.nickname,
          status: response.data.status,
          role: response.data.role
        }

        // 清空头像文件选择
        avatarFile.value = null
      }
    } catch (error: unknown) {
      console.error('加载用户信息失败:', error)
      ElMessage.error('加载用户信息失败')
      void router.push('/home')
    } finally {
      pageLoading.value = false
    }
  }

  /**
   * 提交表单
   */
  const submitForm = async (profileFormRef: FormInstance | undefined) => {
    if (!profileFormRef) return

    await profileFormRef.validate(async (valid) => {
      if (valid) {
        submitLoading.value = true
        try {
          // 构建更新数据
          const updateData: UpdateUserFormValues = {}

          // 检测基本信息变化（不包含 avatar）
          const changes = detectChanges(originalData, userInfo.value, [
            'nickname',
            'status',
            'role'
          ])

          // 昵称变更
          if ('nickname' in changes) {
            updateData.nickname = changes.nickname as string
          }

          // 密码修改
          if (passwordData.value.password) {
            updateData.password = passwordData.value.password
            updateData.confirm_password = passwordData.value.confirmPassword
          }

          // 管理员才能修改状态和角色
          if (isAdmin.value) {
            if ('status' in changes) {
              updateData.status = changes.status as number
            }
            if ('role' in changes) {
              updateData.role = changes.role as string
            }
          }

          // 新头像文件
          if (avatarFile.value) {
            updateData.avatar = avatarFile.value
          }

          // 如果没有任何修改
          if (Object.keys(updateData).length === 0) {
            ElMessage.info('没有任何修改')
            return
          }

          const response = await updateUser(targetUserId.value, updateData)

          ElMessage.success(response.message || '更新成功！')

          // 如果编辑的是当前用户，更新 authStore 中的用户信息
          if (isEditingSelf.value) {
            const updatedUserResponse = await getUserInfo(targetUserId.value)
            if (updatedUserResponse.data) {
              authStore.user = updatedUserResponse.data
            }
          }

          // 清空密码字段
          passwordData.value.password = ''
          passwordData.value.confirmPassword = ''

          // 重新加载用户信息
          await loadUserInfo()
        } catch (error: unknown) {
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

  /**
   * 取消/返回
   */
  const cancel = () => {
    router.back()
  }

  return {
    pageLoading,
    submitLoading,
    targetUser,
    userInfo,
    avatarFile,
    passwordData,
    formData,
    targetUserId,
    isEditingSelf,
    isAdmin,
    loadUserInfo,
    submitForm,
    cancel
  }
}
