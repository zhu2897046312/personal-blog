import request from './request'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest, ListCategoriesRequest, UpdateCategoryStatusRequest } from '@/types/category'
import type { Response, PaginationResponseData } from '@/types/common'

// 创建分类
export function createCategory(data: CreateCategoryRequest) {
  return request.post<Response<Category>>('/categories', data)
}

// 更新分类
export function updateCategory(id: number, data: UpdateCategoryRequest) {
  return request.put<Response<Category>>(`/categories/${id}`, data)
}

// 删除分类
export function deleteCategory(id: number) {
  return request.delete<Response<null>>(`/categories/${id}`)
}

// 获取单个分类
export function getCategory(id: number) {
  return request.get<Response<Category>>(`/categories/${id}`)
}

// 获取分类列表
export function getCategories(params?: ListCategoriesRequest) {
  return request.get<Response<PaginationResponseData<Category>>>('/categories', { params })
}

// 更新分类状态
export function updateCategoryStatus(id: number, data: UpdateCategoryStatusRequest) {
  return request.patch<Response<null>>(`/categories/${id}/status`, data)
}
