<template>
  <div class="cart-page">
    <header class="header">
      <button class="back" @click="$router.back()">←</button>
      <span>购物车</span>
      <span></span>
    </header>

    <div v-if="!cartStore.items.length" class="empty">
      <div class="icon">🛒</div>
      <p>购物车是空的</p>
      <router-link to="/" class="btn-primary">去逛逛</router-link>
    </div>

    <template v-else>
      <div class="merchant-info">
        <h3>{{ cartStore.merchantName }}</h3>
      </div>

      <div class="cart-list">
        <div v-for="item in cartStore.items" :key="item.id" class="cart-item">
          <img :src="item.food_image" :alt="item.food_name" />
          <div class="info">
            <h4>{{ item.food_name }}</h4>
            <div class="price-row">
              <span class="price">¥{{ item.price }}</span>
              <div class="quantity-control">
                <button @click="decrease(item)">-</button>
                <span>{{ item.quantity }}</span>
                <button @click="increase(item)">+</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="cart-summary">
        <div class="summary-row">
          <span>商品总额</span>
          <span>¥{{ cartStore.totalAmount.toFixed(2) }}</span>
        </div>
        <div class="summary-row">
          <span>配送费</span>
          <span>¥{{ cartStore.deliveryFee.toFixed(2) }}</span>
        </div>
        <div class="summary-row total">
          <span>合计</span>
          <span class="total-price">¥{{ cartStore.finalAmount.toFixed(2) }}</span>
        </div>
      </div>

      <div class="address-section" @click="goToAddress">
        <div v-if="selectedAddress" class="address-info">
          <div class="addr-header">
            <span class="name">{{ selectedAddress.name }}</span>
            <span class="phone">{{ selectedAddress.phone }}</span>
          </div>
          <div class="addr-detail">{{ selectedAddress.address }}</div>
        </div>
        <div v-else class="no-address">
          请选择收货地址 →
        </div>
      </div>

      <div class="remark-section">
        <input v-model="remark" placeholder="订单备注（选填）" />
      </div>

      <div class="cart-footer">
        <div class="total">
          <span>合计：</span>
          <span class="price">¥{{ cartStore.finalAmount.toFixed(2) }}</span>
        </div>
        <button
          :class="['submit-btn', { disabled: !selectedAddress }]"
          @click="submitOrder"
          :disabled="!selectedAddress || submitting"
        >
          {{ submitting ? '提交中...' : '提交订单' }}
        </button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { getAddressList, createOrder } from '../api'

const router = useRouter()
const cartStore = useCartStore()

const addresses = ref([])
const selectedAddress = ref(null)
const remark = ref('')
const submitting = ref(false)

const fetchAddresses = async () => {
  try {
    const res = await getAddressList()
    if (res.code === 0 && res.data) {
      addresses.value = res.data
      selectedAddress.value = res.data.find(a => a.is_default) || res.data[0]
    }
  } catch (err) {
    console.error('获取地址失败', err)
  }
}

const decrease = async (item) => {
  if (item.quantity > 1) {
    await cartStore.updateQuantity(item.id, item.quantity - 1)
  } else {
    await cartStore.removeItem(item.id)
  }
}

const increase = async (item) => {
  await cartStore.updateQuantity(item.id, item.quantity + 1)
}

const goToAddress = () => {
  router.push('/address?from=cart')
}

const submitOrder = async () => {
  if (!selectedAddress.value) return

  submitting.value = true
  try {
    const res = await createOrder({
      address_id: selectedAddress.value.id,
      delivery_fee: cartStore.deliveryFee,
      remark: remark.value
    })
    if (res.code === 0) {
      cartStore.clear()
      router.push(`/order-list`)
    } else {
      alert(res.message || '下单失败')
    }
  } catch (err) {
    alert(err || '下单失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  cartStore.fetchCart()
  fetchAddresses()
})
</script>

<style scoped>
.cart-page {
  padding-bottom: 80px;
}

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

.empty {
  text-align: center;
  padding: 100px 20px;
}

.empty .icon {
  font-size: 80px;
  margin-bottom: 20px;
}

.empty p {
  color: #999;
  margin-bottom: 30px;
}

.merchant-info {
  padding: 15px;
  background: white;
  margin-bottom: 10px;
}

.merchant-info h3 {
  font-size: 16px;
}

.cart-list {
  background: white;
}

.cart-item {
  display: flex;
  padding: 15px;
  border-bottom: 1px solid #f5f5f5;
}

.cart-item img {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  object-fit: cover;
  margin-right: 15px;
}

.cart-item .info {
  flex: 1;
}

.cart-item h4 {
  font-size: 15px;
  margin-bottom: 20px;
}

.price-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price {
  font-size: 16px;
  color: #ff6b6b;
  font-weight: bold;
}

.quantity-control {
  display: flex;
  align-items: center;
  gap: 10px;
}

.quantity-control button {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 1px solid #ddd;
  background: white;
  font-size: 16px;
  cursor: pointer;
}

.cart-summary {
  background: white;
  padding: 15px;
  margin-top: 10px;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  color: #666;
}

.summary-row.total {
  border-top: 1px solid #eee;
  margin-top: 10px;
  padding-top: 15px;
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.total-price {
  color: #ff6b6b;
  font-size: 18px;
}

.address-section {
  background: white;
  padding: 15px;
  margin-top: 10px;
  cursor: pointer;
}

.addr-header {
  margin-bottom: 5px;
}

.addr-header .name {
  font-weight: bold;
  margin-right: 10px;
}

.addr-header .phone {
  color: #666;
}

.addr-detail {
  color: #666;
  font-size: 14px;
}

.no-address {
  color: #ff6b6b;
}

.remark-section {
  background: white;
  padding: 15px;
  margin-top: 10px;
}

.remark-section input {
  width: 100%;
  padding: 10px;
  border: 1px solid #eee;
  border-radius: 8px;
}

.cart-footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  display: flex;
  align-items: center;
  padding: 10px 15px;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.1);
}

.cart-footer .total {
  flex: 1;
}

.cart-footer .price {
  font-size: 20px;
  color: #ff6b6b;
  font-weight: bold;
}

.submit-btn {
  background: #ff6b6b;
  color: white;
  border: none;
  padding: 12px 30px;
  border-radius: 20px;
  font-size: 16px;
}

.submit-btn.disabled {
  background: #ccc;
}
</style>
