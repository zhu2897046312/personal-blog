<template>
  <div class="tags-container">
    <div class="header">
      <el-button type="primary" @click="handleAdd">新建标签</el-button>
      <el-button @click="handleBatchAdd">批量添加</el-button>
    </div>

    <el-table v-loading="loading" :data="tags" style="width: 100%">
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="created_at" label="创建时间" />
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
      :title="form.id ? '编辑标签' : '新建标签'"
      width="30%"
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

    <!-- 批量添加对话框 -->
    <el-dialog
      v-model="batchDialogVisible"
      title="批量添加标签"
      width="50%"
    >
      <el-form
        ref="batchFormRef"
        :model="batchForm"
        :rules="batchRules"
        label-width="80px"
      >
        <el-form-item label="标签" prop="names">
          <el-input
            v-model="batchForm.names"
            type="textarea"
            :rows="4"
            placeholder="请输入标签名称，多个标签用逗号分隔"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="batchDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleBatchSubmit">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessageBox } from 'element-plus'
import { useTagStore } from '@/stores'
import type { Tag } from '@/types'
import { storeToRefs } from 'pinia'

const tagStore = useTagStore()
const { tags, loading } = storeToRefs(tagStore)

const dialogVisible = ref(false)
const batchDialogVisible = ref(false)
const formRef = ref<FormInstance>()
const batchFormRef = ref<FormInstance>()

const form = reactive({
  id: 0,
  name: ''
})

const batchForm = reactive({
  names: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' }
  ]
}

const batchRules: FormRules = {
  names: [
    { required: true, message: '请输入标签名称', trigger: 'blur' }
  ]
}

// 获取数据
const fetchData = async () => {
  await tagStore.fetchTags()
}

// 重置表单
const resetForm = () => {
  form.id = 0
  form.name = ''
}

// 重置批量添加表单
const resetBatchForm = () => {
  batchForm.names = ''
}

// 新增标签
const handleAdd = () => {
  resetForm()
  dialogVisible.value = true
}

// 批量添加
const handleBatchAdd = () => {
  resetBatchForm()
  batchDialogVisible.value = true
}

// 编辑标签
const handleEdit = (row: Tag) => {
  Object.assign(form, row)
  dialogVisible.value = true
}

// 删除标签
const handleDelete = (row: Tag) => {
  ElMessageBox.confirm(
    '确认删除该标签吗？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    const success = await tagStore.removeTag(row.id)
    if (success) {
      fetchData()
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      const data = {
        name: form.name
      }
      
      let success
      if (form.id) {
        success = await tagStore.editTag(form.id, data)
      } else {
        success = await tagStore.addTag(data)
      }
      
      if (success) {
        dialogVisible.value = false
        fetchData()
      }
    }
  })
}

// 提交批量添加表单
const handleBatchSubmit = async () => {
  if (!batchFormRef.value) return

  await batchFormRef.value.validate(async (valid) => {
    if (valid) {
      const names = batchForm.names.split(',').map(name => name.trim()).filter(Boolean)
      const success = await tagStore.addTags(names)
      
      if (success.length > 0) {
        batchDialogVisible.value = false
        fetchData()
      }
    }
  })
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.tags-container {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
