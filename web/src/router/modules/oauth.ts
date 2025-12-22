import type { RouteRecordRaw } from 'vue-router'

const oauthRoutes: RouteRecordRaw[] = [
  {
    path: '/oauth/clients',
    name: 'OAuthClients',
    component: () => import('@/views/oauth/Clients.vue'),
    meta: {
      title: 'OAuth 客户端管理',
      requiresAuth: true
    }
  },
  {
    path: '/oauth/authorize',
    name: 'OAuthAuthorize',
    component: () => import('@/views/oauth/Authorize.vue'),
    meta: {
      title: 'OAuth 授权',
      requiresAuth: false  // 页面内部处理登录检查
    }
  }
]

export default oauthRoutes

