import { SearchRequest, StatusRequest } from './common'

// 评论实体接口
export interface Comment {
  id: number
  post_id: number
  content: string
  parent_id?: number
  created_at: string
  updated_at: string
  status: number
}

// 创建评论请求
export interface CreateCommentRequest {
  post_id: number
  content: string
  parent_id?: number
}

// 评论列表请求
export interface ListCommentsRequest extends SearchRequest {
  post_id: number
}

// 更新评论状态请求
export interface UpdateCommentStatusRequest extends StatusRequest {}