<template>
  <div class="posts-container">
    <div class="header">
      <el-button type="primary" @click="handleAdd">新建文章</el-button>
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item>
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索文章标题"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="searchForm.category_id"
            placeholder="选择分类"
            clearable
          >
            <el-option
              v-for="item in categories"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select
            v-model="searchForm.status"
            placeholder="文章状态"
            clearable
          >
            <el-option label="公开" :value="1" />
            <el-option label="草稿" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table v-loading="loading" :data="posts" style="width: 100%">
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column prop="category.name" label="分类" width="120" />
      <el-table-column label="标签" min-width="200">
        <template #default="{ row }">
          <el-tag
            v-for="tag in row.tags"
            :key="tag.id"
            class="tag"
            size="small"
          >
            {{ tag.name }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '公开' : '草稿' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="view_count" label="浏览量" width="100" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑文章' : '新建文章'"
      width="70%"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="15"
          />
        </el-form-item>
        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类">
            <el-option
              v-for="item in categories"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="标签">
          <el-select
            v-model="form.tags"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择标签"
          >
            <el-option
              v-for="item in tags"
              :key="item.id"
              :label="item.name"
              :value="item.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">公开</el-radio>
            <el-radio :label="2">草稿</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessageBox } from 'element-plus'
import { usePostStore } from '@/stores'
import { useCategoryStore } from '@/stores'
import { useTagStore } from '@/stores'
import type { Post, CreatePostRequest, UpdatePostRequest, ListPostsRequest } from '@/types/post'
import type { ListTagsRequest } from '@/types/tag'

const postStore = usePostStore()
const categoryStore = useCategoryStore()
const tagStore = useTagStore()

const { posts, loading, total, currentPage, pageSize } = storeToRefs(postStore)
const { categories } = storeToRefs(categoryStore)
const { tags } = storeToRefs(tagStore)

const dialogVisible = ref(false)
const formRef = ref<FormInstance>()

const searchForm = reactive<ListPostsRequest>({
  keyword: '',
  category_id: undefined,
  status: undefined,
  page: 1,
  page_size: 10
})

const form = reactive<CreatePostRequest & { id?: number }>({
  id: undefined,
  title: '',
  content: '',
  category_id: 0,
  tags: [],
  status: 1
})

const rules: FormRules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 获取数据
const fetchData = async () => {
  await postStore.fetchPosts(searchForm)
}

// 获取数据
const fetchTagsData = async () => {
  const params: ListTagsRequest = {
    page: currentPage.value,
    page_size: pageSize.value
  }
  await tagStore.fetchTags(params)
}

// 重置表单
const resetForm = () => {
  form.id = undefined
  form.title = ''
  form.content = ''
  form.category_id = 0
  form.tags = []
  form.status = 1
}

// 重置搜索
const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.category_id = undefined
  searchForm.status = undefined
  searchForm.page = 1
  searchForm.page_size = 10
  fetchData()
}

// 搜索
const handleSearch = () => {
  searchForm.page = 1
  fetchData()
}

const handleAdd = () => {
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: Post) => {
  form.id = row.id
  form.title = row.title
  form.content = row.content
  form.category_id = row.category_id
  form.status = row.status
  form.tags = row.tags?.map(tag => tag.id.toString()) || []
  dialogVisible.value = true
}

const handleDelete = async (row: Post) => {
  try {
    await ElMessageBox.confirm('确认删除该文章吗？', '提示', {
      type: 'warning'
    })
    await postStore.removePost(row.id)
    await fetchData()
  } catch (error) {
    // 用户取消删除操作
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      const postData: CreatePostRequest | UpdatePostRequest = {
        title: form.title,
        content: form.content,
        category_id: form.category_id,
        tags: form.tags,
        status: form.status
      }

      if (form.id) {
        await postStore.editPost(form.id, postData)
      } else {
        await postStore.addPost(postData)
      }
      dialogVisible.value = false
      resetForm()
      await fetchData()
    }
  })
}

const handleSizeChange = (size: number) => {
  searchForm.page_size = size
  searchForm.page = 1
  fetchData()
}

const handleCurrentChange = (page: number) => {
  searchForm.page = page
  fetchData()
}

onMounted(async () => {
  await Promise.all([
    categoryStore.fetchCategories({ page: 1, page_size: 10 }),
    fetchTagsData(),
    fetchData()
  ])
})
</script>

<style scoped>
.posts-container {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.search-form {
  flex: 1;
  margin-left: 20px;
}

.tag {
  margin-right: 5px;
  margin-bottom: 5px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

:deep(.el-dialog__body) {
  padding-top: 20px;
}
</style>
