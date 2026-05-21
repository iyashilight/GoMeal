<template>
  <div class="manage-page">
    <header class="header">
      <button class="back" @click="$router.push('/')">←</button>
      <span>商家管理</span>
      <button class="add-btn" @click="showForm = true; editingFood = null">+ 新增</button>
    </header>

    <!-- 商品列表 -->
    <div class="food-list">
      <div v-if="foods.length === 0" class="empty">
        <p>还没有商品</p>
        <button class="btn-primary" @click="showForm = true; editingFood = null">添加第一个商品</button>
      </div>

      <div v-for="food in foods" :key="food.id" class="food-card">
        <img :src="food.image" :alt="food.name" class="food-img" />
        <div class="food-info">
          <h4>{{ food.name }}</h4>
          <p class="desc">{{ food.description }}</p>
          <div class="meta">
            <span class="price">¥{{ food.price }}</span>
            <span v-if="food.old_price" class="old-price">¥{{ food.old_price }}</span>
            <span class="stock">库存: {{ food.stock }}</span>
          </div>
        </div>
        <div class="food-actions">
          <label class="switch">
            <input type="checkbox" :checked="food.status === 1" @change="toggleStatus(food)" />
            <span class="slider"></span>
          </label>
          <span class="status-text">{{ food.status === 1 ? '上架' : '下架' }}</span>
          <button class="edit-btn" @click="startEdit(food)">编辑</button>
          <button class="del-btn" @click="confirmDelete(food)">删除</button>
        </div>
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <div v-if="showForm" class="modal" @click="showForm = false">
      <div class="modal-content" @click.stop>
        <h3>{{ editingFood ? '编辑商品' : '新增商品' }}</h3>

        <div class="form-item">
          <label>商品名称</label>
          <input v-model="form.name" placeholder="请输入商品名称" />
        </div>

        <div class="form-item">
          <label>商品描述</label>
          <input v-model="form.description" placeholder="请输入商品描述" />
        </div>

        <div class="form-item">
          <label>商品图片</label>
          <input v-model="form.image" placeholder="请输入图片URL" />
        </div>

        <div class="form-row">
          <div class="form-item half">
            <label>价格 (¥)</label>
            <input v-model.number="form.price" type="number" step="0.01" min="0" placeholder="0" />
          </div>
          <div class="form-item half">
            <label>原价 (¥)</label>
            <input v-model.number="form.old_price" type="number" step="0.01" min="0" placeholder="0" />
          </div>
        </div>

        <div class="form-item">
          <label>库存</label>
          <input v-model.number="form.stock" type="number" min="0" placeholder="0" />
        </div>

        <div class="modal-actions">
          <button class="cancel-btn" @click="showForm = false">取消</button>
          <button class="save-btn" @click="handleSave" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getMyFoods, createFood, updateFood, deleteFood, setFoodStatus } from '../api'

const foods = ref([])
const showForm = ref(false)
const editingFood = ref(null)
const saving = ref(false)

const form = reactive({
  name: '',
  description: '',
  image: '',
  price: 0,
  old_price: 0,
  stock: 0,
  category_id: 1
})

const resetForm = () => {
  form.name = ''
  form.description = ''
  form.image = ''
  form.price = 0
  form.old_price = 0
  form.stock = 0
  form.category_id = 1
}

const fetchFoods = async () => {
  try {
    const res = await getMyFoods()
    if (res.code === 0) {
      foods.value = res.data || []
    }
  } catch (err) {
    console.error('获取商品列表失败', err)
  }
}

const startEdit = (food) => {
  editingFood.value = food
  form.name = food.name
  form.description = food.description || ''
  form.image = food.image || ''
  form.price = food.price
  form.old_price = food.old_price || 0
  form.stock = food.stock
  form.category_id = food.category_id
  showForm.value = true
}

const handleSave = async () => {
  if (!form.name.trim()) {
    alert('请输入商品名称')
    return
  }
  if (!form.price || form.price <= 0) {
    alert('请输入有效价格')
    return
  }

  saving.value = true
  try {
    if (editingFood.value) {
      const res = await updateFood(editingFood.value.id, { ...form })
      if (res.code === 0) {
        await fetchFoods()
        showForm.value = false
        resetForm()
      }
    } else {
      const res = await createFood({ ...form })
      if (res.code === 0) {
        await fetchFoods()
        showForm.value = false
        resetForm()
      }
    }
  } catch (err) {
    alert(err || '操作失败')
  } finally {
    saving.value = false
  }
}

const toggleStatus = async (food) => {
  const newStatus = food.status === 1 ? 2 : 1
  try {
    const res = await setFoodStatus(food.id, { status: newStatus })
    if (res.code === 0) {
      food.status = newStatus
    }
  } catch (err) {
    alert(err || '操作失败')
  }
}

const confirmDelete = async (food) => {
  if (!confirm(`确定要删除「${food.name}」吗？`)) return
  try {
    const res = await deleteFood(food.id)
    if (res.code === 0) {
      await fetchFoods()
    }
  } catch (err) {
    alert(err || '删除失败')
  }
}

onMounted(() => {
  fetchFoods()
})
</script>

<style scoped>
.manage-page {
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

.add-btn {
  background: #ff6b6b;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
}

.food-list {
  padding: 10px;
}

.empty {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.empty p {
  margin-bottom: 20px;
  font-size: 16px;
}

.food-card {
  display: flex;
  background: white;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  align-items: center;
}

.food-img {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  object-fit: cover;
  margin-right: 12px;
}

.food-info {
  flex: 1;
  min-width: 0;
}

.food-info h4 {
  font-size: 15px;
  margin-bottom: 4px;
}

.food-info .desc {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta {
  font-size: 13px;
}

.price {
  color: #ff6b6b;
  font-weight: bold;
}

.old-price {
  color: #999;
  text-decoration: line-through;
  margin-left: 6px;
  font-size: 12px;
}

.stock {
  color: #666;
  margin-left: 10px;
}

.food-actions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  margin-left: 10px;
  min-width: 60px;
}

.status-text {
  font-size: 11px;
  color: #999;
}

.edit-btn,
.del-btn {
  border: none;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  cursor: pointer;
  width: 100%;
}

.edit-btn {
  background: #e8f5e9;
  color: #4caf50;
}

.del-btn {
  background: #fbe9e7;
  color: #ff5722;
}

/* Toggle switch */
.switch {
  position: relative;
  display: inline-block;
  width: 36px;
  height: 20px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #ccc;
  border-radius: 20px;
  transition: 0.3s;
}

.slider:before {
  content: "";
  position: absolute;
  height: 16px;
  width: 16px;
  left: 2px;
  bottom: 2px;
  background: white;
  border-radius: 50%;
  transition: 0.3s;
}

.switch input:checked + .slider {
  background: #4caf50;
}

.switch input:checked + .slider:before {
  transform: translateX(16px);
}

/* Modal */
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
  z-index: 200;
}

.modal-content {
  background: white;
  border-radius: 16px;
  padding: 24px;
  width: 90%;
  max-width: 420px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-content h3 {
  font-size: 18px;
  margin-bottom: 20px;
}

.form-item {
  margin-bottom: 16px;
}

.form-item label {
  display: block;
  font-size: 13px;
  color: #666;
  margin-bottom: 4px;
}

.form-item input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  box-sizing: border-box;
}

.form-item input:focus {
  border-color: #ff6b6b;
}

.form-row {
  display: flex;
  gap: 12px;
}

.form-row .half {
  flex: 1;
}

.modal-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}

.cancel-btn,
.save-btn {
  flex: 1;
  padding: 12px;
  border-radius: 20px;
  font-size: 15px;
  cursor: pointer;
  border: none;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.save-btn {
  background: linear-gradient(135deg, #ff6b6b 0%, #ff8e8e 100%);
  color: white;
  font-weight: bold;
}

.save-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}
</style>
