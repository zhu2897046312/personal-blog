<template>
  <div class="users-container">
    <el-table v-loading="loading" :data="users" style="width: 100%">
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column label="头像" width="100">
        <template #default="{ row }">
          <el-avatar :src="row.avatar" />
        </template>
      </el-table-column>
      <el-table-column prop="role" label="角色">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'info'">
            {{ row.role === 'admin' ? '管理员' : '普通用户' }}
          </el-tag>
        </template>
      </el-table-column>
      <!-- <el-table-column prop="created_at" label="注册时间" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-switch
            v-model="row.status"
            :active-value="1"
            :inactive-value="0"
            @change="handleStatusChange(row)"
          />
        </template>
      </el-table-column> -->
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import type { ListUsersRequest } from '@/types'
import { ElMessageBox } from 'element-plus'

const userStore = useUserStore()
const { users, loading, total } = storeToRefs(userStore)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)

// 获取用户列表
const fetchUsers = async () => {
  const params: ListUsersRequest = {
    page: currentPage.value,
    page_size: pageSize.value
  }
  await userStore.fetchUsers(params)
}

// 处理分页变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchUsers()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchUsers()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchUsers()
}

// const handleStatusChange = async (row: any) => {
//   try {
//     await userStore.updateUser({ id: row.id, status: row.status })
//   } catch (error) {
//     console.error('Failed to update user status', error)
//   }
// }

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.users-container {
  padding: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-avatar) {
  background-color: #f5f7fa;
}
</style>
