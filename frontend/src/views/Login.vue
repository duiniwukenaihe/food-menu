<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <!-- 北极熊图片 -->
        <div class="polar-bear-section">
          <div class="bear-avatar">
            🐻‍❄️
          </div>
          <h1 class="welcome-title">欢迎来到美食世界</h1>
          <p class="welcome-subtitle">登录后开始您的美味之旅</p>
        </div>

        <!-- 登录表单 -->
        <div class="form-section">
          <el-form
            ref="formRef"
            :model="loginForm"
            :rules="rules"
            label-position="top"
            size="large"
            @submit.prevent="handleLogin"
          >
            <el-form-item label="用户名" prop="username">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                :prefix-icon="User"
                clearable
              />
            </el-form-item>
            
            <el-form-item label="密码" prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                :prefix-icon="Lock"
                show-password
                clearable
                @keyup.enter="handleLogin"
              />
            </el-form-item>

            <el-form-item>
              <el-button 
                type="primary" 
                size="large" 
                :loading="loading"
                style="width: 100%"
                @click="handleLogin"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>

          <!-- 测试账号提示 -->
          <div class="test-accounts">
            <h3>测试账号</h3>
            <div class="account-item">
              <span class="account-type">管理员：</span>
              <span>admin / admin123</span>
            </div>
            <div class="account-item">
              <span class="account-type">普通用户：</span>
              <span>user / user123</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 熊掌图案 -->
      <div class="bear-paw-decoration">
        <div class="paw-print">🐾</div>
        <div class="paw-text">美味从这里开始</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

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
      router.push('/')
    } else {
      ElMessage.error('登录失败，请检查用户名和密码')
    }
  } catch (error) {
    console.error('Login validation failed:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  position: relative;
  width: 100%;
  max-width: 400px;
}

.login-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  animation: slideUp 0.6s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.polar-bear-section {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
  padding: 40px 20px;
  text-align: center;
}

.bear-avatar {
  font-size: 4rem;
  margin-bottom: 16px;
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-10px);
  }
  60% {
    transform: translateY(-5px);
  }
}

.welcome-title {
  color: #1565c0;
  font-size: 1.5rem;
  margin: 0 0 8px 0;
  font-weight: 600;
}

.welcome-subtitle {
  color: #64b5f6;
  font-size: 0.9rem;
  margin: 0;
}

.form-section {
  padding: 40px 30px;
}

.test-accounts {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.test-accounts h3 {
  font-size: 0.9rem;
  color: #666;
  margin: 0 0 12px 0;
  text-align: center;
}

.account-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  font-size: 0.85rem;
}

.account-type {
  color: #999;
}

.account-item span:last-child {
  color: #409eff;
  font-family: monospace;
  background: #f0f9ff;
  padding: 2px 6px;
  border-radius: 4px;
}

.bear-paw-decoration {
  position: absolute;
  bottom: -60px;
  left: 50%;
  transform: translateX(-50%);
  text-align: center;
  opacity: 0.8;
}

.paw-print {
  font-size: 2rem;
  animation: wave 3s infinite ease-in-out;
  margin-bottom: 8px;
}

@keyframes wave {
  0%, 100% {
    transform: rotate(-5deg);
  }
  50% {
    transform: rotate(5deg);
  }
}

.paw-text {
  color: white;
  font-size: 0.9rem;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

:deep(.el-form-item__label) {
  color: #666;
  font-weight: 500;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
}

:deep(.el-button) {
  border-radius: 8px;
  font-weight: 500;
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #409eff 0%, #66b3ff 100%);
  border: none;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(135deg, #66b3ff 0%, #409eff 100%);
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.4);
}
</style>