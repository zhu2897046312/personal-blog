import { SearchRequest, StatusRequest } from './common'

export interface Category {
  id: number
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface CreateCategoryRequest {
  name: string
  description?: string
}

export interface UpdateCategoryRequest extends CreateCategoryRequest {}

export interface ListCategoriesRequest extends SearchRequest {}

export interface UpdateCategoryStatusRequest extends StatusRequest {}
