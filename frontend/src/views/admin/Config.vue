<template>
  <div class="admin-config">
    <div class="page-header">
      <h2>系统配置</h2>
      <el-button type="primary" :loading="loading" @click="saveConfig">
        <el-icon><Check /></el-icon>
        保存配置
      </el-button>
    </div>

    <div class="config-sections">
      <!-- 基本配置 -->
      <div class="config-section">
        <h3>基本配置</h3>
        <div class="config-grid">
          <div class="config-item">
            <label>默认荤菜数量</label>
            <el-input-number
              v-model="config.default_meat_count"
              :min="0"
              :max="6"
              size="small"
            />
          </div>
          <div class="config-item">
            <label>默认素菜数量</label>
            <el-input-number
              v-model="config.default_vegetable_count"
              :min="0"
              :max="6"
              size="small"
            />
          </div>
          <div class="config-item">
            <label>菜品最大数量</label>
            <el-input-number
              v-model="config.max_dish_count"
              :min="1"
              :max="20"
              size="small"
            />
          </div>
        </div>
      </div>

      <!-- S3配置 -->
      <div class="config-section">
        <h3>S3对象存储配置</h3>
        <div class="config-grid">
          <div class="config-item full-width">
            <label>S3端点</label>
            <el-input
              v-model="config.s3_endpoint"
              placeholder="例如：https://oss-cn-beijing.aliyuncs.com"
              size="small"
            />
          </div>
          <div class="config-item">
            <label>访问密钥ID</label>
            <el-input
              v-model="config.s3_access_key"
              placeholder="Access Key ID"
              size="small"
              show-password
            />
          </div>
          <div class="config-item">
            <label>访问密钥Secret</label>
            <el-input
              v-model="config.s3_secret_key"
              placeholder="Access Key Secret"
              size="small"
              show-password
            />
          </div>
          <div class="config-item">
            <label>存储桶名称</label>
            <el-input
              v-model="config.s3_bucket"
              placeholder="Bucket名称"
              size="small"
            />
          </div>
          <div class="config-item">
            <label>区域</label>
            <el-input
              v-model="config.s3_region"
              placeholder="例如：oss-cn-beijing"
              size="small"
            />
          </div>
        </div>
      </div>

      <!-- 说明信息 -->
      <div class="config-section">
        <h3>配置说明</h3>
        <div class="help-content">
          <h4>基本配置</h4>
          <ul>
            <li><strong>默认荤菜数量：</strong>随机搭配时默认的荤菜数量</li>
            <li><strong>默认素菜数量：</strong>随机搭配时默认的素菜数量</li>
            <li><strong>菜品最大数量：</strong>单次订单最多可选择的菜品数量</li>
          </ul>
          
          <h4>S3对象存储</h4>
          <p>支持以下S3兼容的对象存储服务：</p>
          <ul>
            <li><strong>阿里云OSS：</strong>端点格式为 https://oss-cn-{region}.aliyuncs.com</li>
            <li><strong>AWS S3：</strong>端点格式为 https://s3.{region}.amazonaws.com</li>
            <li><strong>腾讯云COS：</strong>端点格式为 https://cos.{region}.myqcloud.com</li>
            <li><strong>MinIO：</strong>自建MinIO服务的端点地址</li>
          </ul>
          <p class="warning">
            <el-icon><Warning /></el-icon>
            配置修改后需要重启后端服务才能生效
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Check, Warning } from '@element-plus/icons-vue'
import { api } from '@/utils/api'
import type { SystemConfig } from '@/types'

const loading = ref(false)

const config = reactive({
  default_meat_count: 1,
  default_vegetable_count: 2,
  max_dish_count: 6,
  s3_endpoint: '',
  s3_access_key: '',
  s3_secret_key: '',
  s3_bucket: '',
  s3_region: ''
})

onMounted(loadConfig)

async function loadConfig() {
  try {
    const configs = await api.getConfig()
    
    // 将配置数组转换为对象
    configs.forEach(item => {
      const key = item.config_key
      const value = item.config_value
      
      if (key in config) {
        // 数字类型的配置
        if (['default_meat_count', 'default_vegetable_count', 'max_dish_count'].includes(key)) {
          config[key as keyof typeof config] = parseInt(value) || 0
        } else {
          config[key as keyof typeof config] = value
        }
      }
    })
  } catch (error) {
    console.error('Failed to load config:', error)
    ElMessage.error('加载配置失败')
  }
}

async function saveConfig() {
  try {
    loading.value = true
    
    // 验证配置
    if (config.max_dish_count < 1) {
      ElMessage.error('菜品最大数量不能小于1')
      return
    }
    
    if (config.default_meat_count + config.default_vegetable_count > config.max_dish_count) {
      ElMessage.error('默认荤素数量之和不能超过最大菜品数量')
      return
    }
    
    // 转换为API需要的格式
    const configData: Record<string, string> = {}
    Object.entries(config).forEach(([key, value]) => {
      configData[key] = String(value)
    })
    
    await api.updateConfig(configData)
    ElMessage.success('配置保存成功')
  } catch (error) {
    console.error('Failed to save config:', error)
    ElMessage.error('保存配置失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admin-config {
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

.page-header h2 {
  margin: 0;
  color: #262626;
}

.config-sections {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.config-section {
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  padding: 20px;
}

.config-section h3 {
  margin: 0 0 20px 0;
  color: #262626;
  font-size: 1.1rem;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 10px;
}

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.config-item.full-width {
  grid-column: 1 / -1;
}

.config-item label {
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.help-content {
  line-height: 1.6;
}

.help-content h4 {
  color: #262626;
  margin: 20px 0 10px 0;
  font-size: 1rem;
}

.help-content h4:first-child {
  margin-top: 0;
}

.help-content ul {
  margin: 10px 0;
  padding-left: 20px;
}

.help-content li {
  margin-bottom: 6px;
  color: #666;
}

.help-content p {
  color: #666;
  margin: 10px 0;
}

.warning {
  color: #e6a23c;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  padding: 12px;
  background: #fdf6ec;
  border: 1px solid #f5dab1;
  border-radius: 4px;
}

.warning .el-icon {
  font-size: 16px;
}

:deep(.el-input-number) {
  width: 100%;
}
</style>