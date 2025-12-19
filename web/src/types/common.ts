/**
 * 统一成功响应结构（ginx.OK 返回）
 */
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

/**
 * RFC7807 Problem Details（ginx.Fail 返回）
 */
export interface ProblemDetails {
  type: string
  title: string
  status: number
  detail: string
  extensions?: Record<string, unknown>
}

/**
 * 分页响应结构
 */
export interface PaginationResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}
