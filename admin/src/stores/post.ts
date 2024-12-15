import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Post, CreatePostRequest, UpdatePostRequest, ListPostsRequest, UpdatePostStatusRequest } from '@/types/post'
import { createPost, deletePost, getPost, getPosts, updatePost, updatePostStatus } from '@/api/posts'
import { ElMessage } from 'element-plus'

export const usePostStore = defineStore('post', () => {
  // 状态
  const posts = ref<Post[]>([])
  const total = ref(0)
  const currentPost = ref<Post | null>(null)
  const loading = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 获取文章列表
  const fetchPosts = async (params: ListPostsRequest) => {
    loading.value = true
    try {
      const { data } = await getPosts(params)
      posts.value = data.data.items
      total.value = data.data.total
      currentPage.value = params.page
      pageSize.value = params.page_size
      return data
    } catch (error) {
      ElMessage.error('获取文章列表失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 获取文章详情
  const fetchPost = async (id: number) => {
    loading.value = true
    try {
      const { data } = await getPost(id)
      currentPost.value = data.data
      return data
    } catch (error) {
      ElMessage.error('获取文章详情失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 创建文章
  const addPost = async (data: CreatePostRequest) => {
    try {
      const { data: newPost } = await createPost(data)
      await fetchPosts({ 
        page: currentPage.value,
        page_size: pageSize.value
      })
      ElMessage.success('创建成功')
      return newPost
    } catch (error) {
      ElMessage.error('创建失败')
      return null
    }
  }

  // 更新文章
  const editPost = async (id: number, postData: UpdatePostRequest) => {
    try {
      const { data : response} = await updatePost(id, postData)
      const updatedPost = response.data // 访问嵌套的updatedPost

      const index = posts.value.findIndex(post => post.id === id)
      if (index !== -1) {
        posts.value[index] = updatedPost
      }
      if (currentPost.value?.id === id) {
        currentPost.value = updatedPost
      }
      ElMessage.success('更新成功')
      return updatedPost
    } catch (error) {
      ElMessage.error('更新失败')
      return null
    }
  }

  // 更新文章状态
  const changePostStatus = async (id: number, data: UpdatePostStatusRequest) => {
    try {
      await updatePostStatus(id, data)
      const index = posts.value.findIndex(post => post.id === id)
      if (index !== -1) {
        posts.value[index].status = data.status
      }
      if (currentPost.value?.id === id) {
        currentPost.value.status = data.status
      }
      ElMessage.success('更新状态成功')
      return true
    } catch (error) {
      ElMessage.error('更新状态失败')
      return false
    }
  }

  // 删除文章
  const removePost = async (id: number) => {
    try {
      await deletePost(id)
      const index = posts.value.findIndex(post => post.id === id)
      if (index !== -1) {
        posts.value.splice(index, 1)
        total.value--
      }
      if (currentPost.value?.id === id) {
        currentPost.value = null
      }
      ElMessage.success('删除成功')
      return true
    } catch (error) {
      ElMessage.error('删除失败')
      return false
    }
  }

  // 清空当前文章
  const clearCurrentPost = () => {
    currentPost.value = null
  }

  return {
    // 状态
    posts,
    total,
    currentPost,
    loading,
    currentPage,
    pageSize,
    // 方法
    fetchPosts,
    fetchPost,
    addPost,
    editPost,
    changePostStatus,
    removePost,
    clearCurrentPost,
  }
})
