import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Comment, CreateCommentRequest, ListCommentsRequest, UpdateCommentStatusRequest } from '@/types/comment'
import { createComment, deleteComment, getComments, updateCommentStatus } from '@/api/comments'
import { ElMessage } from 'element-plus'

export const useCommentStore = defineStore('comment', () => {
  // 状态
  const comments = ref<Comment[]>([])
  const total = ref(0)
  const loading = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(10)

  // 获取评论列表
  const fetchComments = async (params: ListCommentsRequest) => {
    loading.value = true
    try {
      const { data } = await getComments(params)
      comments.value = data.data.items
      total.value = data.data.total
      currentPage.value = params.page
      pageSize.value = params.page_size
      return data
    } catch (error) {
      ElMessage.error('获取评论列表失败')
      return null
    } finally {
      loading.value = false
    }
  }

  // 创建评论
  const addComment = async (data: CreateCommentRequest) => {
    try {
      const { data: newComment } = await createComment(data)
      await fetchComments({ 
        post_id: data.post_id,
        page: currentPage.value,
        page_size: pageSize.value
      })
      ElMessage.success('评论成功')
      return newComment
    } catch (error) {
      ElMessage.error('评论失败')
      return null
    }
  }

  // 更新评论状态
  const changeCommentStatus = async (id: number, data: UpdateCommentStatusRequest) => {
    try {
      await updateCommentStatus(id, data)
      const index = comments.value.findIndex(comment => comment.id === id)
      if (index !== -1) {
        comments.value[index].status = data.status
      }
      ElMessage.success('更新状态成功')
      return true
    } catch (error) {
      ElMessage.error('更新状态失败')
      return false
    }
  }

  // 删除评论
  const removeComment = async (id: number) => {
    try {
      await deleteComment(id)
      const index = comments.value.findIndex(comment => comment.id === id)
      if (index !== -1) {
        comments.value.splice(index, 1)
        total.value--
      }
      ElMessage.success('删除成功')
      return true
    } catch (error) {
      ElMessage.error('删除失败')
      return false
    }
  }

  // 清空评论列表
  const clearComments = () => {
    comments.value = []
    total.value = 0
    currentPage.value = 1
    pageSize.value = 10
  }

  return {
    // 状态
    comments,
    total,
    loading,
    currentPage,
    pageSize,
    // 方法
    fetchComments,
    addComment,
    changeCommentStatus,
    removeComment,
    clearComments,
  }
})
