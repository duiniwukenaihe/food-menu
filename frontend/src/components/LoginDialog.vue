<template>
  <el-dialog
    v-model="dialogVisible"
    title="用户登录"
    width="400px"
    :before-close="handleClose"
  >
    <div class="login-container">
      <!-- 北极熊图片 -->
      <div class="polar-bear-image">
        <img src="/polar-bear.jpg" alt="北极熊" class="bear-img" />
        <p class="bear-text">🐻‍❄️ 欢迎来到美食世界</p>
      </div>

      <!-- 登录表单 -->
      <el-form
        ref="formRef"
        :model="loginForm"
        :rules="rules"
        label-width="80px"
        @submit.prevent="handleLogin"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
      </el-form>

      <div class="login-actions">
        <el-button :loading="loading" type="primary" @click="handleLogin">
          登录
        </el-button>
        <el-button @click="handleClose">取消</el-button>
      </div>

      <div class="login-tips">
        <p>测试账号：admin / admin123</p>
        <p>普通用户：user / user123</p>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores'

interface Props {
  modelValue: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const loginForm = reactive({
  username: '',
  password: ''
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

async function handleLogin() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true
    
    const success = await userStore.login(loginForm.username, loginForm.password)
    
    if (success) {
      ElMessage.success('登录成功！')
      emit('success')
      handleClose()
    } else {
      ElMessage.error('登录失败，请检查用户名和密码')
    }
  } catch (error) {
    console.error('Login validation failed:', error)
  } finally {
    loading.value = false
  }
}

function handleClose() {
  dialogVisible.value = false
  // 重置表单
  loginForm.username = ''
  loginForm.password = ''
  formRef.value?.resetFields()
}

// 监听对话框打开，重置表单
watch(dialogVisible, (newValue) => {
  if (newValue) {
    loginForm.username = ''
    loginForm.password = ''
    formRef.value?.resetFields()
  }
})
</script>

<style scoped>
.login-container {
  text-align: center;
}

.polar-bear-image {
  margin-bottom: 2rem;
}

.bear-img {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #e0f2fe;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem;
  font-size: 4rem;
}

.bear-text {
  color: #666;
  font-size: 1.1rem;
  margin: 0;
}

.login-form {
  text-align: left;
}

.login-actions {
  margin-top: 2rem;
  display: flex;
  justify-content: center;
  gap: 1rem;
}

.login-tips {
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
  font-size: 0.9rem;
  color: #999;
}

.login-tips p {
  margin: 0.25rem 0;
}
</style>