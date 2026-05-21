import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { public: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/merchant/:id',
    name: 'Merchant',
    component: () => import('../views/Merchant.vue')
  },
  {
    path: '/cart',
    name: 'Cart',
    component: () => import('../views/Cart.vue')
  },
  {
    path: '/order/:id',
    name: 'Order',
    component: () => import('../views/Order.vue')
  },
  {
    path: '/order-list',
    name: 'OrderList',
    component: () => import('../views/OrderList.vue')
  },
  {
    path: '/address',
    name: 'Address',
    component: () => import('../views/Address.vue')
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue')
  },
  {
    path: '/merchant/register',
    name: 'MerchantRegister',
    component: () => import('../views/MerchantRegister.vue')
  },
  {
    path: '/merchant/manage',
    name: 'MerchantManage',
    component: () => import('../views/MerchantManage.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫 - 直接从 localStorage 读取 token，避免 Pinia 未就绪问题
function getToken() {
  try {
    return localStorage.getItem('token') || ''
  } catch {
    return ''
  }
}

router.beforeEach((to, from, next) => {
  const token = getToken()
  if (!to.meta.public && !token) {
    next('/login')
  } else if ((to.path === '/login' || to.path === '/register') && token) {
    next('/')
  } else {
    next()
  }
})

export default router
