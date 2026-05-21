<template>
  <div class="login-page">
    <div class="login-box">
      <div class="logo">🍔 外卖系统</div>
      <h2>欢迎回来</h2>
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
            placeholder="密码"
            required
          />
        </div>
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
      <p class="link">
        还没有账号？<router-link to="/register">立即注册</router-link>
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
  password: ''
})
const loading = ref(false)
const error = ref('')

const handleSubmit = async () => {
  if (!form.value.phone || !form.value.password) {
    alert('请填写完整信息')
    return
  }

  loading.value = true
  try {
    await userStore.loginUser(form.value)
    router.push('/')
  } catch (err) {
    alert(err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
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
  border-color: #667eea;
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
  color: #667eea;
  text-decoration: none;
}
</style>
