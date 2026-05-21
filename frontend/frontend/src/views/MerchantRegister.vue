<template>
  <div class="register-page">
    <header class="header">
      <button class="back" @click="$router.back()">←</button>
      <span>商家注册</span>
      <span></span>
    </header>

    <div class="form-wrap">
      <div class="form-card">
        <h3>店铺信息</h3>

        <div class="form-item">
          <label>店铺名称</label>
          <input v-model="form.name" placeholder="请输入店铺名称" maxlength="100" />
        </div>

        <div class="form-item">
          <label>联系电话</label>
          <input v-model="form.phone" placeholder="请输入联系电话" />
        </div>

        <div class="form-item">
          <label>店铺地址</label>
          <input v-model="form.address" placeholder="请输入店铺地址" maxlength="255" />
        </div>

        <div class="form-item">
          <label>店铺公告</label>
          <textarea v-model="form.notice" placeholder="请输入店铺公告（选填）" maxlength="500" rows="3"></textarea>
        </div>

        <div class="form-row">
          <div class="form-item half">
            <label>起送价 (¥)</label>
            <input v-model.number="form.min_price" type="number" step="0.01" min="0" placeholder="0" />
          </div>
          <div class="form-item half">
            <label>配送费 (¥)</label>
            <input v-model.number="form.delivery_fee" type="number" step="0.01" min="0" placeholder="0" />
          </div>
        </div>

        <button class="submit-btn" @click="handleRegister" :disabled="submitting">
          {{ submitting ? '提交中...' : '提交注册' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { merchantRegister } from '../api'

const router = useRouter()
const userStore = useUserStore()

const submitting = ref(false)

const form = reactive({
  name: '',
  phone: '',
  address: '',
  notice: '',
  min_price: 0,
  delivery_fee: 0
})

const handleRegister = async () => {
  if (!form.name.trim()) {
    alert('请输入店铺名称')
    return
  }
  if (!form.phone.trim()) {
    alert('请输入联系电话')
    return
  }
  if (!form.address.trim()) {
    alert('请输入店铺地址')
    return
  }

  submitting.value = true
  try {
    const res = await merchantRegister(form)
    if (res.code === 0) {
      // 更新用户信息中的 merchant_id
      await userStore.fetchUserInfo()
      alert('注册成功！')
      router.push('/')
    }
  } catch (err) {
    alert(err || '注册失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.register-page {
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

.form-wrap {
  padding: 15px;
}

.form-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
}

.form-card h3 {
  font-size: 18px;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #f5f5f5;
}

.form-item {
  margin-bottom: 18px;
}

.form-item label {
  display: block;
  font-size: 14px;
  color: #666;
  margin-bottom: 6px;
  font-weight: 500;
}

.form-item input,
.form-item textarea {
  width: 100%;
  padding: 12px 15px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-item input:focus,
.form-item textarea:focus {
  border-color: #ff6b6b;
}

.form-item textarea {
  resize: vertical;
  font-family: inherit;
}

.form-row {
  display: flex;
  gap: 12px;
}

.form-row .half {
  flex: 1;
}

.submit-btn {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #ff6b6b 0%, #ff8e8e 100%);
  color: white;
  border: none;
  border-radius: 25px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  margin-top: 10px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 107, 107, 0.4);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}
</style>
