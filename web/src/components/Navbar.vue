<template>
  <nav class="navbar">
    <div class="navbar__container">
      <!-- 左侧 Logo/标题 -->
      <div class="navbar__brand">
        <h1 class="navbar__title">GoAuth</h1>
        <span class="navbar__subtitle">用户管理系统</span>
      </div>

      <!-- 右侧用户区域 -->
      <div class="navbar__user">
        <div class="navbar__user-info">
          <el-avatar :size="avatarSize" :src="user?.avatar" :icon="Avatar" class="navbar__avatar" />
          <span class="navbar__nickname">{{ user?.nickname || '用户' }}</span>
        </div>

        <div class="navbar__actions">
          <el-button v-if="user?.role === 'admin'" type="info" :icon="User" @click="goToUsers" class="navbar__button">
            用户列表
          </el-button>
          <el-button v-if="user?.role === 'admin'" type="warning" :icon="Key" @click="goToOAuthClients"
            class="navbar__button">
            OAuth 客户端
          </el-button>
          <el-button type="primary" :icon="Edit" @click="goToProfile" class="navbar__button">
            个人中心
          </el-button>
          <el-button type="danger" :icon="SwitchButton" @click="handleLogout" class="navbar__button">
            退出登录
          </el-button>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Edit, SwitchButton, User, Key, Avatar } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const authStore = useAuthStore()
const { handleLogout: logout } = useAuth()

// 头像尺寸
const avatarSize = 40

const user = computed(() => authStore.user)

const goToProfile = () => {
  router.push('/profile')
}

const goToUsers = () => {
  router.push('/users')
}

const goToOAuthClients = () => {
  router.push('/oauth-clients')
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    // 执行登出业务逻辑
    const result = await logout()
    
    if (result.success) {
      // 登出成功：显示提示并跳转到登录页
      ElMessage.success(result.message || '已退出登录')
      router.push('/login')
    }
  } catch {
    // 用户取消操作
  }
}
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: var(--navbar-height);
  background: var(--color-navbar-background);
  backdrop-filter: blur(10px);
  box-shadow: var(--shadow-navbar);
  z-index: 1000;
  border-bottom: var(--border-width-thin) solid var(--color-navbar-border);
}

.navbar__container {
  max-width: var(--container-max-width-xlarge);
  height: 100%;
  margin: 0 auto;
  padding: 0 var(--spacing-lg);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.navbar__brand {
  display: flex;
  align-items: baseline;
  gap: var(--spacing-sm-lg);
}

.navbar__title {
  margin: 0;
  font-size: var(--font-size-title);
  font-weight: 700;
  background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.5px;
}

.navbar__subtitle {
  font-size: var(--font-size-xxs);
  color: var(--color-text-tertiary);
  font-weight: 500;
}

.navbar__user {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
}

.navbar__user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm-lg);
}

.navbar__avatar {
  border: var(--border-width-medium) solid var(--color-background-lighter);
  box-shadow: var(--shadow-card);
}

.navbar__nickname {
  font-size: var(--font-size-base);
  font-weight: 600;
  color: var(--color-text-primary);
}

.navbar__actions {
  display: flex;
  gap: var(--spacing-sm-lg);
}

.navbar__button {
  border-radius: var(--border-radius-large);
  font-weight: 500;
  transition: all 0.3s ease;
}

.navbar__button:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-button-hover);
}

/* 响应式设计 */
/* 平板端：对应 --breakpoint-tablet (768px) */
@media (max-width: 768px) {
  .navbar__container {
    padding: 0 var(--spacing-md);
  }

  .navbar__subtitle {
    display: none;
  }

  .navbar__nickname {
    display: none;
  }

  .navbar__actions {
    gap: var(--spacing-sm);
  }

  .navbar__button {
    padding: var(--spacing-sm) var(--spacing-sm-lg);
    font-size: var(--font-size-sm);
  }

  .navbar__button :deep(.el-icon) {
    margin-right: 0;
  }

  .navbar__button :deep(span) {
    display: none;
  }
}

/* 移动端：对应 --breakpoint-mobile (480px) */
@media (max-width: 480px) {
  .navbar {
    height: var(--navbar-height-mobile);
  }

  .navbar__title {
    font-size: var(--font-size-lg);
  }

  .navbar__avatar {
    width: 32px;
    height: 32px;
  }
}
</style>
