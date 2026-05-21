<template>
  <div class="order-list-page">
    <header class="header">
      <button class="back" @click="$router.push('/')">←</button>
      <span>我的订单</span>
      <span></span>
    </header>

    <div class="tabs">
      <div
        v-for="tab in tabs"
        :key="tab.value"
        :class="['tab', { active: activeTab === tab.value }]"
        @click="activeTab = tab.value"
      >
        {{ tab.label }}
      </div>
    </div>

    <div class="order-list">
      <div v-for="order in filteredOrders" :key="order.id" class="order-card" @click="goToDetail(order.id)">
        <div class="order-header">
          <div class="order-no">
            <span>订单号：{{ order.order_no }}</span>
          </div>
          <span :class="['status', `status-${order.status}`]">{{ statusText(order.status) }}</span>
        </div>

        <div class="order-items">
          <div v-for="item in order.items" :key="item.food_id" class="item">
            <img :src="item.food_image" />
            <div class="info">
              <span class="name">{{ item.food_name }}</span>
              <span class="count">x{{ item.quantity }}</span>
            </div>
            <span class="price">¥{{ (item.price * item.quantity).toFixed(2) }}</span>
          </div>
        </div>

        <div class="order-footer">
          <span class="time">{{ formatTime(order.created_at) }}</span>
          <div class="total">
            <span>共{{ order.items.length }}件 实付</span>
            <span class="amount">¥{{ order.total_amount.toFixed(2) }}</span>
          </div>
        </div>

        <div class="order-actions" @click.stop>
          <button v-if="order.status === 0" class="btn-pay" @click="goToDetail(order.id)">去支付</button>
          <button v-if="order.status === 0" class="btn-cancel" @click="cancel(order.id)">取消订单</button>
        </div>
      </div>
    </div>

    <div v-if="filteredOrders.length === 0" class="empty">
      <div class="icon">📋</div>
      <p>暂无订单</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getOrderList, cancelOrder } from '../api'

const router = useRouter()

const orders = ref([])
const activeTab = ref(-1)

const tabs = [
  { label: '全部', value: -1 },
  { label: '待支付', value: 0 },
  { label: '进行中', value: 1 },
  { label: '已完成', value: 4 }
]

const statusText = (status) => {
  const map = { 0: '待支付', 1: '已支付', 2: '已接单', 3: '配送中', 4: '已完成', 5: '已取消', 6: '已退款' }
  return map[status] || '未知'
}

const filteredOrders = computed(() => {
  if (activeTab.value === -1) return orders.value
  if (activeTab.value === 1) {
    return orders.value.filter(o => o.status === 1 || o.status === 2 || o.status === 3)
  }
  return orders.value.filter(o => o.status === activeTab.value)
})

const formatTime = (time) => {
  return new Date(time).toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const fetchOrders = async () => {
  try {
    const res = await getOrderList({ page: 1, size: 50 })
    if (res.code === 0) {
      orders.value = res.data.items || []
    }
  } catch (err) {
    console.error('获取订单失败', err)
  }
}

const goToDetail = (id) => {
  router.push(`/order/${id}`)
}

const cancel = async (orderId) => {
  if (!confirm('确定要取消这个订单吗？')) return
  try {
    const res = await cancelOrder(orderId)
    if (res.code === 0) {
      fetchOrders()
    }
  } catch (err) {
    alert(err)
  }
}

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background: white;
  border-bottom: 1px solid #eee;
}

.back {
  font-size: 20px;
  background: none;
  border: none;
  cursor: pointer;
}

.tabs {
  display: flex;
  background: white;
  border-bottom: 1px solid #eee;
}

.tab {
  flex: 1;
  text-align: center;
  padding: 15px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
}

.tab.active {
  color: #ff6b6b;
  border-bottom-color: #ff6b6b;
}

.order-list {
  padding: 10px;
}

.order-card {
  background: white;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 10px;
  cursor: pointer;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.order-no {
  font-size: 13px;
  color: #666;
}

.status {
  font-size: 13px;
}

.status-0 { color: #ff6b6b; }
.status-1, .status-2, .status-3 { color: #1890ff; }
.status-4 { color: #52c41a; }
.status-5 { color: #999; }

.order-items {
  margin-bottom: 15px;
}

.order-items .item {
  display: flex;
  align-items: center;
  padding: 10px 0;
}

.order-items .item img {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  object-fit: cover;
  margin-right: 10px;
}

.order-items .info {
  flex: 1;
}

.order-items .name {
  display: block;
  font-size: 14px;
}

.order-items .count {
  font-size: 12px;
  color: #999;
}

.order-items .price {
  color: #333;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 10px;
  border-top: 1px solid #f5f5f5;
}

.order-footer .time {
  font-size: 12px;
  color: #999;
}

.order-footer .total {
  font-size: 13px;
}

.order-footer .amount {
  font-size: 16px;
  color: #333;
  font-weight: bold;
}

.order-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

.order-actions button {
  padding: 6px 15px;
  border-radius: 15px;
  font-size: 13px;
  cursor: pointer;
  border: 1px solid #ddd;
  background: white;
}

.btn-pay {
  background: #ff6b6b !important;
  color: white !important;
  border-color: #ff6b6b !important;
}

.empty {
  text-align: center;
  padding: 100px 20px;
}

.empty .icon {
  font-size: 60px;
  margin-bottom: 15px;
}

.empty p {
  color: #999;
}
</style>
