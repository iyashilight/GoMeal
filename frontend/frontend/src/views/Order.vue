<template>
  <div class="order-page">
    <header class="header">
      <button class="back" @click="$router.push('/order-list')">←</button>
      <span>订单详情</span>
      <span></span>
    </header>

    <div v-if="order" class="order-content">
      <!-- 订单状态 -->
      <div :class="['status-section', { pending: order.status === 0 }]">
        <div class="status-icon">{{ order.status === 0 ? '⏳' : '✓' }}</div>
        <div class="status-text">{{ statusText(order.status) }}</div>
        <div v-if="order.status === 0" class="status-action">
          <button class="pay-btn" @click="handlePay" :disabled="paying">
            {{ paying ? '支付中...' : '去支付' }}
          </button>
        </div>
      </div>

      <!-- 商品列表 -->
      <div class="items-section">
        <div v-for="item in order.items" :key="item.food_id" class="item">
          <img :src="item.food_image" />
          <div class="info">
            <span class="name">{{ item.food_name }}</span>
            <span class="count">x{{ item.quantity }}</span>
          </div>
          <span class="price">¥{{ (item.price * item.quantity).toFixed(2) }}</span>
        </div>
      </div>

      <!-- 金额明细 -->
      <div class="amount-section">
        <div class="row">
          <span>商品总额</span>
          <span>¥{{ order.total_amount.toFixed(2) }}</span>
        </div>
        <div class="row">
          <span>配送费</span>
          <span>¥{{ order.delivery_fee.toFixed(2) }}</span>
        </div>
        <div class="row total">
          <span>实付金额</span>
          <span class="price">¥{{ (order.total_amount + order.delivery_fee).toFixed(2) }}</span>
        </div>
      </div>

      <!-- 订单信息 -->
      <div class="order-info-section">
        <h4>订单信息</h4>
        <p>订单编号：{{ order.order_no }}</p>
        <p>下单时间：{{ formatTime(order.created_at) }}</p>
        <p v-if="order.remark">订单备注：{{ order.remark }}</p>
      </div>

      <!-- 支付信息 -->
      <div v-if="payment" class="payment-section">
        <h4>支付信息</h4>
        <p>支付方式：{{ payment.method || '-' }}</p>
        <p>交易单号：{{ payment.trade_no || '-' }}</p>
        <p>支付状态：
          <span :class="payment.status === 1 ? 'paid' : 'unpaid'">
            {{ payment.status === 1 ? '已支付' : '未支付' }}
          </span>
        </p>
        <p v-if="payment.paid_at">支付时间：{{ payment.paid_at }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getOrderDetail, payOrder, getPayment } from '../api'

const route = useRoute()
const order = ref(null)
const payment = ref(null)
const paying = ref(false)

const statusText = (status) => {
  const map = { 0: '待支付', 1: '已支付', 2: '已接单', 3: '配送中', 4: '已完成', 5: '已取消', 6: '已退款' }
  return map[status] || '未知'
}

const formatTime = (time) => {
  return new Date(time).toLocaleString('zh-CN')
}

const fetchOrder = async () => {
  try {
    const res = await getOrderDetail(route.params.id)
    if (res.code === 0) {
      order.value = res.data
    }
  } catch (err) {
    console.error('获取订单失败', err)
  }
}

const fetchPayment = async () => {
  try {
    const res = await getPayment(route.params.id)
    if (res.code === 0) {
      payment.value = res.data
    }
  } catch (err) {
    // 支付信息可能不存在，静默处理
  }
}

const handlePay = async () => {
  paying.value = true
  try {
    const res = await payOrder(route.params.id)
    if (res.code === 0) {
      payment.value = res.data
      order.value.status = 1 // 更新订单状态为已支付
    }
  } catch (err) {
    console.error('支付失败', err)
  } finally {
    paying.value = false
  }
}

onMounted(() => {
  fetchOrder()
  fetchPayment()
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

.order-content {
  padding: 15px;
}

.status-section {
  text-align: center;
  padding: 30px;
  background: linear-gradient(135deg, #52c41a 0%, #73d13d 100%);
  border-radius: 12px;
  color: white;
  margin-bottom: 15px;
}

.status-icon {
  width: 50px;
  height: 50px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: #52c41a;
  margin: 0 auto 10px;
}

.status-text {
  font-size: 18px;
}

.status-section.pending {
  background: linear-gradient(135deg, #ff9500 0%, #ffb340 100%);
}

.status-action {
  margin-top: 15px;
}

.pay-btn {
  background: white;
  color: #ff6b6b;
  border: none;
  padding: 12px 40px;
  border-radius: 25px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: transform 0.2s;
}

.pay-btn:hover {
  transform: scale(1.05);
}

.pay-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.items-section,
.amount-section,
.order-info-section {
  background: white;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 15px;
}

.item {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #f5f5f5;
}

.item:last-child {
  border-bottom: none;
}

.item img {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  object-fit: cover;
  margin-right: 10px;
}

.item .info {
  flex: 1;
}

.item .name {
  display: block;
  font-size: 14px;
}

.item .count {
  font-size: 12px;
  color: #999;
}

.amount-section .row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
}

.amount-section .total {
  border-top: 1px solid #eee;
  margin-top: 10px;
  padding-top: 15px;
  font-weight: bold;
}

.amount-section .price {
  color: #ff6b6b;
  font-size: 18px;
}

.order-info-section h4,
.payment-section h4 {
  margin-bottom: 10px;
  font-size: 14px;
}

.order-info-section p,
.payment-section p {
  font-size: 13px;
  color: #666;
  line-height: 1.8;
}

.payment-section {
  background: white;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 15px;
}

.paid {
  color: #52c41a;
  font-weight: bold;
}

.unpaid {
  color: #ff9500;
  font-weight: bold;
}
</style>
