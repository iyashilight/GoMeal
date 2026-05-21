import axios from 'axios'
import { useUserStore } from '../stores/user'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

api.interceptors.request.use(config => {
  let token
  try {
    token = useUserStore().token
  } catch {
    token = localStorage.getItem('token') || ''
  }
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  response => response.data,
  error => {
    // 401: token无效，需要重新登录
    if (error.response?.status === 401) {
      try {
        useUserStore().logout()
      } catch {
        localStorage.removeItem('token')
      }
      window.location.href = '/login'
      return Promise.reject('登录已过期，请重新登录')
    }
    // 其他错误：优先使用后端返回的消息
    const msg = error.response?.data?.message
      || error.response?.data?.error
      || (typeof error.response?.data === 'string' ? error.response.data : null)
      || '网络错误，请稍后重试'
    return Promise.reject(msg)
  }
)

export default api

// 用户相关
export const login = (data) => api.post('/auth/login', data)
export const register = (data) => api.post('/auth/register', data)
export const getUserInfo = () => api.get('/user/info')

// 商家相关
export const getMerchantList = () => api.get('/merchants')
export const getMerchantDetail = (id) => api.get(`/merchants/${id}`)
export const getFoodDetail = (id) => api.get(`/foods/${id}`)

// 购物车相关
export const getCart = () => api.get('/cart')
export const addToCart = (data) => api.post('/cart', data)
export const updateCartQuantity = (id, data) => api.put(`/cart/${id}`, data)
export const removeFromCart = (id) => api.delete(`/cart/${id}`)
export const clearCart = () => api.post('/cart/clear')

// 地址相关
export const getAddressList = () => api.get('/addresses')
export const createAddress = (data) => api.post('/addresses', data)
export const updateAddress = (id, data) => api.put(`/addresses/${id}`, data)
export const deleteAddress = (id) => api.delete(`/addresses/${id}`)

// 订单相关
export const createOrder = (data) => api.post('/orders', data)
export const getOrderList = (params) => api.get('/orders', { params })
export const getOrderDetail = (id) => api.get(`/orders/${id}`)
export const cancelOrder = (id) => api.put(`/orders/${id}/cancel`)

// 支付相关
export const payOrder = (id) => api.post(`/orders/${id}/pay`)
export const getPayment = (id) => api.get(`/orders/${id}/payment`)

// 商家注册
export const merchantRegister = (data) => api.post('/merchant/register', data)

// 商家商品管理
export const getMyFoods = () => api.get('/merchant/foods')
export const createFood = (data) => api.post('/merchant/foods', data)
export const updateFood = (id, data) => api.put(`/merchant/foods/${id}`, data)
export const deleteFood = (id) => api.delete(`/merchant/foods/${id}`)
export const setFoodStatus = (id, data) => api.put(`/merchant/foods/${id}/status`, data)

// 初始化数据
export const initData = () => api.post('/init')
