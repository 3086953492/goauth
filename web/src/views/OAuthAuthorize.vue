<template>
  <div class="oauth-authorize-page">
    <el-card class="oauth-authorize-page__card">
      <template #header>
        <div class="oauth-authorize-page__header">
          <h2>授权确认</h2>
          <p>第三方应用请求访问您的账号</p>
        </div>
      </template>

      <!-- 无效请求提示 -->
      <div v-if="!isValidRequest" class="oauth-authorize-page__error">
        <el-alert title="无效的授权请求" type="error" :closable="false" show-icon>
          <template #default>
            <p>授权请求参数不完整，请联系应用提供方。</p>
          </template>
        </el-alert>
      </div>

      <!-- 正常授权信息展示 -->
      <div v-else class="oauth-authorize-page__content">
        <!-- 用户信息 -->
        <div class="oauth-authorize-page__section">
          <div class="oauth-authorize-page__section-title">当前登录用户</div>
          <div class="oauth-authorize-page__user-info">
            <el-avatar :size="avatarSize" :src="currentUser?.avatar">
              {{ currentUser?.nickname?.[0] || currentUser?.username?.[0] }}
            </el-avatar>
            <div class="oauth-authorize-page__user-details">
              <div class="oauth-authorize-page__user-name">{{ currentUser?.nickname || currentUser?.username }}</div>
              <div class="oauth-authorize-page__user-username">@{{ currentUser?.username }}</div>
            </div>
          </div>
        </div>

        <!-- 客户端信息 -->
        <div class="oauth-authorize-page__section">
          <div class="oauth-authorize-page__section-title">授权给</div>
          <div class="oauth-authorize-page__client-info">
            <div class="oauth-authorize-page__client-item">
              <el-icon>
                <Key />
              </el-icon>
              <span>客户端 ID: {{ oauthParams.client_id }}</span>
            </div>
            <div class="oauth-authorize-page__client-item">
              <el-icon>
                <Link />
              </el-icon>
              <span>回调地址: {{ oauthParams.redirect_uri }}</span>
            </div>
          </div>
        </div>

        <!-- 权限范围 -->
        <div v-if="oauthParams.scope" class="oauth-authorize-page__section">
          <div class="oauth-authorize-page__section-title">请求的权限</div>
          <div class="oauth-authorize-page__scope-list">
            <el-tag v-for="scope in scopeList" :key="scope" type="info" size="large">
              {{ scope }}
            </el-tag>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="oauth-authorize-page__actions">
          <el-button type="primary" size="large" :loading="authorizing" @click="handleAuthorize">
            {{ authorizing ? '授权中...' : '确认授权' }}
          </el-button>
          <el-button size="large" @click="handleCancel">
            取消
          </el-button>
        </div>

        <!-- 提示信息 -->
        <div class="oauth-authorize-page__warning">
          <el-icon>
            <Warning />
          </el-icon>
          <span>授权后，该应用将能够访问您的账号信息</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Key, Link, Warning } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { usePermission } from '@/composables/usePermission'
import { confirmAuthorize } from '@/api/oauth'

const route = useRoute()
const authStore = useAuthStore()
const { checkLogin } = usePermission()

// 头像尺寸（对应 --icon-size-medium）
const avatarSize = 50

// OAuth 参数
const oauthParams = ref({
  client_id: '',
  redirect_uri: '',
  response_type: '',
  scope: '',
  state: ''
})

// 授权中状态
const authorizing = ref(false)

// 当前用户
const currentUser = computed(() => authStore.user)

// 请求是否有效
const isValidRequest = computed(() => {
  return !!(oauthParams.value.client_id && oauthParams.value.redirect_uri)
})

// 权限列表
const scopeList = computed(() => {
  if (!oauthParams.value.scope) return []
  return oauthParams.value.scope.split(' ').filter(s => s.trim())
})

// 初始化
onMounted(() => {
  // 读取 OAuth 参数（先读取参数，以便在重定向时保留）
  oauthParams.value = {
    client_id: (route.query.client_id as string) || '',
    redirect_uri: (route.query.redirect_uri as string) || '',
    response_type: (route.query.response_type as string) || 'code',
    scope: (route.query.scope as string) || '',
    state: (route.query.state as string) || ''
  }

  // 检查登录状态（如果未登录，会重定向到登录页并保留完整的查询参数）
  if (!checkLogin()) {
    return
  }

  // 参数校验
  if (!isValidRequest.value) {
    ElMessage.error('授权请求参数不完整')
  }
})

// 确认授权
const handleAuthorize = async () => {
  if (!isValidRequest.value) {
    ElMessage.error('授权请求参数不完整')
    return
  }

  // 检查登录状态（统一由 checkLogin 处理提示和跳转）
  if (!checkLogin()) {
    return
  }

  authorizing.value = true

  try {
    // 调用授权 API（使用 Cookie 鉴权）
    const response = await confirmAuthorize({
      client_id: oauthParams.value.client_id,
      redirect_uri: oauthParams.value.redirect_uri,
      response_type: oauthParams.value.response_type || 'code',
      scope: oauthParams.value.scope || undefined,
      state: oauthParams.value.state || undefined
    })

    // 授权成功，跳转到回调地址
    if (response.data.redirect_uri) {
      window.location.href = response.data.redirect_uri
    } else {
      ElMessage.error('授权失败：未返回回调地址')
    }
  } catch (error: any) {
    console.error('授权失败:', error)
    // 错误已在拦截器中处理
  } finally {
    authorizing.value = false
  }
}

// 取消授权
const handleCancel = () => {
  // 构造错误回调
  if (oauthParams.value.redirect_uri) {
    const redirectUrl = new URL(oauthParams.value.redirect_uri)
    redirectUrl.searchParams.append('error', 'access_denied')
    redirectUrl.searchParams.append('error_description', '用户拒绝授权')

    if (oauthParams.value.state) {
      redirectUrl.searchParams.append('state', oauthParams.value.state)
    }

    window.location.href = redirectUrl.toString()
  } else {
    ElMessage.info('已取消授权')
  }
}
</script>

<style scoped>
.oauth-authorize-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--color-page-background-alt);
  background-image:
    radial-gradient(circle at var(--pattern-size-dot-small) var(--pattern-size-dot-small), rgba(0, 0, 0, 0.03) 2%, transparent 0%),
    radial-gradient(circle at var(--pattern-size-dot-large) var(--pattern-size-dot-large), rgba(0, 0, 0, 0.03) 2%, transparent 0%);
  background-size: var(--pattern-size-grid) var(--pattern-size-grid);
  padding: var(--spacing-lg);
}

.oauth-authorize-page__card {
  width: 100%;
  max-width: var(--container-max-width-small);
  border-radius: var(--border-radius-xlarge);
  box-shadow: var(--shadow-auth-card);
  background: var(--color-card-background);
}

.oauth-authorize-page__header {
  text-align: center;
}

.oauth-authorize-page__header h2 {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: var(--font-size-display);
  font-weight: 600;
  color: var(--color-text-primary);
}

.oauth-authorize-page__header p {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-tertiary);
}

.oauth-authorize-page__error {
  padding: var(--spacing-lg) 0;
}

.oauth-authorize-page__content {
  padding: var(--spacing-sm-md) 0;
}

.oauth-authorize-page__section {
  margin-bottom: var(--spacing-lg-xl);
  padding-bottom: var(--spacing-lg);
  border-bottom: var(--border-width-thin) solid var(--color-border-lighter);
}

.oauth-authorize-page__section:last-of-type {
  border-bottom: none;
}

.oauth-authorize-page__section-title {
  font-size: var(--font-size-base);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--spacing-md);
}

.oauth-authorize-page__user-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.oauth-authorize-page__user-details {
  flex: 1;
}

.oauth-authorize-page__user-name {
  font-size: var(--font-size-lg);
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: var(--spacing-xs);
}

.oauth-authorize-page__user-username {
  font-size: var(--font-size-sm);
  color: var(--color-text-tertiary);
}

.oauth-authorize-page__client-info {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm-lg);
}

.oauth-authorize-page__client-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm-lg);
  background: var(--color-background-lighter);
  border-radius: var(--border-radius-large);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  word-break: break-all;
}

.oauth-authorize-page__client-item .el-icon {
  font-size: var(--icon-size-small);
  color: var(--color-info-icon);
  flex-shrink: 0;
}

.oauth-authorize-page__scope-list {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm-md);
}

.oauth-authorize-page__actions {
  display: flex;
  justify-content: center;
  gap: var(--spacing-md);
  margin-top: var(--spacing-lg-xl);
  margin-bottom: var(--spacing-lg);
}

.oauth-authorize-page__actions .el-button {
  min-width: var(--button-min-width);
  height: var(--button-height-large);
  font-size: var(--font-size-base);
  font-weight: 500;
  border-radius: var(--border-radius-large);
}

.oauth-authorize-page__warning {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  font-size: var(--font-size-xxs);
  color: var(--color-warning-text);
  padding: var(--spacing-sm-lg);
  background: var(--color-warning-background);
  border-radius: var(--border-radius-large);
}

.oauth-authorize-page__warning .el-icon {
  font-size: var(--icon-size-xsmall);
}

:deep(.el-avatar) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-size: var(--font-size-xl);
  font-weight: 600;
}

:deep(.el-tag) {
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-sm);
}

/* 响应式设计 */
/* 小屏幕：对应 --breakpoint-mobile (480px) */
@media (max-width: 480px) {
  .oauth-authorize-page {
    padding: var(--spacing-md);
  }

  .oauth-authorize-page__card {
    border-radius: var(--border-radius-large);
  }

  .oauth-authorize-page__header h2 {
    font-size: var(--font-size-title);
  }

  .oauth-authorize-page__user-info {
    flex-direction: column;
    text-align: center;
  }

  .oauth-authorize-page__actions {
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .oauth-authorize-page__actions .el-button {
    width: 100%;
  }
}
</style>
