import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Tag, CreateTagRequest, UpdateTagRequest, CreateTagsRequest } from '@/types/tag'
import { createTag, createTags, deleteTag, getPostTags, getTags, updateTag } from '@/api/tags'
import { ElMessage } from 'element-plus'

export const useTagStore = defineStore('tag', () => {
  // 状态
  const tags = ref<Tag[]>([])
  const postTags = ref<Tag[]>([])
  const loading = ref(false)

  // 获取标签列表
  const fetchTags = async () => {
    loading.value = true
    try {
      const { data: response } = await getTags()
      tags.value = response.data
      return response
    } catch (error) {
      ElMessage.error('获取标签列表失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 获取文章标签
  const fetchPostTags = async (postId: number) => {
    loading.value = true
    try {
      const { data: response } = await getPostTags(postId)
      postTags.value = response.data.items
      return response
    } catch (error) {
      ElMessage.error('获取文章标签失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 创建标签
  const addTag = async (data: CreateTagRequest) => {
    try {
      const { data: response } = await createTag(data)
      const newTag = response.data
      tags.value.push(newTag)
      ElMessage.success('创建成功')
      return newTag
    } catch (error) {
      ElMessage.error('创建失败')
      return null
    }
  }

  // 批量创建标签
  const addTags = async (data: CreateTagsRequest) => {
    try {
      const { data: response } = await createTags(data)
      const newTags = response.data
      tags.value.push(...newTags)
      ElMessage.success('批量创建成功')
      return newTags
    } catch (error) {
      ElMessage.error('批量创建失败')
      return null
    }
  }

  // 更新标签
  const editTag = async (id: number, data: UpdateTagRequest) => {
    try {
      const { data: response } = await updateTag(id, data)
      const updatedTag = response.data
      const index = tags.value.findIndex(tag => tag.id === id)
      if (index !== -1) {
        tags.value[index] = updatedTag
      }
      // 同时更新文章标签列表
      const postTagIndex = postTags.value.findIndex(tag => tag.id === id)
      if (postTagIndex !== -1) {
        postTags.value[postTagIndex] = updatedTag
      }
      ElMessage.success('更新成功')
      return updatedTag
    } catch (error) {
      ElMessage.error('更新失败')
      return null
    }
  }

  // 删除标签
  const removeTag = async (id: number) => {
    try {
      await deleteTag(id)
      // 从标签列表中移除
      const index = tags.value.findIndex(tag => tag.id === id)
      if (index !== -1) {
        tags.value.splice(index, 1)
      }
      // 从文章标签列表中移除
      const postTagIndex = postTags.value.findIndex(tag => tag.id === id)
      if (postTagIndex !== -1) {
        postTags.value.splice(postTagIndex, 1)
      }
      ElMessage.success('删除成功')
      return true
    } catch (error) {
      ElMessage.error('删除失败')
      return false
    }
  }

  // 清空文章标签
  const clearPostTags = () => {
    postTags.value = []
  }

  return {
    // 状态
    tags,
    postTags,
    loading,
    // 方法
    fetchTags,
    fetchPostTags,
    addTag,
    addTags,
    editTag,
    removeTag,
    clearPostTags,
  }
})