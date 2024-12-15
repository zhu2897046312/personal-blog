import request from './request'
import type { Tag, CreateTagRequest, UpdateTagRequest, CreateTagsRequest } from '@/types/tag'
import type { Response, PaginationResponseData } from '@/types/common'

// 创建标签
export function createTag(data: CreateTagRequest) {
  return request.post<Response<Tag>>('/tags', data)
}

// 批量创建标签
export function createTags(data: CreateTagsRequest) {
  return request.post<Response<Tag[]>>('/tags/batch', data)
}

// 更新标签
export function updateTag(id: number, data: UpdateTagRequest) {
  return request.put<Response<Tag>>(`/tags/${id}`, data)
}

// 删除标签
export function deleteTag(id: number) {
  return request.delete<Response<null>>(`/tags/${id}`)
}

// 获取单个标签
export function getTag(id: number) {
  return request.get<Response<Tag>>(`/tags/${id}`)
}

// 获取标签列表
export function getTags() {
  return request.get<Response<Tag[]>>('/tags')
}

// 获取文章的标签
export function getPostTags(postId: number) {
  return request.get<Response<PaginationResponseData<Tag>>>(`/posts/${postId}/tags`)
}
