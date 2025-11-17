<template>
  <div class="admin-categories">
    <div class="page-header">
      <div class="header-left">
        <h2>分类管理</h2>
        <span class="total-count">共 {{ categories.length }} 个分类</span>
      </div>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加分类
      </el-button>
    </div>

    <!-- 分类列表 -->
    <div class="categories-grid">
      <div
        v-for="category in categories"
        :key="category.id"
        class="category-card"
      >
        <div class="category-info">
          <h3>{{ category.name }}</h3>
          <p>{{ category.description || '暂无描述' }}</p>
          <div class="category-meta">
            <span>创建时间：{{ formatDate(category.created_at) }}</span>
          </div>
        </div>
        <div class="category-actions">
          <el-button type="primary" size="small" @click="editCategory(category)">
            编辑
          </el-button>
          <el-button type="danger" size="small" @click="deleteCategory(category)">
            删除
          </el-button>
        </div>
      </div>
    </div>

    <!-- 创建/编辑对话框 -->
    <CategoryDialog
      v-model="showCreateDialog"
      :category="editingCategory"
      @success="handleDialogSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { api } from '@/utils/api'
import CategoryDialog from '@/components/CategoryDialog.vue'
import type { Category } from '@/types'

const categories = ref<Category[]>([])
const showCreateDialog = ref(false)
const editingCategory = ref<Category | null>(null)

onMounted(loadCategories)

async function loadCategories() {
  try {
    categories.value = await api.getCategories()
  } catch (error) {
    console.error('Failed to load categories:', error)
    ElMessage.error('加载分类列表失败')
  }
}

function editCategory(category: Category) {
  editingCategory.value = { ...category }
  showCreateDialog.value = true
}

async function deleteCategory(category: Category) {
  try {
    await ElMessageBox.confirm(
      `确定要删除分类 "${category.name}" 吗？此操作不可恢复。`,
      '删除分类',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await api.deleteCategory(category.id)
    ElMessage.success('分类删除成功')
    await loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete category:', error)
      ElMessage.error('删除失败')
    }
  }
}

function handleDialogSuccess() {
  showCreateDialog.value = false
  editingCategory.value = null
  loadCategories()
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.admin-categories {
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

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.category-card {
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  padding: 20px;
  background: #fafafa;
  transition: box-shadow 0.3s;
}

.category-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.category-info h3 {
  margin: 0 0 8px 0;
  color: #262626;
  font-size: 1.1rem;
}

.category-info p {
  margin: 0 0 12px 0;
  color: #666;
  line-height: 1.4;
}

.category-meta {
  font-size: 0.85rem;
  color: #999;
}

.category-actions {
  margin-top: 16px;
  display: flex;
  gap: 8px;
}
</style>