<template>
  <div class="error-container">
    <el-card class="error-card">
      <div class="error-content">
        <!-- 错误图标 -->
        <div class="error-icon">
          <el-icon :size="80" color="#f56c6c">
            <CircleClose />
          </el-icon>
        </div>

        <!-- 错误标题 -->
        <h1 class="error-title">{{ displayError }}</h1>

        <!-- 错误描述 -->
        <p class="error-description">{{ displayDescription }}</p>

        <!-- 操作按钮 -->
        <div class="error-actions">
          <el-button type="default" size="large" @click="goBack">
            返回上一页
          </el-button>
          <el-button type="primary" size="large" @click="goHome">
            回到首页
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CircleClose } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()

// 错误码到友好提示的映射
const errorMessages: Record<string, string> = {
  'invalid_request': '无效的请求',
  'unauthorized_client': '未授权的客户端',
  'access_denied': '访问被拒绝',
  'unsupported_response_type': '不支持的响应类型',
  'invalid_scope': '无效的权限范围',
  'server_error': '服务器错误',
  'temporarily_unavailable': '服务暂时不可用',
  'invalid_client': '无效的客户端',
  'invalid_grant': '无效的授权',
  'unsupported_grant_type': '不支持的授权类型'
}

// 从 query 中读取错误信息
const error = computed(() => {
  const errorParam = route.query.error
  return typeof errorParam === 'string' ? errorParam : ''
})

const errorDescription = computed(() => {
  const descParam = route.query.error_description
  return typeof descParam === 'string' ? descParam : ''
})

// 显示的错误标题（优先使用映射的友好提示）
const displayError = computed(() => {
  if (error.value && errorMessages[error.value]) {
    return errorMessages[error.value]
  }
  return error.value || '发生错误'
})

// 显示的错误描述
const displayDescription = computed(() => {
  return errorDescription.value || '请求发生未知错误，请稍后重试。'
})

// 返回上一页
const goBack = () => {
  router.back()
}

// 回到首页
const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
.error-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #edf2f7;
  background-image:
    radial-gradient(circle at 0% 0%, rgba(255, 255, 255, 0.8) 0, transparent 55%),
    radial-gradient(circle at 100% 100%, rgba(15, 23, 42, 0.08) 0, transparent 55%);
  padding: 20px;
}

.error-card {
  width: 100%;
  max-width: 600px;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  background: #ffffff;
  overflow: hidden;
}

.error-content {
  padding: 40px 20px;
  text-align: center;
}

.error-icon {
  margin-bottom: 24px;
  animation: shake 0.5s ease-in-out;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-10px); }
  75% { transform: translateX(10px); }
}

.error-title {
  margin: 0 0 16px 0;
  font-size: 32px;
  font-weight: 600;
  color: #303133;
  line-height: 1.4;
  word-break: break-word;
}

.error-description {
  margin: 0 0 32px 0;
  font-size: 16px;
  color: #606266;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
}

.error-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}

.error-actions .el-button {
  min-width: 140px;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
}

/* 响应式设计 */
@media (max-width: 600px) {
  .error-container {
    padding: 16px;
  }

  .error-content {
    padding: 32px 16px;
  }

  .error-title {
    font-size: 24px;
  }

  .error-description {
    font-size: 14px;
  }

  .error-actions {
    flex-direction: column;
    gap: 12px;
  }

  .error-actions .el-button {
    width: 100%;
  }
}

/* 使用中性浅灰 + 轻微蓝灰点缀的现代背景 */
.error-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 40%, rgba(255, 255, 255, 0.3) 0%, transparent 55%),
    radial-gradient(circle at 80% 70%, rgba(148, 163, 184, 0.25) 0%, transparent 55%);
  pointer-events: none;
}
</style>

