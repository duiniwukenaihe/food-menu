<template>
  <div class="admin-dishes">
    <div class="page-header">
      <div class="header-left">
        <h2>菜品管理</h2>
        <span class="total-count">共 {{ dishStore.total }} 道菜品</span>
      </div>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加菜品
      </el-button>
    </div>

    <!-- 搜索和筛选 -->
    <div class="filters-section">
      <el-row :gutter="16">
        <el-col :span="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索菜品名称"
            clearable
            @clear="loadDishes"
            @keyup.enter="loadDishes"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="selectedCategory" placeholder="选择分类" clearable>
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="selectedStatus" placeholder="状态" clearable>
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="loadDishes">搜索</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 菜品列表 -->
    <div class="dishes-table">
      <el-table
        v-loading="loading"
        :data="dishes"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="图片" width="100">
          <template #default="{ row }">
            <el-image
              :src="row.image_url || '/placeholder-food.jpg'"
              :alt="row.name"
              style="width: 60px; height: 60px; object-fit: cover; border-radius: 4px;"
              fit="cover"
            />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="菜品名称" min-width="150" />
        <el-table-column prop="category.name" label="分类" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.category?.name || '--' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">
            <span class="price">¥{{ row.price }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
              {{ row.is_active ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="应季" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_seasonal" type="warning" size="small">应季</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="editDish(row)">
              编辑
            </el-button>
            <el-button 
              :type="row.is_active ? 'warning' : 'success'" 
              size="small" 
              @click="toggleStatus(row)"
            >
              {{ row.is_active ? '禁用' : '启用' }}
            </el-button>
            <el-button type="danger" size="small" @click="deleteDish(row)">
              删除
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

    <!-- 创建/编辑对话框 -->
    <DishDialog
      v-model="showCreateDialog"
      :dish="editingDish"
      :categories="categories"
      @success="handleDialogSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { useDishStore } from '@/stores'
import { api } from '@/utils/api'
import DishDialog from '@/components/DishDialog.vue'
import type { Dish, Category } from '@/types'

const dishStore = useDishStore()

const loading = ref(false)
const dishes = ref<Dish[]>([])
const categories = ref<Category[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchQuery = ref('')
const selectedCategory = ref<number>()
const selectedStatus = ref<boolean>()

const showCreateDialog = ref(false)
const editingDish = ref<Dish | null>(null)

onMounted(async () => {
  await Promise.all([
    loadCategories(),
    loadDishes()
  ])
})

async function loadCategories() {
  try {
    categories.value = await api.getCategories()
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

async function loadDishes() {
  try {
    loading.value = true
    const response = await api.getDishes({
      page: currentPage.value,
      limit: pageSize.value,
      category_id: selectedCategory.value,
      search: searchQuery.value
    })
    
    dishes.value = response.dishes || []
    total.value = response.total || 0
    
    // 过滤状态
    if (selectedStatus.value !== undefined) {
      dishes.value = dishes.value.filter(dish => dish.is_active === selectedStatus.value)
      total.value = dishes.value.length
    }
  } catch (error) {
    console.error('Failed to load dishes:', error)
    ElMessage.error('加载菜品列表失败')
  } finally {
    loading.value = false
  }
}

function resetFilters() {
  searchQuery.value = ''
  selectedCategory.value = undefined
  selectedStatus.value = undefined
  currentPage.value = 1
  loadDishes()
}

function editDish(dish: Dish) {
  editingDish.value = { ...dish }
  showCreateDialog.value = true
}

async function toggleStatus(dish: Dish) {
  try {
    await api.updateDish(dish.id, { is_active: !dish.is_active })
    ElMessage.success(`菜品已${dish.is_active ? '禁用' : '启用'}`)
    await loadDishes()
  } catch (error) {
    console.error('Failed to toggle dish status:', error)
    ElMessage.error('操作失败')
  }
}

async function deleteDish(dish: Dish) {
  try {
    await ElMessageBox.confirm(
      `确定要删除菜品 "${dish.name}" 吗？此操作不可恢复。`,
      '删除菜品',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await api.deleteDish(dish.id)
    ElMessage.success('菜品删除成功')
    await loadDishes()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete dish:', error)
      ElMessage.error('删除失败')
    }
  }
}

function handleDialogSuccess() {
  showCreateDialog.value = false
  editingDish.value = null
  loadDishes()
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleString('zh-CN')
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  loadDishes()
}

function handleCurrentChange(page: number) {
  currentPage.value = page
  loadDishes()
}
</script>

<style scoped>
.admin-dishes {
  background: white;
  border-radius: 8px;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
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

.filters-section {
  margin-bottom: 24px;
  padding: 16px;
  background: #fafafa;
  border-radius: 6px;
}

.dishes-table {
  margin-bottom: 24px;
}

.price {
  color: #f56c6c;
  font-weight: bold;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
</style>