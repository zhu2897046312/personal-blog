// 标签实体
export interface Tag {
  id: number
  name: string
  created_at: string
  updated_at: string
}

// 创建标签请求
export interface CreateTagRequest {
  name: string
}

// 更新标签请求
export interface UpdateTagRequest {
  name: string
}

// 批量创建标签请求
export interface CreateTagsRequest {
  names: string[]
}
