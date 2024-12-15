import request from './request'
import type { Post, CreatePostRequest, UpdatePostRequest, ListPostsRequest, UpdatePostStatusRequest } from '@/types/post'
import type { Response, PaginationResponseData } from '@/types/common'

// 创建文章
export function createPost(data: CreatePostRequest) {
  return request.post<Response<Post>>('/posts', data)
}

// 更新文章
export function updatePost(id: number, data: UpdatePostRequest) {
  return request.put<Response<Post>>(`/posts/${id}`, data)
}

// 删除文章
export function deletePost(id: number) {
  return request.delete<Response<null>>(`/posts/${id}`)
}

// 获取单篇文章
export function getPost(id: number) {
  return request.get<Response<Post>>(`/posts/${id}`)
}

// 获取文章列表
export function getPosts(params: ListPostsRequest) {
  return request.get<Response<PaginationResponseData<Post>>>('/posts', { params })
}

// 更新文章状态
export function updatePostStatus(id: number, data: UpdatePostStatusRequest) {
  return request.patch<Response<null>>(`/posts/${id}/status`, data)
}
