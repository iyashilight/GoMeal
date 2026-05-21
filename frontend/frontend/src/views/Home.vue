<template>
  <div class="home">
    <header class="header">
      <div class="logo">🍔 外卖系统</div>
      <div class="user" @click="showMenu = !showMenu">
        <span>{{ userStore.userInfo?.nickname || '用户' }}</span>
        <div v-if="showMenu" class="dropdown">
          <router-link to="/order-list">我的订单</router-link>
          <router-link to="/address">收货地址</router-link>
          <router-link to="/profile">个人中心</router-link>
          <router-link v-if="userStore.userInfo?.merchant_id" to="/merchant/manage">商家管理</router-link>
          <router-link v-if="!userStore.userInfo?.merchant_id" to="/merchant/register">商家注册</router-link>
          <a @click="logout">退出登录</a>
        </div>
      </div>
    </header>

    <div class="search-bar">
      <input v-model="searchKey" placeholder="搜索商家或商品" />
      <button>搜索</button>
    </div>

    <div class="merchant-list">
      <div
        v-for="merchant in filteredMerchants"
        :key="merchant.id"
        class="merchant-card"
        @click="goToMerchant(merchant.id)"
      >
        <img :src="merchant.logo" :alt="merchant.name" />
        <div class="info">
          <h3>{{ merchant.name }}</h3>
          <div class="rating">
            <span class="stars">⭐ {{ merchant.rating }}</span>
            <span>月售 {{ merchant.sales }}</span>
          </div>
          <div class="delivery">
            <span>起送 ¥{{ merchant.min_price }}</span>
            <span>配送 ¥{{ merchant.delivery_fee }}</span>
          </div>
          <div v-if="merchant.notice" class="notice">
            {{ merchant.notice }}
          </div>
        </div>
      </div>
    </div>

    <!-- 底部购物车入口 -->
    <div v-if="cartStore.items.length > 0" class="cart-float" @click="goToCart">
      <div class="cart-icon">🛒</div>
      <div class="cart-info">
        <span class="total">¥{{ cartStore.finalAmount.toFixed(2) }}</span>
        <span class="count">{{ cartStore.totalQuantity }}件商品</span>
      </div>
      <button class="checkout-btn">去结算</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useCartStore } from '../stores/cart'
import { getMerchantList } from '../api'

const router = useRouter()
const userStore = useUserStore()
const cartStore = useCartStore()

const merchants = ref([])
const searchKey = ref('')
const showMenu = ref(false)

const filteredMerchants = computed(() => {
  if (!searchKey.value) return merchants.value
  return merchants.value.filter(m =>
    m.name.includes(searchKey.value) ||
    m.notice?.includes(searchKey.value)
  )
})

const fetchMerchants = async () => {
  try {
    const res = await getMerchantList()
    if (res.code === 0) {
      merchants.value = res.data
    }
  } catch (err) {
    console.error('获取商家列表失败', err)
  }
}

const goToMerchant = (id) => {
  router.push(`/merchant/${id}`)
}

const goToCart = () => {
  router.push('/cart')
}

const logout = () => {
  userStore.logout()
  router.push('/login')
}

onMounted(() => {
  fetchMerchants()
  cartStore.fetchCart()
  userStore.fetchUserInfo()
})
</script>

<style scoped>
.home {
  padding-bottom: 80px;
}

.header {
  background: linear-gradient(135deg, #ff6b6b 0%, #ff8e8e 100%);
  color: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-size: 20px;
  font-weight: bold;
}

.user {
  position: relative;
  cursor: pointer;
  padding: 5px 10px;
}

.dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  padding: 8px 0;
  min-width: 120px;
  z-index: 100;
}

.dropdown a {
  display: block;
  padding: 10px 20px;
  color: #333;
  text-decoration: none;
}

.dropdown a:hover {
  background: #f5f5f5;
}

.search-bar {
  display: flex;
  padding: 15px;
  gap: 10px;
  background: white;
}

.search-bar input {
  flex: 1;
  padding: 12px 15px;
  border: 1px solid #ddd;
  border-radius: 20px;
  font-size: 14px;
}

.search-bar button {
  padding: 12px 20px;
  background: #ff6b6b;
  color: white;
  border: none;
  border-radius: 20px;
  cursor: pointer;
}

.merchant-list {
  padding: 10px;
}

.merchant-card {
  display: flex;
  background: white;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  cursor: pointer;
  transition: transform 0.2s;
}

.merchant-card:hover {
  transform: translateY(-2px);
}

.merchant-card img {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  object-fit: cover;
  margin-right: 15px;
}

.merchant-card .info {
  flex: 1;
}

.merchant-card h3 {
  font-size: 16px;
  margin-bottom: 8px;
}

.rating, .delivery {
  display: flex;
  gap: 15px;
  font-size: 13px;
  color: #666;
  margin-bottom: 5px;
}

.stars {
  color: #ff9500;
}

.notice {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.cart-float {
  position: fixed;
  bottom: 20px;
  left: 20px;
  right: 20px;
  background: #333;
  border-radius: 30px;
  padding: 10px 15px;
  display: flex;
  align-items: center;
  color: white;
  cursor: pointer;
  box-shadow: 0 4px 15px rgba(0,0,0,0.3);
}

.cart-icon {
  font-size: 24px;
  margin-right: 10px;
}

.cart-info {
  flex: 1;
}

.cart-info .total {
  font-size: 18px;
  font-weight: bold;
}

.cart-info .count {
  font-size: 12px;
  color: #ccc;
  margin-left: 8px;
}

.checkout-btn {
  background: #ff6b6b;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 20px;
  font-size: 14px;
}
</style>