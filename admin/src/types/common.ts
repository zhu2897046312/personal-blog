export interface PaginationRequest {
  page: number
  page_size: number
}

export interface SearchRequest extends PaginationRequest {
  keyword?: string
}

export interface IDRequest {
  id: number
}

export interface StatusRequest {
  status: 0 | 1  // 0:禁用 1:启用
}

export interface SortRequest {
  order_by?: 'created_at' | 'updated_at'
  order?: 'asc' | 'desc'
}

export interface PaginationResponseData<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface Response<T> {
  code: number
  message: string
  data: T
}
