import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { createOAuthClient } from '@/api/oauth_client'
import type { CreateOAuthClientRequest } from '@/types/oauth_client'

// 配置字段默认值（与后端保持一致，单位：秒）
export const DEFAULT_AUTH_CODE_EXPIRE = 300       // 授权码过期时间：5分钟
export const DEFAULT_ACCESS_TOKEN_EXPIRE = 3600   // 访问令牌过期时间：1小时
export const DEFAULT_REFRESH_TOKEN_EXPIRE = 2592000 // 刷新令牌过期时间：30天
export const DEFAULT_SUBJECT_LENGTH = 16          // 用户标识长度
export const DEFAULT_SUBJECT_PREFIX = 'usr_'      // 用户标识前缀

/**
 * OAuth 客户端表单管理相关的组合式函数
 */
export function useOAuthClientForm() {
  const loading = ref(false)

  const formData = reactive<CreateOAuthClientRequest & {
    // 编辑模式下用于轮换的密钥字段（可选）
    new_client_secret?: string
    new_access_token_secret?: string
    new_refresh_token_secret?: string
    new_subject_secret?: string
  }>({
    // 必填密钥字段
    client_secret: '',
    access_token_secret: '',
    refresh_token_secret: '',
    subject_secret: '',

    // 必填基本字段
    name: '',
    redirect_uris: [''],
    grant_types: [],
    scopes: ['profile'],
    status: 1,

    // 可选基本字段
    description: '',
    logo: '',

    // 可选配置字段（带默认值）
    auth_code_expire: DEFAULT_AUTH_CODE_EXPIRE,
    access_token_expire: DEFAULT_ACCESS_TOKEN_EXPIRE,
    refresh_token_expire: DEFAULT_REFRESH_TOKEN_EXPIRE,
    subject_length: DEFAULT_SUBJECT_LENGTH,
    subject_prefix: DEFAULT_SUBJECT_PREFIX
  })

  /**
   * 生成随机密钥（32位字符串）
   */
  const generateSecret = (): string => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    let result = ''
    for (let i = 0; i < 32; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    return result
  }

  /**
   * 初始化所有密钥字段（创建模式）
   */
  const initSecrets = () => {
    formData.client_secret = generateSecret()
    formData.access_token_secret = generateSecret()
    formData.refresh_token_secret = generateSecret()
    formData.subject_secret = generateSecret()
  }

  /**
   * 初始化 client_secret（兼容旧版）
   */
  const initClientSecret = () => {
    initSecrets()
  }

  /**
   * 重新生成指定密钥
   */
  const regenerateSecret = (field: 'client_secret' | 'access_token_secret' | 'refresh_token_secret' | 'subject_secret') => {
    formData[field] = generateSecret()
    ElMessage.success('密钥已重新生成')
  }

  /**
   * 重新生成 client_secret（兼容旧版）
   */
  const regenerateClientSecret = () => {
    regenerateSecret('client_secret')
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
        client_secret: formData.client_secret,
        access_token_secret: formData.access_token_secret,
        refresh_token_secret: formData.refresh_token_secret,
        subject_secret: formData.subject_secret,
        name: formData.name,
        redirect_uris: filteredRedirectUris,
        grant_types: formData.grant_types,
        scopes: formData.scopes,
        status: formData.status,
        description: formData.description,
        logo: formData.logo,
        auth_code_expire: formData.auth_code_expire,
        access_token_expire: formData.access_token_expire,
        refresh_token_expire: formData.refresh_token_expire,
        subject_length: formData.subject_length,
        subject_prefix: formData.subject_prefix
      }

      await createOAuthClient(requestData)
      ElMessage.success('创建 OAuth 客户端成功')
      return true
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
    formData.client_secret = generateSecret()
    formData.access_token_secret = generateSecret()
    formData.refresh_token_secret = generateSecret()
    formData.subject_secret = generateSecret()
    formData.description = ''
    formData.logo = ''
    formData.redirect_uris = ['']
    formData.grant_types = []
    formData.scopes = ['profile']
    formData.status = 1
    formData.auth_code_expire = DEFAULT_AUTH_CODE_EXPIRE
    formData.access_token_expire = DEFAULT_ACCESS_TOKEN_EXPIRE
    formData.refresh_token_expire = DEFAULT_REFRESH_TOKEN_EXPIRE
    formData.subject_length = DEFAULT_SUBJECT_LENGTH
    formData.subject_prefix = DEFAULT_SUBJECT_PREFIX
  }

  return {
    loading,
    formData,
    generateSecret,
    initSecrets,
    initClientSecret,
    regenerateSecret,
    regenerateClientSecret,
    addRedirectUri,
    removeRedirectUri,
    submitForm,
    resetForm,
    // 导出默认值供组件使用
    DEFAULT_AUTH_CODE_EXPIRE,
    DEFAULT_ACCESS_TOKEN_EXPIRE,
    DEFAULT_REFRESH_TOKEN_EXPIRE,
    DEFAULT_SUBJECT_LENGTH,
    DEFAULT_SUBJECT_PREFIX
  }
}
