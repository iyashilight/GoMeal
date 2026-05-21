<template>
  <div class="profile-page">
    <header class="header">
      <button class="back" @click="$router.push('/')">←</button>
      <span>个人中心</span>
      <span></span>
    </header>

    <div class="user-info">
      <div class="avatar">{{ userStore.userInfo?.nickname?.[0] || '用' }}</div>
      <div class="info">
        <h3>{{ userStore.userInfo?.nickname || '用户' }}</h3>
        <p>{{ userStore.userInfo?.phone }}</p>
      </div>
    </div>

    <div class="menu-list">
      <router-link to="/order-list" class="menu-item">
        <span class="icon">📋</span>
        <span class="text">我的订单</span>
        <span class="arrow">→</span>
      </router-link>
      <router-link to="/address" class="menu-item">
        <span class="icon">📍</span>
        <span class="text">收货地址</span>
        <span class="arrow">→</span>
      </router-link>
      <router-link v-if="userStore.userInfo?.merchant_id" to="/merchant/manage" class="menu-item">
        <span class="icon">🏪</span>
        <span class="text">商家管理</span>
        <span class="arrow">→</span>
      </router-link>
      <router-link v-if="!userStore.userInfo?.merchant_id" to="/merchant/register" class="menu-item">
        <span class="icon">📝</span>
        <span class="text">商家注册</span>
        <span class="arrow">→</span>
      </router-link>
    </div>

    <button class="logout-btn" @click="logout">退出登录</button>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const logout = () => {
  if (confirm('确定要退出登录吗？')) {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: #f5f5f5;
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

.user-info {
  display: flex;
  align-items: center;
  padding: 30px 20px;
  background: linear-gradient(135deg, #ff6b6b 0%, #ff8e8e 100%);
  color: white;
}

.avatar {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  background: white;
  color: #ff6b6b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: bold;
  margin-right: 20px;
}

.info h3 {
  font-size: 20px;
  margin-bottom: 5px;
}

.info p {
  opacity: 0.9;
}

.menu-list {
  margin-top: 15px;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 18px 20px;
  background: white;
  border-bottom: 1px solid #f5f5f5;
  text-decoration: none;
  color: #333;
}

.menu-item .icon {
  font-size: 20px;
  margin-right: 15px;
}

.menu-item .text {
  flex: 1;
}

.menu-item .arrow {
  color: #999;
}

.logout-btn {
  display: block;
  width: 90%;
  margin: 40px auto;
  padding: 15px;
  background: white;
  border: none;
  border-radius: 8px;
  color: #ff6b6b;
  font-size: 16px;
  cursor: pointer;
}
</style>
