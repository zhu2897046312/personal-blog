import request from './request'
import type { User, LoginRequest, UpdateUserRequest, PaginationResponseData, ListUsersRequest } from '@/types'
import type { Response} from '@/types/common'


export function login(data: LoginRequest) {
  return request.post<{ token: string; user: User }>('/users/login', data)
}

export function register(data: LoginRequest & { email: string; nickname: string }) {
  return request.post<Response<User>>('/users/register', data)
}

export function getProfile() {
  return request.get<Response<User>>('/users/profile')
}

export function updateProfile(data: UpdateUserRequest) {
  return request.put<Response<User>>('/users/profile', data)
}

export function changePassword(data: { old_password: string; new_password: string }) {
  return request.put<Response<null>>('/users/password', data)
}

export function getUsers(params: ListUsersRequest) {
  return request.get<Response<PaginationResponseData<User>>>('/users', { params })
}
