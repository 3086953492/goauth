<template>
  <div class="home-wrapper">
    <Navbar />
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
    </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'

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
</script>

<style scoped>
.home-wrapper {
  min-height: 100vh;
  background: 
    linear-gradient(135deg, rgba(245, 247, 250, 0.8) 0%, rgba(228, 231, 235, 0.9) 100%),
    repeating-linear-gradient(45deg, transparent, transparent 35px, rgba(0, 0, 0, 0.02) 35px, rgba(0, 0, 0, 0.02) 70px),
    repeating-linear-gradient(-45deg, transparent, transparent 35px, rgba(0, 0, 0, 0.01) 35px, rgba(0, 0, 0, 0.01) 70px),
    #f8f9fa;
}

.home-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 84px 20px 20px;
}

.user-card {
  width: 100%;
  max-width: 600px;
  border-radius: 24px;
  box-shadow: 
    0 2px 8px rgba(0, 0, 0, 0.04),
    0 8px 24px rgba(0, 0, 0, 0.06),
    0 16px 48px rgba(0, 0, 0, 0.08);
  background: #ffffff;
  border: 1px solid rgba(255, 255, 255, 0.8);
  overflow: hidden;
}

:deep(.el-card__body) {
  padding: 32px;
}

.user-info {
  display: flex;
  gap: 24px;
  align-items: flex-start;
}

.avatar {
  flex-shrink: 0;
  border: 4px solid #f5f7fa;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
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
  background: #f8f9fa;
  padding: 18px 20px;
  border-radius: 14px;
  font-size: 14px;
  line-height: 1.8;
  border: 1px solid #e9ecef;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
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
</style>

