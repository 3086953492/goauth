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
  name: string
  client_secret: string
  description?: string
  logo?: string
  redirect_uris: string[]
  grant_types: string[]
  scopes: string[]
  status: number
}

export interface UpdateOAuthClientRequest {
  name: string
  description?: string
  logo?: string
  redirect_uris?: string[]
  grant_types?: string[]
  scopes?: string[]
  status?: number
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
  created_at: string
  updated_at: string
}

export type OAuthClientFormMode = 'create' | 'edit' | 'view'

