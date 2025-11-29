import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import router from './router'
import App from './App.vue'
import './styles/reset.css'
import './styles/global.css'
import './styles/variables.css'
import { useAuthStore } from './stores/useAuthStore'
import { startTokenRefresh } from './composables/useTokenRefresh'

const app = createApp(App)
const pinia = createPinia()

// 注册所有 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(pinia)

// 在路由挂载前初始化认证状态（从 sessionStorage 恢复）
const authStore = useAuthStore()
const hasValidSession = authStore.initFromSession()

// 如果恢复了有效的登录态，启动令牌自动刷新
if (hasValidSession) {
  startTokenRefresh()
}

app.use(router)
app.use(ElementPlus)

app.mount('#app')
