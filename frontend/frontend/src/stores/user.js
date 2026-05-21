import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, register, getUserInfo } from '../api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)

  const isLoggedIn = computed(() => !!token.value)

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const loginUser = async (data) => {
    const res = await login(data)
    if (res.code === 0) {
      setToken(res.data.token)
      userInfo.value = res.data.user
      return res.data
    }
    throw new Error(res.message)
  }

  const registerUser = async (data) => {
    const res = await register(data)
    if (res.code === 0) {
      setToken(res.data.token)
      userInfo.value = res.data.user
      return res.data
    }
    throw new Error(res.message)
  }

  const fetchUserInfo = async () => {
    if (!token.value) return
    try {
      const res = await getUserInfo()
      if (res.code === 0) {
        userInfo.value = res.data
      }
    } catch {
      // 401 is handled by axios interceptor (redirect to login)
      // other errors are transient network issues, don't logout
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    loginUser,
    registerUser,
    fetchUserInfo,
    logout
  }
})
