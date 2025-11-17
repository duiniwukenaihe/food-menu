<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑菜品' : '添加菜品'"
    width="800px"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      size="default"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="菜品名称" prop="name">
            <el-input v-model="form.name" placeholder="请输入菜品名称" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="分类" prop="category_id">
            <el-select v-model="form.category_id" placeholder="请选择分类" style="width: 100%">
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="价格" prop="price">
            <el-input-number
              v-model="form.price"
              :min="0"
              :precision="2"
              placeholder="请输入价格"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="状态">
            <el-switch
              v-model="form.is_active"
              active-text="启用"
              inactive-text="禁用"
            />
            <el-switch
              v-model="form.is_seasonal"
              active-text="应季推荐"
              inactive-text="普通菜品"
              style="margin-left: 20px"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="菜品描述">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入菜品描述"
        />
      </el-form-item>

      <el-form-item label="图片URL">
        <el-input
          v-model="form.image_url"
          placeholder="请输入图片URL"
        />
        <div v-if="form.image_url" class="image-preview">
          <el-image
            :src="form.image_url"
            alt="预览"
            style="width: 100px; height: 100px; object-fit: cover; border-radius: 4px;"
            fit="cover"
          />
        </div>
      </el-form-item>

      <el-form-item label="视频URL">
        <el-input
          v-model="form.video_url"
          placeholder="请输入视频URL"
        />
      </el-form-item>

      <el-form-item label="制作步骤">
        <el-input
          v-model="form.cooking_steps"
          type="textarea"
          :rows="6"
          placeholder="请输入制作步骤，每行一个步骤"
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
import type { Dish, Category, CreateDishRequest, UpdateDishRequest } from '@/types'

interface Props {
  modelValue: boolean
  dish?: Dish | null
  categories: Category[]
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  dish: null
})

const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()
const loading = ref(false)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const isEdit = computed(() => !!props.dish)

const form = reactive<CreateDishRequest>({
  name: '',
  description: '',
  category_id: 0,
  price: 0,
  image_url: '',
  video_url: '',
  cooking_steps: '',
  is_seasonal: false
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入菜品名称', trigger: 'blur' },
    { min: 2, max: 50, message: '菜品名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '价格不能小于0', trigger: 'blur' }
  ]
}

// 监听dish变化，初始化表单
watch(() => props.dish, (dish) => {
  if (dish) {
    form.name = dish.name
    form.description = dish.description || ''
    form.category_id = dish.category_id
    form.price = dish.price
    form.image_url = dish.image_url || ''
    form.video_url = dish.video_url || ''
    form.cooking_steps = dish.cooking_steps || ''
    form.is_seasonal = dish.is_seasonal
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
  form.category_id = 0
  form.price = 0
  form.image_url = ''
  form.video_url = ''
  form.cooking_steps = ''
  form.is_seasonal = false
  formRef.value?.resetFields()
}

async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    if (isEdit.value && props.dish) {
      const updateData: UpdateDishRequest = {
        name: form.name,
        description: form.description,
        category_id: form.category_id,
        price: form.price,
        image_url: form.image_url,
        video_url: form.video_url,
        cooking_steps: form.cooking_steps,
        is_seasonal: form.is_seasonal,
        is_active: props.dish.is_active
      }

      await api.updateDish(props.dish.id, updateData)
      ElMessage.success('菜品更新成功')
    } else {
      await api.createDish(form)
      ElMessage.success('菜品创建成功')
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
.image-preview {
  margin-top: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>