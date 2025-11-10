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

