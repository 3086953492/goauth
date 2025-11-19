<template>
  <div class="home-page">
    <Navbar />
    <div class="home-page__container">
      <el-card class="home-page__user-card">
      <div class="home-page__user-info">
        <el-avatar :size="80" :src="user?.avatar" :icon="Avatar" class="home-page__avatar" />
        
        <div class="home-page__info-content">
          <h2 class="home-page__nickname">{{ user?.nickname || '未知用户' }}</h2>
          <p class="home-page__username">@{{ user?.username }}</p>
          
          <div class="home-page__details">
            <el-tag type="info" size="large">{{ user?.role || 'user' }}</el-tag>
            <el-tag :type="user?.status === 1 ? 'success' : 'danger'" size="large" class="home-page__status-tag">
              {{ user?.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </div>
          
          <div class="home-page__meta">
            <p><span class="home-page__meta-label">用户ID:</span> {{ user?.id }}</p>
            <p><span class="home-page__meta-label">注册时间:</span> {{ formatDate(user?.created_at) }}</p>
          </div>
        </div>
      </div>
    </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/useAuthStore'
import Navbar from '@/components/Navbar.vue'
import { Avatar } from '@element-plus/icons-vue'

const authStore = useAuthStore()

const user = computed(() => authStore.user)

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
.home-page {
  min-height: 100vh;
  background: 
    linear-gradient(135deg, rgba(245, 247, 250, 0.8) 0%, rgba(228, 231, 235, 0.9) 100%),
    repeating-linear-gradient(45deg, transparent, transparent var(--pattern-size-small), rgba(0, 0, 0, 0.02) var(--pattern-size-small), rgba(0, 0, 0, 0.02) var(--pattern-size-large)),
    repeating-linear-gradient(-45deg, transparent, transparent var(--pattern-size-small), rgba(0, 0, 0, 0.01) var(--pattern-size-small), rgba(0, 0, 0, 0.01) var(--pattern-size-large)),
    var(--color-background-light);
}

.home-page__container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: var(--page-padding-top) var(--spacing-lg) var(--spacing-lg);
}

.home-page__user-card {
  width: 100%;
  max-width: var(--container-max-width-small);
  border-radius: var(--border-radius-card-large);
  box-shadow: var(--shadow-card-layered);
  background: var(--color-card-background);
  border: var(--border-width-thin) solid var(--color-border-white-translucent);
  overflow: hidden;
}

.home-page__user-card :deep(.el-card__body) {
  padding: var(--spacing-xl);
}

.home-page__user-info {
  display: flex;
  gap: var(--spacing-lg);
  align-items: flex-start;
}

.home-page__avatar {
  flex-shrink: 0;
  border: var(--border-width-thick) solid var(--color-background-lighter);
  box-shadow: var(--shadow-avatar);
}

.home-page__info-content {
  flex: 1;
}

.home-page__nickname {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: var(--font-size-display);
  font-weight: 600;
  color: var(--color-text-primary);
}

.home-page__username {
  margin: 0 0 var(--spacing-md) 0;
  font-size: var(--font-size-base);
  color: var(--color-text-tertiary);
}

.home-page__details {
  display: flex;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-lg);
}

.home-page__status-tag {
  font-weight: 500;
}

.home-page__meta {
  background: var(--color-background-light);
  padding: var(--spacing-md-lg) var(--spacing-lg);
  border-radius: var(--border-radius-xxlarge);
  font-size: var(--font-size-sm);
  line-height: 1.8;
  border: var(--border-width-thin) solid var(--color-border-lightest);
  box-shadow: var(--shadow-meta);
}

.home-page__meta p {
  margin: 0;
  color: var(--color-text-secondary);
}

.home-page__meta-label {
  font-weight: 600;
  color: var(--color-text-primary);
  margin-right: var(--spacing-sm);
}

/* 响应式设计 */
/* 平板端：对应 --breakpoint-tablet (768px) */
@media (max-width: 768px) {
  .home-page__container {
    padding: var(--page-padding-top-tablet) var(--spacing-md) var(--spacing-md);
  }

  .home-page__user-card {
    border-radius: var(--border-radius-xlarge);
  }

  .home-page__user-card :deep(.el-card__body) {
    padding: var(--spacing-lg);
  }

  .home-page__user-info {
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: var(--spacing-md);
  }

  .home-page__info-content {
    width: 100%;
  }

  .home-page__nickname {
    font-size: var(--font-size-title);
  }

  .home-page__username {
    font-size: var(--font-size-sm);
  }

  .home-page__details {
    justify-content: center;
  }

  .home-page__meta {
    padding: var(--spacing-md) var(--spacing-md);
    text-align: left;
  }
}

/* 移动端：对应 --breakpoint-mobile (480px) */
@media (max-width: 480px) {
  .home-page__container {
    padding: var(--page-padding-top-mobile) var(--spacing-sm) var(--spacing-sm);
  }

  .home-page__user-card {
    border-radius: var(--border-radius-large);
  }

  .home-page__user-card :deep(.el-card__body) {
    padding: var(--spacing-md);
  }

  .home-page__nickname {
    font-size: var(--font-size-xl);
  }

  .home-page__details {
    flex-direction: column;
    gap: var(--spacing-xs);
    align-items: center;
  }
}
</style>

