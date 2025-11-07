<template>
  <div class="home-container">
    <el-card class="user-card">
      <div class="user-info">
        <el-avatar :size="80" :src="user?.avatar || defaultAvatar" class="avatar" />
        
        <div class="info-content">
          <h2 class="nickname">{{ user?.nickname || '未知用户' }}</h2>
          <p class="username">@{{ user?.username }}</p>
          
          <div class="user-details">
            <el-tag type="info" size="large">{{ user?.role || 'user' }}</el-tag>
            <el-tag :type="user?.status === 1 ? 'success' : 'danger'" size="large" class="status-tag">
              {{ user?.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </div>
          
          <div class="user-meta">
            <p><span class="label">用户ID:</span> {{ user?.id }}</p>
            <p><span class="label">注册时间:</span> {{ formatDate(user?.created_at) }}</p>
          </div>
        </div>
      </div>
      
      <el-divider />
      
      <div class="actions">
        <el-button type="danger" :icon="SwitchButton" @click="handleLogout">
          退出登录
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { SwitchButton } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

const formatDate = (dateString?: string) => {
  if (!dateString) return '未知'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
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
.home-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.user-card {
  width: 100%;
  max-width: 600px;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  background: #ffffff;
}

.user-info {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

.avatar {
  flex-shrink: 0;
  border: 4px solid #f0f2f5;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.info-content {
  flex: 1;
}

.nickname {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.username {
  margin: 0 0 16px 0;
  font-size: 16px;
  color: #909399;
}

.user-details {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.status-tag {
  font-weight: 500;
}

.user-meta {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.8;
}

.user-meta p {
  margin: 0;
  color: #606266;
}

.label {
  font-weight: 600;
  color: #303133;
  margin-right: 8px;
}

.actions {
  display: flex;
  justify-content: center;
  padding-top: 8px;
}

:deep(.el-divider) {
  margin: 24px 0;
}

:deep(.el-button) {
  border-radius: 8px;
  padding: 12px 32px;
  font-size: 16px;
  font-weight: 500;
}
</style>

