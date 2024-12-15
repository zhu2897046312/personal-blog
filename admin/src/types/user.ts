import { SearchRequest } from './common'

export interface User {
  id: number
  username: string
  nickname: string
  email: string
  avatar: string
  role: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface UpdateUserRequest {
  nickname?: string
  email?: string
  avatar?: string
  password?: string
}

// 用户列表请求
export interface ListUsersRequest extends SearchRequest { }
