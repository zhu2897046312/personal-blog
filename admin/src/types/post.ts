import type { User } from './user'
import type { Category } from './category'
import type { Tag } from './tag'
import { SearchRequest } from './common'

// 文章实体
export interface Post {
  id: number
  title: string
  content: string
  user_id: number
  category_id: number
  status: number // 1:公开 2:草稿
  view_count: number
  created_at: string
  updated_at: string
  user?: User
  category?: Category
  tags?: Tag[]
}

// 创建文章请求
export interface CreatePostRequest {
  title: string
  content: string
  category_id: number
  tags?: string[]
  status: number // 1:公开 2:草稿
}

// 更新文章请求
export interface UpdatePostRequest {
  title: string
  content: string
  category_id: number
  tags?: string[]
  status: number // 1:公开 2:草稿
}

// 文章列表请求
export interface ListPostsRequest extends SearchRequest {
  category_id?: number
  tag?: string
  status?: number // 1:公开 2:草稿
}

// 更新文章状态请求
export interface UpdatePostStatusRequest {
  status: number // 1:公开 2:草稿
}
