import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { listUsers } from '@/api/user'
import type { UserListResponse } from '@/types/user'

/**
 * 用户列表管理相关的组合式函数
 */
export function useUserList() {
  const loading = ref(false)
  const userList = ref<UserListResponse[]>([])

  const filters = ref({
    status: undefined as number | string | undefined,
    role: undefined as string | undefined
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

  return {
    loading,
    userList,
    filters,
    pagination,
    fetchUserList,
    handleFilterChange,
    handlePageChange,
    handleSizeChange
  }
}

