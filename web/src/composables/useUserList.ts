import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listUsers, deleteUser } from '@/api/user'
import type { UserListResponse } from '@/types/user'
import { useAuthStore } from '@/stores/auth'

/**
 * 用户列表管理相关的组合式函数
 */
export function useUserList() {
  const loading = ref(false)
  const userList = ref<UserListResponse[]>([])

  const filters = ref({
    status: undefined as number | string | undefined,
    role: undefined as string | undefined,
    nickname: undefined as string | undefined
  })

  const pagination = ref({
    page: 1,
    pageSize: 10,
    total: 0
  })

  /**
   * 获取用户列表
   */
  const fetchUserList = async () => {
    loading.value = true
    try {
      const params: any = {
        page: pagination.value.page,
        page_size: pagination.value.pageSize
      }

      if (filters.value.status !== undefined && filters.value.status !== '') {
        params.status = filters.value.status
      }
      if (filters.value.role !== undefined && filters.value.role !== '') {
        params.role = filters.value.role
      }
      if (filters.value.nickname !== undefined && filters.value.nickname !== '') {
        params.nickname = filters.value.nickname
      }

      const response = await listUsers(params)
      if (response.success && response.data) {
        userList.value = response.data.items
        pagination.value.total = response.data.total
        pagination.value.page = response.data.page
        pagination.value.pageSize = response.data.pageSize
      } else {
        ElMessage.error(response.message || '获取用户列表失败')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '获取用户列表失败')
    } finally {
      loading.value = false
    }
  }

  /**
   * 处理筛选变化
   */
  const handleFilterChange = () => {
    pagination.value.page = 1
    fetchUserList()
  }

  /**
   * 处理页码变化
   */
  const handlePageChange = (page: number) => {
    pagination.value.page = page
    fetchUserList()
  }

  /**
   * 处理页面大小变化
   */
  const handleSizeChange = (size: number) => {
    pagination.value.pageSize = size
    pagination.value.page = 1
    fetchUserList()
  }

  /**
   * 删除用户
   */
  const handleDeleteUser = async (userId: number) => {
    const authStore = useAuthStore()
    
    // 检查是否删除当前登录用户
    if (authStore.user?.id === userId) {
      ElMessage.warning('不能删除当前登录用户')
      return
    }

    try {
      await ElMessageBox.confirm(
        '确定要删除该用户吗？此操作不可恢复。',
        '删除用户',
        {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning',
          confirmButtonClass: 'el-button--danger'
        }
      )

      // 执行删除
      const response = await deleteUser(userId)
      if (response.success) {
        ElMessage.success('删除用户成功')
        // 重新获取用户列表
        await fetchUserList()
      } else {
        ElMessage.error(response.message || '删除用户失败')
      }
    } catch (error: any) {
      // 用户取消操作
      if (error === 'cancel') {
        return
      }
      ElMessage.error(error.message || '删除用户失败')
    }
  }

  return {
    loading,
    userList,
    filters,
    pagination,
    fetchUserList,
    handleFilterChange,
    handlePageChange,
    handleSizeChange,
    handleDeleteUser
  }
}

