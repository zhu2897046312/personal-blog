import request from './request'
import type { User, LoginRequest, UpdateUserRequest, PaginationResponseData, ListUsersRequest } from '@/types'
import type { Response} from '@/types/common'

// 用户登录
export function login(data: LoginRequest) {
  return request.post<{ token: string; user: User }>('/users/login', data)
}

// 用户注册
export function register(data: LoginRequest & { email: string; nickname: string }) {
  return request.post<Response<User>>('/users/register', data)
}

// 获取用户信息
export function getProfile() {
  return request.get<Response<User>>('/users/profile')
}

// 更新用户信息
export function updateProfile(data: UpdateUserRequest) {
  return request.put<Response<User>>('/users/profile', data)
}

// 更新用户密码
export function changePassword(data: { old_password: string; new_password: string }) {
  return request.put<Response<null>>('/users/password', data)
}

// 获取用户列表
export function getUsers(params: ListUsersRequest) {
  return request.get<Response<PaginationResponseData<User>>>('/users', { params })
}

// 删除用户
export function deleteUser(id: number) {
  return request.delete<Response<null>>(`/users/${id}`)
}
