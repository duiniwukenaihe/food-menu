<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑分类' : '添加分类'"
    width="500px"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="80px"
    >
      <el-form-item label="分类名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入分类名称" />
      </el-form-item>
      
      <el-form-item label="分类描述">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入分类描述（可选）"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { api } from '@/utils/api'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '@/types'

interface Props {
  modelValue: boolean
  category?: Category | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  category: null
})

const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()
const loading = ref(false)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const isEdit = computed(() => !!props.category)

const form = reactive<CreateCategoryRequest>({
  name: '',
  description: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 2, max: 20, message: '分类名称长度在 2 到 20 个字符', trigger: 'blur' }
  ]
}

// 监听category变化，初始化表单
watch(() => props.category, (category) => {
  if (category) {
    form.name = category.name
    form.description = category.description || ''
  } else {
    resetForm()
  }
})

// 监听对话框打开，重置表单
watch(dialogVisible, (newValue) => {
  if (newValue && !isEdit.value) {
    resetForm()
  }
})

function resetForm() {
  form.name = ''
  form.description = ''
  formRef.value?.resetFields()
}

async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    if (isEdit.value && props.category) {
      const updateData: UpdateCategoryRequest = {
        name: form.name,
        description: form.description
      }

      await api.updateCategory(props.category.id, updateData)
      ElMessage.success('分类更新成功')
    } else {
      await api.createCategory(form)
      ElMessage.success('分类创建成功')
    }

    emit('success')
  } catch (error) {
    console.error('Submit failed:', error)
    ElMessage.error('操作失败，请重试')
  } finally {
    loading.value = false
  }
}

function handleClose() {
  dialogVisible.value = false
  resetForm()
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>