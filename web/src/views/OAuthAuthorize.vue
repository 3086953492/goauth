<template>
  <div class="authorize-container">
    <el-card class="authorize-card">
      <template #header>
        <div class="card-header">
          <h2>授权确认</h2>
          <p>第三方应用请求访问您的账号</p>
        </div>
      </template>

      <!-- 无效请求提示 -->
      <div v-if="!isValidRequest" class="error-message">
        <el-alert title="无效的授权请求" type="error" :closable="false" show-icon>
          <template #default>
            <p>授权请求参数不完整，请联系应用提供方。</p>
          </template>
        </el-alert>
      </div>

      <!-- 正常授权信息展示 -->
      <div v-else class="authorize-content">
        <!-- 用户信息 -->
        <div class="info-section">
          <div class="section-title">当前登录用户</div>
          <div class="user-info">
            <el-avatar :size="50" :src="currentUser?.avatar">
              {{ currentUser?.nickname?.[0] || currentUser?.username?.[0] }}
            </el-avatar>
            <div class="user-details">
              <div class="user-name">{{ currentUser?.nickname || currentUser?.username }}</div>
              <div class="user-username">@{{ currentUser?.username }}</div>
            </div>
          </div>
        </div>

        <!-- 客户端信息 -->
        <div class="info-section">
          <div class="section-title">授权给</div>
          <div class="client-info">
            <div class="client-id">
              <el-icon>
                <Key />
              </el-icon>
              <span>客户端 ID: {{ oauthParams.client_id }}</span>
            </div>
            <div class="redirect-uri">
              <el-icon>
                <Link />
              </el-icon>
              <span>回调地址: {{ oauthParams.redirect_uri }}</span>
            </div>
          </div>
        </div>

        <!-- 权限范围 -->
        <div v-if="oauthParams.scope" class="info-section">
          <div class="section-title">请求的权限</div>
          <div class="scope-list">
            <el-tag v-for="scope in scopeList" :key="scope" type="info" size="large">
              {{ scope }}
            </el-tag>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons">
          <el-button type="primary" size="large" :loading="authorizing" @click="handleAuthorize">
            {{ authorizing ? '授权中...' : '确认授权' }}
          </el-button>
          <el-button size="large" @click="handleCancel">
            取消
          </el-button>
        </div>

        <!-- 提示信息 -->
        <div class="warning-text">
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
import { useAuthStore } from '@/stores/auth'
import { usePermission } from '@/composables/usePermission'
import { API_BASE_URL } from '@/constants'

const route = useRoute()
const authStore = useAuthStore()
const { checkLogin } = usePermission()

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
const handleAuthorize = () => {
  if (!isValidRequest.value) {
    ElMessage.error('授权请求参数不完整')
    return
  }

  // 从 localStorage 获取 token
  const token = localStorage.getItem('token')
  if (!token) {
    ElMessage.error('未登录或登录已过期')
    return
  }

  // 设置 Cookie (后端中间件会从 Cookie 中读取 token)
  document.cookie = `access_token=${token}; path=/; SameSite=Lax`

  // 构建后端授权 URL
  const authUrl = new URL(`${API_BASE_URL}/api/v1/oauth/authorization`, window.location.origin)
  authUrl.searchParams.append('client_id', oauthParams.value.client_id)
  authUrl.searchParams.append('redirect_uri', oauthParams.value.redirect_uri)
  authUrl.searchParams.append('response_type', oauthParams.value.response_type || 'code')
  
  if (oauthParams.value.scope) {
    authUrl.searchParams.append('scope', oauthParams.value.scope)
  }
  
  if (oauthParams.value.state) {
    authUrl.searchParams.append('state', oauthParams.value.state)
  }

  // 直接跳转到后端授权接口，后端会处理授权并重定向到第三方应用
  window.location.href = authUrl.toString()
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
.authorize-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: #f0f2f5;
  background-image:
    radial-gradient(circle at 25px 25px, rgba(0, 0, 0, 0.03) 2%, transparent 0%),
    radial-gradient(circle at 75px 75px, rgba(0, 0, 0, 0.03) 2%, transparent 0%);
  background-size: 100px 100px;
  padding: 20px;
}

.authorize-card {
  width: 100%;
  max-width: 600px;
  border-radius: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08);
  background: #ffffff;
}

.card-header {
  text-align: center;
}

.card-header h2 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.card-header p {
  margin: 0;
  font-size: 14px;
  color: #909399;
}

.error-message {
  padding: 20px 0;
}

.authorize-content {
  padding: 10px 0;
}

.info-section {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.info-section:last-of-type {
  border-bottom: none;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: 18px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.user-username {
  font-size: 14px;
  color: #909399;
}

.client-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.client-id,
.redirect-uri {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  font-size: 14px;
  color: #606266;
  word-break: break-all;
}

.client-id .el-icon,
.redirect-uri .el-icon {
  font-size: 18px;
  color: #409eff;
  flex-shrink: 0;
}

.scope-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 30px;
  margin-bottom: 20px;
}

.action-buttons .el-button {
  min-width: 140px;
  height: 44px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
}

.warning-text {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 13px;
  color: #e6a23c;
  padding: 12px;
  background: #fdf6ec;
  border-radius: 8px;
}

.warning-text .el-icon {
  font-size: 16px;
}

:deep(.el-avatar) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-size: 20px;
  font-weight: 600;
}

:deep(.el-tag) {
  padding: 8px 16px;
  font-size: 14px;
}
</style>
