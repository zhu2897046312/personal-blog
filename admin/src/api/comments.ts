import request from './request'
import type { Comment, CreateCommentRequest, ListCommentsRequest, UpdateCommentStatusRequest } from '@/types/comment'
import type { Response, PaginationResponseData } from '@/types/common'

// 创建评论
export function createComment(data: CreateCommentRequest) {
  return request.post<Response<Comment>>('/comments', data)
}

// 获取评论列表
export function getComments(params: ListCommentsRequest) {
  return request.get<Response<PaginationResponseData<Comment>>>('/comments', { params })
}

// 删除评论
export function deleteComment(id: number) {
  return request.delete<Response<null>>(`/comments/${id}`)
}

// 更新评论状态
export function updateCommentStatus(id: number, data: UpdateCommentStatusRequest) {
  return request.patch<Response<null>>(`/comments/${id}/status`, data)
}
