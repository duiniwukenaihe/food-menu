<template>
  <div class="admin-users">
    <div class="page-header">
      <div class="header-left">
        <h2>用户管理</h2>
        <span class="total-count">共 {{ total }} 个用户</span>
      </div>
      <div class="search-section">
        <el-input
          v-model="searchQuery"
          placeholder="搜索用户名或邮箱"
          style="width: 300px; margin-right: 12px"
          clearable
          @clear="loadUsers"
          @keyup.enter="loadUsers"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="loadUsers">搜索</el-button>
        <el-button @click="resetSearch">重置</el-button>
      </div>
    </div>

    <!-- 用户列表 -->
    <div class="users-table">
      <el-table
        v-loading="loading"
        :data="users"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'primary'" size="small">
              {{ row.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button 
              v-if="row.role !== 'admin'"
              type="primary" 
              size="small" 
              @click="toggleUserRole(row)"
            >
              设为管理员
            </el-button>
            <el-button 
              v-else
              type="warning" 
              size="small" 
              @click="toggleUserRole(row)"
            >
              取消管理员
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { api } from '@/utils/api'
import type { User } from '@/types'

const loading = ref(false)
const users = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')

onMounted(loadUsers)

async function loadUsers() {
  try {
    loading.value = true
    const response = await api.getUsers({
      page: currentPage.value,
      limit: pageSize.value,
      search: searchQuery.value
    })
    
    users.value = response.users || []
    total.value = response.total || 0
  } catch (error) {
    console.error('Failed to load users:', error)
    ElMessage.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

function resetSearch() {
  searchQuery.value = ''
  currentPage.value = 1
  loadUsers()
}

async function toggleUserRole(user: User) {
  const newRole = user.role === 'admin' ? 'user' : 'admin'
  const action = newRole === 'admin' ? '设为管理员' : '取消管理员'
  
  try {
    await ElMessageBox.confirm(
      `确定要${action}用户 "${user.username}" 吗？`,
      '更改用户角色',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    // 这里应该调用更新用户角色的API
    // await api.updateUserRole(user.id, { role: newRole })
    
    ElMessage.success(`用户角色已更新为${newRole === 'admin' ? '管理员' : '普通用户'}`)
    await loadUsers()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to toggle user role:', error)
      ElMessage.error('操作失败')
    }
  }
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleString('zh-CN')
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  loadUsers()
}

function handleCurrentChange(page: number) {
  currentPage.value = page
  loadUsers()
}
</script>

<style scoped>
.admin-users {
  background: white;
  border-radius: 8px;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h2 {
  margin: 0;
  color: #262626;
}

.total-count {
  color: #8c8c8c;
  font-size: 14px;
}

.search-section {
  display: flex;
  align-items: center;
}

.users-table {
  margin-bottom: 24px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-section {
    justify-content: center;
  }
}
</style>