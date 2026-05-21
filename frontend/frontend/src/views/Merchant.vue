<template>
  <div class="merchant-page">
    <header class="header">
      <button class="back" @click="$router.back()">←</button>
      <span>{{ merchant?.name || '商家详情' }}</span>
      <span></span>
    </header>

    <div v-if="merchant" class="merchant-info">
      <img :src="merchant.logo" class="logo" />
      <div class="info">
        <h1>{{ merchant.name }}</h1>
        <div class="rating">
          <span>⭐ {{ merchant.rating }}</span>
          <span>月售 {{ merchant.sales }}</span>
        </div>
        <div class="notice">{{ merchant.notice }}</div>
        <div class="delivery-info">
          <span>起送 ¥{{ merchant.min_price }}</span>
          <span>配送 ¥{{ merchant.delivery_fee }}</span>
        </div>
      </div>
    </div>

    <div class="menu">
      <div class="category-sidebar">
        <div
          v-for="cat in categories"
          :key="cat.id"
          :class="['category-item', { active: activeCategory === cat.id }]"
          @click="scrollToCategory(cat.id)"
        >
          {{ cat.name }}
        </div>
      </div>

      <div class="food-list">
        <div
          v-for="cat in categories"
          :key="cat.id"
          :id="'cat-' + cat.id"
          class="category-section"
        >
          <h3>{{ cat.name }}</h3>
          <div
            v-for="food in cat.foods"
            :key="food.id"
            class="food-item"
          >
            <img :src="food.image" :alt="food.name" />
            <div class="food-info">
              <h4>{{ food.name }}</h4>
              <p class="desc">{{ food.description }}</p>
              <div class="food-meta">
                <span class="sales">月售 {{ food.sales }}</span>
              </div>
              <div class="food-price">
                <span class="price">¥{{ food.price }}</span>
                <span v-if="food.old_price" class="old-price">¥{{ food.old_price }}</span>
                <div class="quantity-control">
                  <button
                    v-if="getQuantity(food.id) > 0"
                    @click="decrease(food.id)"
                    class="btn-minus"
                  >-</button>
                  <span v-if="getQuantity(food.id) > 0" class="quantity">{{ getQuantity(food.id) }}</span>
                  <button @click="increase(food)" class="btn-plus">+</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部购物车栏 -->
    <div class="cart-bar">
      <div class="cart-icon" @click="showCartDetail = true">
        <span class="icon">🛒</span>
        <span v-if="totalQuantity > 0" class="badge">{{ totalQuantity }}</span>
      </div>
      <div class="cart-info">
        <div class="price">
          <span class="current">¥{{ totalAmount.toFixed(2) }}</span>
          <span class="delivery">配送费 ¥{{ merchant?.delivery_fee || 0 }}</span>
        </div>
      </div>
      <button
        :class="['submit-btn', { disabled: totalAmount < (merchant?.min_price || 0) }]"
        @click="goToCart"
      >
        {{ totalAmount < (merchant?.min_price || 0) ? `¥${merchant?.min_price}起送` : '去结算' }}
      </button>
    </div>

    <!-- 购物车详情弹窗 -->
    <div v-if="showCartDetail" class="cart-modal" @click="showCartDetail = false">
      <div class="cart-detail" @click.stop>
        <div class="cart-header">
          <span>已选商品</span>
          <span class="clear" @click="clearCart">清空</span>
        </div>
        <div class="cart-items">
          <div v-for="item in cartItems" :key="item.id" class="cart-item">
            <span class="name">{{ item.food_name }}</span>
            <span class="price">¥{{ item.price }}</span>
            <div class="quantity-control">
              <button @click="decrease(item.food_id)">-</button>
              <span>{{ item.quantity }}</span>
              <button @click="increase(foods.find(f => f.id === item.food_id))">+</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { getMerchantDetail } from '../api'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()

const merchant = ref(null)
const categories = ref([])
const foods = ref([])
const activeCategory = ref(null)
const showCartDetail = ref(false)

const cartItems = computed(() => cartStore.items)
const totalQuantity = computed(() => cartStore.totalQuantity)
const totalAmount = computed(() => cartStore.totalAmount)

const getQuantity = (foodId) => {
  const item = cartStore.items.find(i => i.food_id === foodId)
  return item?.quantity || 0
}

const increase = async (food) => {
  await cartStore.addItem({
    merchant_id: parseInt(route.params.id),
    food_id: food.id,
    quantity: 1
  })
}

const decrease = async (foodId) => {
  const item = cartStore.items.find(i => i.food_id === foodId)
  if (item) {
    if (item.quantity > 1) {
      await cartStore.updateQuantity(item.id, item.quantity - 1)
    } else {
      await cartStore.removeItem(item.id)
    }
  }
}

const clearCart = () => {
  cartStore.clear()
  showCartDetail.value = false
}

const scrollToCategory = (id) => {
  activeCategory.value = id
  const el = document.getElementById('cat-' + id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' })
  }
}

const goToCart = () => {
  if (totalAmount.value >= (merchant.value?.min_price || 0)) {
    router.push('/cart')
  }
}

const fetchMerchant = async () => {
  try {
    const res = await getMerchantDetail(route.params.id)
    if (res.code === 0) {
      merchant.value = res.data
      categories.value = res.data.categories || []
      foods.value = categories.value.flatMap(c => c.foods || [])
      if (categories.value.length > 0) {
        activeCategory.value = categories.value[0].id
      }
    }
  } catch (err) {
    console.error('获取商家详情失败', err)
  }
}

onMounted(() => {
  fetchMerchant()
  cartStore.fetchCart()
})
</script>

<style scoped>
.merchant-page {
  padding-bottom: 60px;
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

.merchant-info {
  display: flex;
  padding: 15px;
  background: white;
  margin-bottom: 10px;
}

.merchant-info .logo {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  margin-right: 15px;
}

.merchant-info h1 {
  font-size: 18px;
  margin-bottom: 8px;
}

.rating {
  display: flex;
  gap: 15px;
  font-size: 13px;
  color: #666;
}

.notice {
  font-size: 12px;
  color: #ff6b6b;
  margin-top: 5px;
}

.delivery-info {
  display: flex;
  gap: 15px;
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.menu {
  display: flex;
  background: white;
  height: calc(100vh - 200px);
}

.category-sidebar {
  width: 80px;
  background: #f5f5f5;
  overflow-y: auto;
}

.category-item {
  padding: 15px 10px;
  font-size: 13px;
  text-align: center;
  cursor: pointer;
}

.category-item.active {
  background: white;
  color: #ff6b6b;
  border-left: 3px solid #ff6b6b;
}

.food-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.category-section h3 {
  font-size: 14px;
  color: #666;
  padding: 10px 0;
}

.food-item {
  display: flex;
  padding: 10px 0;
  border-bottom: 1px solid #f5f5f5;
}

.food-item img {
  width: 90px;
  height: 90px;
  border-radius: 8px;
  object-fit: cover;
  margin-right: 10px;
}

.food-info {
  flex: 1;
}

.food-info h4 {
  font-size: 15px;
  margin-bottom: 5px;
}

.desc {
  font-size: 12px;
  color: #999;
  margin-bottom: 5px;
}

.sales {
  font-size: 12px;
  color: #999;
}

.food-price {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.price {
  font-size: 16px;
  color: #ff6b6b;
  font-weight: bold;
}

.old-price {
  font-size: 12px;
  color: #999;
  text-decoration: line-through;
  margin-left: 8px;
}

.quantity-control {
  display: flex;
  align-items: center;
  margin-left: auto;
  gap: 10px;
}

.quantity-control button {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: none;
  font-size: 16px;
  cursor: pointer;
}

.btn-minus {
  background: #f5f5f5;
  color: #666;
}

.btn-plus {
  background: #ff6b6b;
  color: white;
}

.quantity {
  font-size: 14px;
  min-width: 20px;
  text-align: center;
}

.cart-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #333;
  color: white;
  display: flex;
  align-items: center;
  padding: 10px 15px;
}

.cart-icon {
  position: relative;
  margin-right: 15px;
}

.cart-icon .icon {
  font-size: 28px;
}

.badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background: #ff6b6b;
  color: white;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 10px;
}

.cart-info {
  flex: 1;
}

.cart-info .current {
  font-size: 20px;
  font-weight: bold;
}

.cart-info .delivery {
  font-size: 12px;
  color: #999;
  margin-left: 8px;
}

.submit-btn {
  background: #ff6b6b;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 20px;
  font-size: 14px;
}

.submit-btn.disabled {
  background: #666;
}

.cart-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 60px;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: flex-end;
}

.cart-detail {
  width: 100%;
  background: white;
  border-radius: 15px 15px 0 0;
  max-height: 60%;
  overflow-y: auto;
}

.cart-header {
  display: flex;
  justify-content: space-between;
  padding: 15px;
  border-bottom: 1px solid #eee;
}

.cart-header .clear {
  color: #999;
}

.cart-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #f5f5f5;
}

.cart-item .name {
  flex: 1;
}

.cart-item .price {
  color: #ff6b6b;
  margin-right: 15px;
}
</style>
