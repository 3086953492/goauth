import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { listOAuthClients } from '@/api/oauth_client'
import type { OAuthClientListResponse } from '@/types/oauth_client'

/**
 * OAuth 客户端列表管理相关的组合式函数
 */
export function useOAuthClientList() {
  const loading = ref(false)
  const clientList = ref<OAuthClientListResponse[]>([])

  const filters = ref({
    name: undefined as string | undefined,
    status: undefined as number | string | undefined
  })

  const pagination = ref({
    page: 1,
    pageSize: 10,
    total: 0
  })

  /**
   * 获取 OAuth 客户端列表
   */
  const fetchClientList = async () => {
    loading.value = true
    try {
      const params: any = {
        page: pagination.value.page,
        page_size: pagination.value.pageSize
      }

      if (filters.value.name !== undefined && filters.value.name !== '') {
        params.name = filters.value.name
      }
      if (filters.value.status !== undefined && filters.value.status !== '') {
        params.status = filters.value.status
      }

      const response = await listOAuthClients(params)
      if (response.success && response.data) {
        clientList.value = response.data.items
        pagination.value.total = response.data.total
        pagination.value.page = response.data.page
        pagination.value.pageSize = response.data.pageSize
      } else {
        ElMessage.error(response.message || '获取 OAuth 客户端列表失败')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '获取 OAuth 客户端列表失败')
    } finally {
      loading.value = false
    }
  }

  /**
   * 处理筛选变化
   */
  const handleFilterChange = () => {
    pagination.value.page = 1
    fetchClientList()
  }

  /**
   * 处理页码变化
   */
  const handlePageChange = (page: number) => {
    pagination.value.page = page
    fetchClientList()
  }

  /**
   * 处理页面大小变化
   */
  const handleSizeChange = (size: number) => {
    pagination.value.pageSize = size
    pagination.value.page = 1
    fetchClientList()
  }

  return {
    loading,
    clientList,
    filters,
    pagination,
    fetchClientList,
    handleFilterChange,
    handlePageChange,
    handleSizeChange
  }
}

