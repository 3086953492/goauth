<template>
  <nav class="navbar">
    <div class="navbar-container">
      <!-- 左侧 Logo/标题 -->
      <div class="navbar-brand">
        <h1 class="brand-title">GoAuth</h1>
        <span class="brand-subtitle">用户管理系统</span>
      </div>

      <!-- 右侧用户区域 -->
      <div class="navbar-user">
        <div class="user-info">
          <el-avatar :size="40" :src="user?.avatar || defaultAvatar" class="user-avatar" />
          <span class="user-nickname">{{ user?.nickname || '用户' }}</span>
        </div>

        <div class="navbar-actions">
          <el-button
            v-if="user?.role === 'admin'"
            type="info"
            :icon="User"
            @click="goToUsers"
            class="nav-button"
          >
            用户列表
          </el-button>
          <el-button type="primary" :icon="Edit" @click="goToProfile" class="nav-button">
            个人中心
          </el-button>
          <el-button type="danger" :icon="SwitchButton" @click="handleLogout" class="nav-button">
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
import { Edit, SwitchButton, User } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { DEFAULT_AVATAR } from '@/constants'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const defaultAvatar = DEFAULT_AVATAR

const goToProfile = () => {
  router.push('/profile')
}

const goToUsers = () => {
  router.push('/users')
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
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
  height: 64px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  z-index: 1000;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.navbar-container {
  max-width: 1400px;
  height: 100%;
  margin: 0 auto;
  padding: 0 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.navbar-brand {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.brand-title {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.5px;
}

.brand-subtitle {
  font-size: 13px;
  color: #909399;
  font-weight: 500;
}

.navbar-user {
  display: flex;
  align-items: center;
  gap: 24px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  border: 2px solid #f5f7fa;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.user-nickname {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.navbar-actions {
  display: flex;
  gap: 12px;
}

.nav-button {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.nav-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .navbar-container {
    padding: 0 16px;
  }

  .brand-subtitle {
    display: none;
  }

  .user-nickname {
    display: none;
  }

  .navbar-actions {
    gap: 8px;
  }

  .nav-button {
    padding: 8px 12px;
    font-size: 14px;
  }

  .nav-button :deep(.el-icon) {
    margin-right: 0;
  }

  .nav-button :deep(span) {
    display: none;
  }
}

@media (max-width: 480px) {
  .navbar {
    height: 56px;
  }

  .brand-title {
    font-size: 20px;
  }

  .user-avatar {
    width: 32px;
    height: 32px;
  }
}
</style>
