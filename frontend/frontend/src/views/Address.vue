<template>
  <div class="address-page">
    <header class="header">
      <button class="back" @click="$router.back()">←</button>
      <span>收货地址</span>
      <span></span>
    </header>

    <div class="address-list">
      <div
        v-for="addr in addresses"
        :key="addr.id"
        :class="['address-card', { default: addr.is_default }]"
        @click="selectAddress(addr)"
      >
        <div class="addr-header">
          <span class="name">{{ addr.name }}</span>
          <span class="phone">{{ addr.phone }}</span>
          <span v-if="addr.is_default" class="default-tag">默认</span>
          <span v-if="addr.tag" class="tag">{{ addr.tag }}</span>
        </div>
        <div class="addr-detail">{{ addr.address }}</div>
        <div class="actions">
          <button @click.stop="editAddress(addr)">编辑</button>
          <button @click.stop="deleteAddr(addr.id)">删除</button>
        </div>
      </div>
    </div>

    <div v-if="addresses.length === 0" class="empty">
      <p>还没有添加地址</p>
    </div>

    <button class="add-btn" @click="openAddForm">+ 新增地址</button>

    <!-- 地址表单弹窗 -->
    <div v-if="showForm" class="modal" @click="showForm = false">
      <div class="form-container" @click.stop>
        <h3>{{ editingId ? '编辑地址' : '新增地址' }}</h3>
        <div class="form-item">
          <input v-model="form.name" placeholder="联系人姓名" />
        </div>
        <div class="form-item">
          <input v-model="form.phone" placeholder="手机号" maxlength="11" />
        </div>
        <div class="form-item">
          <input v-model="form.address" placeholder="详细地址" />
        </div>
        <div class="form-item">
          <select v-model="form.tag">
            <option value="">标签（选填）</option>
            <option value="家">家</option>
            <option value="公司">公司</option>
            <option value="学校">学校</option>
          </select>
        </div>
        <div class="form-item checkbox">
          <label>
            <input v-model="form.is_default" type="checkbox" />
            设为默认地址
          </label>
        </div>
        <div class="form-actions">
          <button class="cancel" @click="showForm = false">取消</button>
          <button class="submit" @click="submit">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getAddressList, createAddress, updateAddress, deleteAddress } from '../api'

const route = useRoute()
const router = useRouter()

const addresses = ref([])
const showForm = ref(false)
const editingId = ref(null)
const form = ref({
  name: '',
  phone: '',
  address: '',
  tag: '',
  is_default: false
})

const fetchAddresses = async () => {
  try {
    const res = await getAddressList()
    if (res.code === 0) {
      addresses.value = res.data || []
    }
  } catch (err) {
    console.error('获取地址失败', err)
  }
}

const selectAddress = (addr) => {
  if (route.query.from === 'cart') {
    router.back()
  }
}

const editAddress = (addr) => {
  editingId.value = addr.id
  form.value = { ...addr }
  showForm.value = true
}

const deleteAddr = async (id) => {
  if (!confirm('确定要删除这个地址吗？')) return
  try {
    const res = await deleteAddress(id)
    if (res.code === 0) {
      fetchAddresses()
    }
  } catch (err) {
    alert(err)
  }
}

const openAddForm = () => {
  editingId.value = null
  form.value = {
    name: '',
    phone: '',
    address: '',
    tag: '',
    is_default: false
  }
  showForm.value = true
}

const submit = async () => {
  if (!form.value.name || !form.value.phone || !form.value.address) {
    alert('请填写完整信息')
    return
  }

  try {
    let res
    if (editingId.value) {
      res = await updateAddress(editingId.value, form.value)
    } else {
      res = await createAddress(form.value)
    }
    if (res.code !== 0) {
      alert(res.message || '操作失败')
      return
    }
    showForm.value = false
    resetForm()
    fetchAddresses()
  } catch (err) {
    alert(err)
  }
}

const resetForm = () => {
  editingId.value = null
  form.value = {
    name: '',
    phone: '',
    address: '',
    tag: '',
    is_default: false
  }
}

onMounted(() => {
  fetchAddresses()
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

.address-list {
  padding: 10px;
}

.address-card {
  background: white;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 10px;
  cursor: pointer;
}

.address-card.default {
  border: 1px solid #ff6b6b;
}

.addr-header {
  margin-bottom: 8px;
}

.addr-header .name {
  font-weight: bold;
  margin-right: 10px;
}

.addr-header .phone {
  color: #666;
}

.default-tag {
  background: #ff6b6b;
  color: white;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 8px;
}

.tag {
  background: #f0f0f0;
  color: #666;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 8px;
}

.addr-detail {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
}

.actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.actions button {
  padding: 5px 12px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

.empty {
  text-align: center;
  padding: 50px;
  color: #999;
}

.add-btn {
  position: fixed;
  bottom: 20px;
  left: 20px;
  right: 20px;
  background: #ff6b6b;
  color: white;
  border: none;
  padding: 15px;
  border-radius: 25px;
  font-size: 16px;
  cursor: pointer;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.form-container {
  background: white;
  padding: 25px;
  border-radius: 16px;
  width: 90%;
  max-width: 400px;
}

.form-container h3 {
  text-align: center;
  margin-bottom: 20px;
}

.form-item {
  margin-bottom: 15px;
}

.form-item input,
.form-item select {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 15px;
}

.form-item.checkbox {
  display: flex;
  align-items: center;
}

.form-item.checkbox input {
  width: auto;
  margin-right: 8px;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 20px;
}

.form-actions button {
  flex: 1;
  padding: 12px;
  border-radius: 8px;
  border: none;
  font-size: 15px;
  cursor: pointer;
}

.form-actions .cancel {
  background: #f5f5f5;
  color: #666;
}

.form-actions .submit {
  background: #ff6b6b;
  color: white;
}
</style>
