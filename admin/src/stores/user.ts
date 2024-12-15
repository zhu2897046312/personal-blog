import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User, ListUsersRequest} from '@/types'
import { getProfile, login, updateProfile, register as registerApi, getUsers } from '@/api/users'
import { ElMessage } from 'element-plus'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const users = ref<User[]>([])
  const token = ref<string | null>(null)
  const loading = ref(false)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 初始化用户状态
  const init = async () => {
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      token.value = storedToken
      try {
        const { data : response } = await getProfile()
        user.value = response.data
      } catch (error) {
        logout()
      }
    }
  }

  // 获取所有用户列表
  const fetchUsers = async (params : ListUsersRequest) => {
    loading.value = true
    try {
      const { data } = await getUsers(params)
      users.value = data.data.items
      total.value = data.data.total
      currentPage.value = params.page
      pageSize.value = params.page_size
      return data
    } catch (error) {
      ElMessage.error('获取分类列表失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 用户登录
  const userLogin = async (username: string, password: string) => {
    try {
      const { data } = await login({ username, password })
      token.value = data.token
      user.value = data.user
      localStorage.setItem('token', data.token)
      ElMessage.success('登录成功')
      return true
    } catch (error) {
      return false
    }
  }

  // 注册
  const register = async (username: string, password: string, email: string, nickname: string) => {
    try {
      const { data } = await registerApi({ username, password, email, nickname })
      ElMessage.success('注册成功')
      return true
    } catch (error) {
      return false
    }
  }

  // 更新用户信息
  const updateUser = async (data: Partial<User>) => {
    try {
      const { data: updatedUser } = await updateProfile(data)
      user.value = updatedUser.data
      ElMessage.success('更新成功')
      return true
    } catch (error) {
      return false
    }
  }

  // // 更新用户状态
  // const updateUserStatus = async (id: number, status: number) => {
  //   try {
  //     await updateProfile({ id, status })
  //     // Update the local users array to reflect the status change
  //     const userIndex = users.value.findIndex(u => u.id === id)
  //     if (userIndex !== -1) {
  //       users.value[userIndex].status = status
  //     }
  //     ElMessage.success('用户状态更新成功')
  //     return true
  //   } catch (error) {
  //     ElMessage.error('更新用户状态失败')
  //     return false
  //   }
  // }

  // 退出登录
  const logout = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  return {
    user,
    token,
    init,
    fetchUsers,
    userLogin,
    register,
    updateUser,
    //updateUserStatus,
    logout,
    loading,
    users,
    total,
    currentPage,
    pageSize
  }
})
