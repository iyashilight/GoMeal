<template>
  <div class="register-page">
    <div class="register-box">
      <div class="logo">🍔 外卖系统</div>
      <h2>注册账号</h2>
      <form @submit.prevent="handleSubmit">
        <div class="form-item">
          <input
            v-model="form.phone"
            type="tel"
            placeholder="手机号"
            maxlength="11"
            required
          />
        </div>
        <div class="form-item">
          <input
            v-model="form.password"
            type="password"
            placeholder="密码（至少6位）"
            minlength="6"
            required
          />
        </div>
        <div class="form-item">
          <input
            v-model="form.confirmPassword"
            type="password"
            placeholder="确认密码"
            required
          />
        </div>
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>
      <p class="link">
        已有账号？<router-link to="/login">立即登录</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  phone: '',
  password: '',
  confirmPassword: ''
})
const loading = ref(false)

const handleSubmit = async () => {
  if (form.value.password !== form.value.confirmPassword) {
    alert('两次输入的密码不一致')
    return
  }

  if (form.value.password.length < 6) {
    alert('密码长度至少6位')
    return
  }

  loading.value = true
  try {
    await userStore.registerUser({
      phone: form.value.phone,
      password: form.value.password
    })
    router.push('/')
  } catch (err) {
    alert(err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.register-box {
  background: white;
  padding: 40px;
  border-radius: 16px;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 10px 40px rgba(0,0,0,0.2);
}

.logo {
  font-size: 48px;
  text-align: center;
  margin-bottom: 10px;
}

h2 {
  text-align: center;
  color: #333;
  margin-bottom: 30px;
}

.form-item {
  margin-bottom: 20px;
}

.form-item input {
  width: 100%;
  padding: 14px 16px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.form-item input:focus {
  outline: none;
  border-color: #f5576c;
}

button {
  width: 100%;
  padding: 14px;
  font-size: 16px;
}

.link {
  text-align: center;
  margin-top: 20px;
  color: #666;
}

.link a {
  color: #f5576c;
  text-decoration: none;
}
</style>
