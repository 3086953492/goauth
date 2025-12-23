export interface OAuthClient {
  id: number
  name: string
  client_id: string
  client_secret: string
  redirect_uris: string[]
  logo: string
  status: number
  created_at: string
  updated_at: string
}

export interface OAuthClientListResponse {
  id: number
  name: string
  logo: string
  status: number
}

export interface CreateOAuthClientRequest {
  // 必填密钥字段
  client_secret: string
  access_token_secret: string
  refresh_token_secret: string
  subject_secret: string

  // 必填基本字段
  name: string
  redirect_uris: string[]
  grant_types: string[]
  scopes: string[]
  status: number

  // 可选基本字段
  description?: string
  logo?: string

  // 可选配置字段（单位：秒）
  auth_code_expire?: number
  access_token_expire?: number
  refresh_token_expire?: number
  subject_length?: number
  subject_prefix?: string
}

export interface UpdateOAuthClientRequest {
  // 必填基本字段
  name: string

  // 可选基本字段
  description?: string
  logo?: string
  redirect_uris?: string[]
  grant_types?: string[]
  scopes?: string[]
  status?: number

  // 可选密钥字段（用于轮换，不传则不更新）
  client_secret?: string
  access_token_secret?: string
  refresh_token_secret?: string
  subject_secret?: string

  // 可选配置字段（单位：秒）
  auth_code_expire?: number
  access_token_expire?: number
  refresh_token_expire?: number
  subject_length?: number
  subject_prefix?: string
}

export interface OAuthClientDetailResponse {
  id: number
  name: string
  description: string
  logo: string
  redirect_uris: string[]
  grant_types: string[]
  scopes: string[]
  status: number

  // 配置字段（不返回密钥，单位：秒）
  auth_code_expire: number
  access_token_expire: number
  refresh_token_expire: number
  subject_length: number
  subject_prefix: string

  created_at: string
  updated_at: string
}

export type OAuthClientFormMode = 'create' | 'edit' | 'view'
