import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest, ListCategoriesRequest, UpdateCategoryStatusRequest } from '@/types/category'
import { createCategory, deleteCategory, getCategories, updateCategory, updateCategoryStatus } from '@/api/categories'
import { ElMessage } from 'element-plus'

export const useCategoryStore = defineStore('category', () => {
  // 状态
  const categories = ref<Category[]>([])
  const total = ref(0)
  const loading = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 获取分类列表
  const fetchCategories = async (params?: ListCategoriesRequest) => {
    loading.value = true
    try {
      const { data } = await getCategories(params)
      categories.value = data.data.items
      total.value = data.data.total
      currentPage.value = params?.page || 1
      pageSize.value = params?.page_size || 10
      return data
    } catch (error) {
      ElMessage.error('获取分类列表失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 创建分类
  const addCategory = async (data: CreateCategoryRequest) => {
    try {
      const { data: newCategory } = await createCategory(data)
      await fetchCategories({ page: currentPage.value, page_size: pageSize.value })
      ElMessage.success('创建成功')
      return newCategory
    } catch (error) {
      ElMessage.error('创建失败')
      return null
    }
  }

  // 更新分类
  const editCategory = async (id: number, categoryData: UpdateCategoryRequest) => {
    try {
      const { data : response } = await updateCategory(id, categoryData)
      const updatedCategory = response.data // 访问嵌套的分类数据

      const index = categories.value.findIndex(cat => cat.id === id)
      if (index !== -1) {
        categories.value[index] = updatedCategory
      }
      ElMessage.success('更新成功')
      return updatedCategory
    } catch (error) {
      ElMessage.error('更新失败')
      return null
    }
  }

  // // 更新分类状态
  // const changeCategoryStatus = async (id: number, data: UpdateCategoryStatusRequest) => {
  //   try {
  //     await updateCategoryStatus(id, data)
  //     const index = categories.value.findIndex(cat => cat.id === id)
  //     if (index !== -1) {
  //       categories.value[index].status = data.status
  //     }
  //     ElMessage.success('更新状态成功')
  //     return true
  //   } catch (error) {
  //     ElMessage.error('更新状态失败')
  //     return false
  //   }
  // }

  // 删除分类
  const removeCategory = async (id: number) => {
    try {
      await deleteCategory(id)
      const index = categories.value.findIndex(cat => cat.id === id)
      if (index !== -1) {
        categories.value.splice(index, 1)
      }
      ElMessage.success('删除成功')
      return true
    } catch (error) {
      ElMessage.error('删除失败')
      return false
    }
  }

  return {
    // 状态
    categories,
    total,
    loading,
    currentPage,
    pageSize,
    // 方法
    fetchCategories,
    addCategory,
    editCategory,
    //changeCategoryStatus,
    removeCategory,
  }
})
