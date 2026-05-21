import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getCart, addToCart as addToCartApi, updateCartQuantity, removeFromCart, clearCart } from '../api'

export const useCartStore = defineStore('cart', () => {
  const cartData = ref(null)
  const loading = ref(false)

  const items = computed(() => cartData.value?.items || [])
  const merchantId = computed(() => cartData.value?.merchant_id)
  const merchantName = computed(() => cartData.value?.merchant?.name || '')
  const deliveryFee = computed(() => cartData.value?.merchant?.delivery_fee || 0)
  const minPrice = computed(() => cartData.value?.merchant?.min_price || 0)
  const totalAmount = computed(() => cartData.value?.total_price || 0)
  const totalQuantity = computed(() => cartData.value?.total_quantity || 0)
  const finalAmount = computed(() => totalAmount.value + deliveryFee.value)

  const canCheckout = computed(() => {
    return items.value.length > 0 && totalAmount.value >= minPrice.value
  })

  const fetchCart = async () => {
    loading.value = true
    try {
      const res = await getCart()
      if (res.code === 0) {
        cartData.value = res.data
      } else {
        cartData.value = null
      }
    } catch {
      cartData.value = null
    } finally {
      loading.value = false
    }
  }

  const addItem = async (data) => {
    const res = await addToCartApi(data)
    if (res.code === 0) {
      await fetchCart()
    }
    return res
  }

  const updateQuantity = async (id, quantity) => {
    const res = await updateCartQuantity(id, { quantity })
    if (res.code === 0) {
      await fetchCart()
    }
  }

  const removeItem = async (id) => {
    const res = await removeFromCart(id)
    if (res.code === 0) {
      await fetchCart()
    }
  }

  const clear = async () => {
    try {
      await clearCart()
    } catch {
      // ignore
    }
    cartData.value = null
  }

  return {
    cartData,
    loading,
    items,
    merchantId,
    merchantName,
    deliveryFee,
    minPrice,
    totalAmount,
    totalQuantity,
    finalAmount,
    canCheckout,
    fetchCart,
    addItem,
    updateQuantity,
    removeItem,
    clear
  }
})
