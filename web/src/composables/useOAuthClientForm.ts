import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { createOAuthClient } from '@/api/oauth_client'
import type { CreateOAuthClientRequest } from '@/types/oauth_client'

/**
 * OAuth 客户端表单管理相关的组合式函数
 */
export function useOAuthClientForm() {
  const loading = ref(false)

  const formData = reactive<CreateOAuthClientRequest>({
    name: '',
    client_secret: '',
    description: '',
    logo: '',
    redirect_uris: [''],
    grant_types: [],
    scopes: ['profile'],
    status: 1
  })

  /**
   * 生成随机的 client_secret（32位字符串）
   */
  const generateClientSecret = (): string => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    let result = ''
    for (let i = 0; i < 32; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    return result
  }

  /**
   * 初始化 client_secret
   */
  const initClientSecret = () => {
    formData.client_secret = generateClientSecret()
  }

  /**
   * 重新生成 client_secret
   */
  const regenerateClientSecret = () => {
    formData.client_secret = generateClientSecret()
    ElMessage.success('密钥已重新生成')
  }

  /**
   * 添加回调地址
   */
  const addRedirectUri = () => {
    formData.redirect_uris.push('')
  }

  /**
   * 删除回调地址
   */
  const removeRedirectUri = (index: number) => {
    if (formData.redirect_uris.length > 1) {
      formData.redirect_uris.splice(index, 1)
    } else {
      ElMessage.warning('至少保留一个回调地址')
    }
  }

  /**
   * 提交表单
   */
  const submitForm = async (): Promise<boolean> => {
    loading.value = true
    try {
      // 过滤空的回调地址
      const filteredRedirectUris = formData.redirect_uris.filter(uri => uri.trim() !== '')
      
      if (filteredRedirectUris.length === 0) {
        ElMessage.error('至少需要一个有效的回调地址')
        return false
      }

      const requestData: CreateOAuthClientRequest = {
        ...formData,
        redirect_uris: filteredRedirectUris
      }

      const response = await createOAuthClient(requestData)
      if (response.success) {
        ElMessage.success('创建 OAuth 客户端成功')
        return true
      } else {
        ElMessage.error(response.message || '创建 OAuth 客户端失败')
        return false
      }
    } catch (error: any) {
      ElMessage.error(error.message || '创建 OAuth 客户端失败')
      return false
    } finally {
      loading.value = false
    }
  }

  /**
   * 重置表单
   */
  const resetForm = () => {
    formData.name = ''
    formData.client_secret = generateClientSecret()
    formData.description = ''
    formData.logo = ''
    formData.redirect_uris = ['']
    formData.grant_types = []
    formData.scopes = ['profile']
    formData.status = 1
  }

  return {
    loading,
    formData,
    generateClientSecret,
    initClientSecret,
    regenerateClientSecret,
    addRedirectUri,
    removeRedirectUri,
    submitForm,
    resetForm
  }
}

