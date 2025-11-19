import type { RouteRecordRaw } from 'vue-router'

const systemRoutes: RouteRecordRaw[] = [
  {
    path: '/error',
    name: 'Error',
    component: () => import('@/views/ErrorPage.vue'),
    meta: {
      title: '错误',
      requiresAuth: false
    }
  }
]

export default systemRoutes

