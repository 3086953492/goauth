import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import authRoutes from './modules/auth'
import userRoutes from './modules/user'
import oauthRoutes from './modules/oauth'
import systemRoutes from './modules/system'

const routes: RouteRecordRaw[] = [
  ...authRoutes,
  ...userRoutes,
  ...oauthRoutes,
  ...systemRoutes
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
  
  // 使用新的 isAuthenticated 计算属性
  const isAuthenticated = authStore.isAuthenticated

  // 需要认证的页面
  if (to.meta.requiresAuth && !isAuthenticated) {
    // 用户未登录，清理状态并重定向到登录页
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

