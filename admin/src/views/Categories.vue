<template>
  <div class="categories-container">
    <div class="header">
      <el-button type="primary" @click="handleAdd">新建分类</el-button>
    </div>

    <el-table v-loading="loading" :data="categories" style="width: 100%">
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="description" label="描述" />
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column prop="updated_at" label="更新时间" />
      <el-table-column label="操作" width="200">
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

    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑分类' : '新建分类'"
      width="50%"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
          />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { useCategoryStore } from '@/stores'
import type { Category, CreateCategoryRequest } from '@/types'

const categoryStore = useCategoryStore()
const { categories, loading } = storeToRefs(categoryStore)

const dialogVisible = ref(false)
const formRef = ref<FormInstance>()

interface CategoryForm extends CreateCategoryRequest {
  id?: number
}

const form = reactive<CategoryForm>({
  id: undefined,
  name: '',
  description: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入分类描述', trigger: 'blur' },
    { min: 2, max: 200, message: '长度在 2 到 200 个字符', trigger: 'blur' }
  ]
}

// 重置表单
const resetForm = () => {
  form.id = undefined
  form.name = ''
  form.description = ''
}

// 新增分类
const handleAdd = () => {
  resetForm()
  dialogVisible.value = true
}

// 编辑分类
const handleEdit = (row: Category) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

// 删除分类
const handleDelete = (row: Category) => {
  ElMessageBox.confirm(
    '确认删除该分类吗？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    const success = await categoryStore.removeCategory(row.id)
    if (success) {
      await categoryStore.fetchCategories()
      ElMessage.success('删除成功')
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      const data: CreateCategoryRequest = {
        name: form.name,
        description: form.description
      }
      
      let success
      if (form.id) {
        success = await categoryStore.editCategory(form.id, data)
      } else {
        success = await categoryStore.addCategory(data)
      }
      
      if (success) {
        dialogVisible.value = false
        await categoryStore.fetchCategories()
      }
    }
  })
}

onMounted(() => {
  categoryStore.fetchCategories()
})
</script>

<style scoped>
.categories-container {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
