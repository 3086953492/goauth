import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: {
      title: '用户注册',
      requiresAuth: false
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: {
      title: '用户登录',
      requiresAuth: false
    }
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: {
      title: '主页',
      requiresAuth: true
    }
  },
  {
    path: '/profile/:id?',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: {
      title: '个人信息',
      requiresAuth: true
    }
  },
  {
    path: '/users',
    name: 'Users',
    component: () => import('@/views/Users.vue'),
    meta: {
      title: '用户列表',
      requiresAuth: true
    }
  },
  {
    path: '/oauth-clients',
    name: 'OAuthClients',
    component: () => import('@/views/OAuthClients.vue'),
    meta: {
      title: 'OAuth 客户端管理',
      requiresAuth: true
    }
  },
  {
    path: '/oauth/authorize',
    name: 'OAuthAuthorize',
    component: () => import('@/views/OAuthAuthorize.vue'),
    meta: {
      title: 'OAuth 授权',
      requiresAuth: false  // 页面内部处理登录检查
    }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = to.meta.title as string
  }

  const authStore = useAuthStore()
  
  // 使用新的 isAuthenticated 计算属性（包含过期判断）
  const isAuthenticated = authStore.isAuthenticated

  // 需要认证的页面
  if (to.meta.requiresAuth && !isAuthenticated) {
    // token 不存在或已过期，清理状态并重定向到登录页
    authStore.logout()
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
    return
  }

  // 已登录用户访问登录或注册页面，重定向到主页
  if ((to.path === '/login' || to.path === '/register') && isAuthenticated) {
    next('/home')
    return
  }

  next()
})

export default router

