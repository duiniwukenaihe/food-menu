<template>
  <el-dialog
    v-model="dialogVisible"
    title="确认订单"
    width="600px"
    :before-close="handleClose"
  >
    <div class="order-container">
      <!-- 订单项目列表 -->
      <div v-if="orderItems.length > 0" class="order-items">
        <div
          v-for="(item, index) in orderItems"
          :key="item.dish.id"
          class="order-item"
        >
          <img 
            :src="item.dish.image_url || '/placeholder-food.jpg'" 
            :alt="item.dish.name" 
            class="item-image"
          />
          <div class="item-info">
            <h4>{{ item.dish.name }}</h4>
            <p class="item-price">¥{{ item.dish.price }}</p>
          </div>
          <div class="item-quantity">
            <el-input-number
              v-model="item.quantity"
              :min="1"
              :max="99"
              size="small"
              @change="updateTotal"
            />
          </div>
          <div class="item-total">
            ¥{{ (item.dish.price * item.quantity).toFixed(2) }}
          </div>
          <el-button 
            type="danger" 
            size="small" 
            @click="removeItem(index)"
            circle
          >
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </div>

      <div v-else class="empty-order">
        <p>订单为空</p>
      </div>

      <!-- 订单总计 -->
      <div class="order-summary">
        <div class="summary-row">
          <span>商品总数：</span>
          <span>{{ totalQuantity }} 件</span>
        </div>
        <div class="summary-row total">
          <span>总计：</span>
          <span class="total-amount">¥{{ totalAmount.toFixed(2) }}</span>
        </div>
      </div>

      <!-- 备注 -->
      <div class="order-note">
        <el-input
          v-model="note"
          type="textarea"
          placeholder="订单备注（可选）"
          :rows="3"
          maxlength="200"
          show-word-limit
        />
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button 
          type="primary" 
          :loading="submitting"
          :disabled="orderItems.length === 0"
          @click="submitOrder"
        >
          确认下单
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'
import { useOrderStore } from '@/stores'
import type { Dish } from '@/types'

interface Props {
  modelValue: boolean
  items: { dish: Dish; quantity: number }[]
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const orderStore = useOrderStore()
const submitting = ref(false)
const note = ref('')

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const orderItems = ref<{ dish: Dish; quantity: number }[]>([])

const totalQuantity = computed(() => 
  orderItems.value.reduce((sum, item) => sum + item.quantity, 0)
)

const totalAmount = computed(() => 
  orderItems.value.reduce((sum, item) => sum + (item.dish.price * item.quantity), 0)
)

// 监听外部传入的items变化
watch(() => props.items, (newItems) => {
  orderItems.value = [...newItems]
}, { immediate: true, deep: true })

function updateTotal() {
  // 总计会自动更新
}

function removeItem(index: number) {
  orderItems.value.splice(index, 1)
  ElMessage.success('已移除商品')
}

async function submitOrder() {
  if (orderItems.value.length === 0) {
    ElMessage.warning('订单为空，无法提交')
    return
  }

  try {
    submitting.value = true
    
    const orderData = orderItems.value.map(item => ({
      dish_id: item.dish.id,
      quantity: item.quantity
    }))

    await orderStore.createOrder(orderData)
    
    ElMessage.success('订单提交成功！')
    emit('success')
    handleClose()
  } catch (error) {
    console.error('Submit order failed:', error)
    ElMessage.error('订单提交失败，请重试')
  } finally {
    submitting.value = false
  }
}

function handleClose() {
  dialogVisible.value = false
  note.value = ''
}
</script>

<style scoped>
.order-container {
  max-height: 500px;
  overflow-y: auto;
}

.order-items {
  margin-bottom: 20px;
}

.order-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border: 1px solid #eee;
  border-radius: 8px;
  margin-bottom: 10px;
  background: #fafafa;
}

.item-image {
  width: 60px;
  height: 60px;
  border-radius: 4px;
  object-fit: cover;
  margin-right: 12px;
}

.item-info {
  flex: 1;
  margin-right: 12px;
}

.item-info h4 {
  margin: 0 0 4px 0;
  font-size: 14px;
  color: #333;
}

.item-price {
  margin: 0;
  color: #666;
  font-size: 12px;
}

.item-quantity {
  margin-right: 12px;
}

.item-total {
  min-width: 80px;
  text-align: right;
  font-weight: bold;
  color: #f56c6c;
  margin-right: 12px;
}

.empty-order {
  text-align: center;
  padding: 40px;
  color: #999;
}

.order-summary {
  border-top: 1px solid #eee;
  padding-top: 16px;
  margin-bottom: 16px;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.summary-row.total {
  font-weight: bold;
  font-size: 16px;
  color: #f56c6c;
}

.total-amount {
  font-size: 18px;
}

.order-note {
  margin-top: 16px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>