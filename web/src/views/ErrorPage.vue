<template>
  <div class="error-page">
    <el-card class="error-page__card">
      <div class="error-page__content">
        <!-- 错误图标 -->
        <div class="error-page__icon">
          <el-icon :size="iconSize">
            <CircleClose />
          </el-icon>
        </div>

        <!-- 错误标题 -->
        <h1 class="error-page__title">{{ displayError }}</h1>

        <!-- 错误描述 -->
        <p class="error-page__description">{{ displayDescription }}</p>

        <!-- 操作按钮 -->
        <div class="error-page__actions">
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

// 图标尺寸（对应 --icon-size-large）
const iconSize = 80

// 错误码到友好提示的映射
const errorMessages: Record<string, string> = {
  // OAuth 相关错误
  'invalid_request': '无效的请求',
  'unauthorized_client': '未授权的客户端',
  'access_denied': '访问被拒绝',
  'unsupported_response_type': '不支持的响应类型',
  'invalid_scope': '无效的权限范围',
  'server_error': '服务器错误',
  'temporarily_unavailable': '服务暂时不可用',
  'invalid_client': '无效的客户端',
  'invalid_grant': '无效的授权',
  'unsupported_grant_type': '不支持的授权类型',
  
  // HTTP 状态码相关错误
  '403': '访问被拒绝',
  '404': '页面不存在',
  '500': '服务器内部错误',
  '502': '网关错误',
  '503': '服务不可用',
  
  // 通用错误类型
  'forbidden': '您没有权限访问此资源',
  'not_found': '请求的资源不存在',
  'internal_error': '服务器内部错误，请稍后重试',
  'network_error': '网络连接失败',
  'timeout': '请求超时'
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
  if (errorDescription.value) {
    return errorDescription.value
  }
  
  // 根据错误类型提供默认描述
  if (error.value === '403' || error.value === 'forbidden') {
    return '您没有访问此页面的权限，请联系管理员或返回上一页。'
  }
  if (error.value === '404' || error.value === 'not_found') {
    return '抱歉，您访问的页面不存在，可能已被删除或移动。'
  }
  if (error.value === '500' || error.value === 'internal_error') {
    return '服务器遇到了一些问题，我们正在努力修复，请稍后再试。'
  }
  if (error.value === 'network_error') {
    return '网络连接失败，请检查您的网络设置后重试。'
  }
  if (error.value === 'timeout') {
    return '请求超时，服务器响应时间过长，请稍后重试。'
  }
  
  return '请求发生未知错误，请稍后重试。'
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
.error-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--color-page-background);
  background-image:
    radial-gradient(circle at 0% 0%, rgba(255, 255, 255, 0.8) 0, transparent 55%),
    radial-gradient(circle at 100% 100%, rgba(15, 23, 42, 0.08) 0, transparent 55%);
  padding: var(--spacing-lg);
}

.error-page__card {
  width: 100%;
  max-width: var(--container-max-width-small);
  border-radius: var(--border-radius-xlarge);
  box-shadow: var(--shadow-xlarge);
  background: var(--color-card-background);
  overflow: hidden;
}

.error-page__content {
  padding: var(--spacing-xxl) var(--spacing-lg);
  text-align: center;
}

.error-page__icon {
  margin-bottom: var(--spacing-lg);
  animation: shake 0.5s ease-in-out;
  color: var(--color-danger-icon);
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(calc(-1 * var(--animation-shake-distance))); }
  75% { transform: translateX(var(--animation-shake-distance)); }
}

.error-page__title {
  margin: 0 0 var(--spacing-md) 0;
  font-size: var(--font-size-heading);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.4;
  word-break: break-word;
}

.error-page__description {
  margin: 0 0 var(--spacing-xl) 0;
  font-size: var(--font-size-base);
  color: var(--color-text-secondary);
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  max-width: var(--container-max-width-xsmall);
  margin-left: auto;
  margin-right: auto;
}

.error-page__actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: center;
  flex-wrap: wrap;
}

.error-page__actions .el-button {
  min-width: var(--button-min-width);
  height: var(--button-height-large);
  font-size: var(--font-size-base);
  font-weight: 500;
  border-radius: var(--border-radius-large);
}

/* 响应式设计 */
/* 小屏幕：对应 --container-max-width-small (600px) */
@media (max-width: 600px) {
  .error-page {
    padding: var(--spacing-md);
  }

  .error-page__content {
    padding: var(--spacing-xl) var(--spacing-md);
  }

  .error-page__title {
    font-size: var(--font-size-title);
  }

  .error-page__description {
    font-size: var(--font-size-sm);
  }

  .error-page__actions {
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .error-page__actions .el-button {
    width: 100%;
  }
}

/* 使用中性浅灰 + 轻微蓝灰点缀的现代背景 */
.error-page::before {
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

