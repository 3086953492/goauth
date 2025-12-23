<template>
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="140px" :disabled="mode === 'view'">
        <!-- 基本信息 -->
        <el-divider content-position="left">基本信息</el-divider>

        <el-form-item label="客户端名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入客户端名称（3-20字符）" maxlength="20" show-word-limit clearable />
        </el-form-item>

        <el-form-item label="应用描述" prop="description">
            <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入应用描述（选填，最多255字符）"
                maxlength="255" show-word-limit />
        </el-form-item>

        <el-form-item label="Logo URL" prop="logo">
            <el-input v-model="formData.logo" placeholder="请输入Logo图片URL（选填）" clearable />
        </el-form-item>

        <el-form-item label="回调地址" prop="redirect_uris" required>
            <div class="oauth-client-form__uris-wrapper">
                <div v-for="(_uri, index) in formData.redirect_uris" :key="index" class="oauth-client-form__uri-item">
                    <el-input v-model="formData.redirect_uris[index]" placeholder="https://example.com/callback"
                        clearable />
                    <el-button v-if="formData.redirect_uris.length > 1" :icon="Delete" type="danger" plain
                        @click="handleRemoveUri(index)">
                        删除
                    </el-button>
                </div>
                <el-button v-if="mode !== 'view'" :icon="Plus" type="primary" plain @click="handleAddUri">
                    添加回调地址
                </el-button>
            </div>
            <div class="oauth-client-form__tip">OAuth 授权成功后的回调地址，至少需要一个有效地址</div>
        </el-form-item>

        <el-form-item label="授权类型" prop="grant_types">
            <el-checkbox-group v-model="formData.grant_types">
                <el-checkbox v-for="grantType in OAUTH_GRANT_TYPES" :key="grantType.value" :label="grantType.value">
                    {{ grantType.label }}
                </el-checkbox>
            </el-checkbox-group>
        </el-form-item>

        <el-form-item label="权限范围" prop="scopes">
            <el-checkbox-group v-model="formData.scopes">
                <el-checkbox v-for="scope in OAUTH_SCOPES" :key="scope.value" :label="scope.value">
                    {{ scope.label }}
                </el-checkbox>
            </el-checkbox-group>
        </el-form-item>

        <el-form-item label="状态" prop="status">
            <el-radio-group v-model="formData.status">
                <el-radio v-for="status in OAUTH_CLIENT_STATUS" :key="status.value" :label="status.value">
                    {{ status.label }}
                </el-radio>
            </el-radio-group>
        </el-form-item>

        <!-- 密钥配置 -->
        <el-divider content-position="left">
            {{ mode === 'create' ? '密钥配置' : '密钥轮换（留空则不更新）' }}
        </el-divider>

        <el-form-item label="客户端密钥" prop="client_secret" :required="mode === 'create'">
            <div class="oauth-client-form__secret-wrapper">
                <el-input v-model="formData.client_secret" :placeholder="mode === 'create' ? '客户端密钥' : '留空则不更新'" :readonly="mode === 'create'">
                    <template #append>
                        <el-button-group>
                            <el-button v-if="mode === 'create'" :icon="CopyDocument" @click="copyToClipboard(formData.client_secret)">复制</el-button>
                            <el-button :icon="RefreshRight" @click="handleRegenerateSecret('client_secret')">
                                {{ mode === 'create' ? '重新生成' : '生成新密钥' }}
                            </el-button>
                        </el-button-group>
                    </template>
                </el-input>
            </div>
            <div class="oauth-client-form__tip">{{ mode === 'create' ? '请妥善保管客户端密钥，创建后将无法再次查看完整密钥' : '如需轮换密钥，请点击生成新密钥按钮' }}</div>
        </el-form-item>

        <el-form-item label="访问令牌密钥" prop="access_token_secret" :required="mode === 'create'">
            <div class="oauth-client-form__secret-wrapper">
                <el-input v-model="formData.access_token_secret" :placeholder="mode === 'create' ? '访问令牌密钥' : '留空则不更新'" :readonly="mode === 'create'">
                    <template #append>
                        <el-button-group>
                            <el-button v-if="mode === 'create'" :icon="CopyDocument" @click="copyToClipboard(formData.access_token_secret)">复制</el-button>
                            <el-button :icon="RefreshRight" @click="handleRegenerateSecret('access_token_secret')">
                                {{ mode === 'create' ? '重新生成' : '生成新密钥' }}
                            </el-button>
                        </el-button-group>
                    </template>
                </el-input>
            </div>
            <div class="oauth-client-form__tip">用于签名访问令牌</div>
        </el-form-item>

        <el-form-item label="刷新令牌密钥" prop="refresh_token_secret" :required="mode === 'create'">
            <div class="oauth-client-form__secret-wrapper">
                <el-input v-model="formData.refresh_token_secret" :placeholder="mode === 'create' ? '刷新令牌密钥' : '留空则不更新'" :readonly="mode === 'create'">
                    <template #append>
                        <el-button-group>
                            <el-button v-if="mode === 'create'" :icon="CopyDocument" @click="copyToClipboard(formData.refresh_token_secret)">复制</el-button>
                            <el-button :icon="RefreshRight" @click="handleRegenerateSecret('refresh_token_secret')">
                                {{ mode === 'create' ? '重新生成' : '生成新密钥' }}
                            </el-button>
                        </el-button-group>
                    </template>
                </el-input>
            </div>
            <div class="oauth-client-form__tip">用于签名刷新令牌</div>
        </el-form-item>

        <el-form-item label="用户标识密钥" prop="subject_secret" :required="mode === 'create'">
            <div class="oauth-client-form__secret-wrapper">
                <el-input v-model="formData.subject_secret" :placeholder="mode === 'create' ? '用户标识密钥' : '留空则不更新'" :readonly="mode === 'create'">
                    <template #append>
                        <el-button-group>
                            <el-button v-if="mode === 'create'" :icon="CopyDocument" @click="copyToClipboard(formData.subject_secret)">复制</el-button>
                            <el-button :icon="RefreshRight" @click="handleRegenerateSecret('subject_secret')">
                                {{ mode === 'create' ? '重新生成' : '生成新密钥' }}
                            </el-button>
                        </el-button-group>
                    </template>
                </el-input>
            </div>
            <div class="oauth-client-form__tip">用于生成用户唯一标识（Subject）</div>
        </el-form-item>

        <!-- 过期时间配置 -->
        <el-divider content-position="left">过期时间配置</el-divider>

        <el-form-item label="授权码过期时间" prop="auth_code_expire">
            <el-input-number v-model="formData.auth_code_expire" :min="60" :max="600" :step="30" />
            <span class="oauth-client-form__unit">秒（60-600，默认300秒=5分钟）</span>
        </el-form-item>

        <el-form-item label="访问令牌过期时间" prop="access_token_expire">
            <el-input-number v-model="formData.access_token_expire" :min="300" :max="86400" :step="300" />
            <span class="oauth-client-form__unit">秒（300-86400，默认3600秒=1小时）</span>
        </el-form-item>

        <el-form-item label="刷新令牌过期时间" prop="refresh_token_expire">
            <el-input-number v-model="formData.refresh_token_expire" :min="3600" :max="31536000" :step="86400" />
            <span class="oauth-client-form__unit">秒（3600-31536000，默认2592000秒=30天）</span>
        </el-form-item>

        <!-- 用户标识配置 -->
        <el-divider content-position="left">用户标识配置</el-divider>

        <el-form-item label="用户标识长度" prop="subject_length">
            <el-input-number v-model="formData.subject_length" :min="8" :max="64" :step="1" />
            <span class="oauth-client-form__unit">字符（8-64，默认16）</span>
        </el-form-item>

        <el-form-item label="用户标识前缀" prop="subject_prefix">
            <el-input v-model="formData.subject_prefix" placeholder="usr_" maxlength="20" style="width: 200px;" />
            <span class="oauth-client-form__unit">（最多20字符，默认 usr_）</span>
        </el-form-item>
    </el-form>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { CopyDocument, RefreshRight, Delete, Plus } from '@element-plus/icons-vue'
import { useOAuthClientForm, DEFAULT_AUTH_CODE_EXPIRE, DEFAULT_ACCESS_TOKEN_EXPIRE, DEFAULT_REFRESH_TOKEN_EXPIRE, DEFAULT_SUBJECT_LENGTH, DEFAULT_SUBJECT_PREFIX } from '@/composables/useOAuthClientForm'
import { OAUTH_GRANT_TYPES, OAUTH_SCOPES, OAUTH_CLIENT_STATUS } from '@/constants'
import type { OAuthClientFormMode, OAuthClientDetailResponse } from '@/types/oauth_client'

interface Props {
    mode?: OAuthClientFormMode
    initialData?: OAuthClientDetailResponse
}

const props = withDefaults(defineProps<Props>(), {
    mode: 'create',
    initialData: undefined
})

const emit = defineEmits<{
    submit: []
    cancel: []
}>()

const formRef = ref<FormInstance>()

const {
    formData,
    initSecrets,
    regenerateSecret,
    addRedirectUri,
    removeRedirectUri
} = useOAuthClientForm()

// 表单验证规则
const formRules: FormRules = {
    name: [
        { required: true, message: '请输入客户端名称', trigger: 'blur' },
        { min: 3, max: 20, message: '客户端名称长度应为3-20字符', trigger: 'blur' }
    ],
    client_secret: [
        { required: props.mode === 'create', message: '客户端密钥不能为空', trigger: 'blur' }
    ],
    access_token_secret: [
        { required: props.mode === 'create', message: '访问令牌密钥不能为空', trigger: 'blur' }
    ],
    refresh_token_secret: [
        { required: props.mode === 'create', message: '刷新令牌密钥不能为空', trigger: 'blur' }
    ],
    subject_secret: [
        { required: props.mode === 'create', message: '用户标识密钥不能为空', trigger: 'blur' }
    ],
    description: [
        { max: 255, message: '应用描述不能超过255字符', trigger: 'blur' }
    ],
    logo: [
        {
            pattern: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/,
            message: '请输入有效的URL',
            trigger: 'blur'
        }
    ],
    redirect_uris: [
        {
            validator: (_rule, value, callback) => {
                if (!value || value.length === 0) {
                    callback()
                    return
                }
                const validUris = value.filter((uri: string) => uri.trim() !== '')
                if (validUris.length === 0) {
                    callback()
                } else {
                    const urlPattern = /^https?:\/\/.+/
                    const invalidUris = validUris.filter((uri: string) => !urlPattern.test(uri))
                    if (invalidUris.length > 0) {
                        callback(new Error('请输入有效的URL（必须以http://或https://开头）'))
                    } else {
                        callback()
                    }
                }
            },
            trigger: 'change'
        }
    ],
    grant_types: [
        { type: 'array', required: props.mode === 'create', message: '请至少选择一种授权类型', trigger: 'change' }
    ],
    scopes: [
        { type: 'array', required: props.mode === 'create', message: '请至少选择一个权限范围', trigger: 'change' }
    ],
    status: [
        { required: props.mode === 'create', message: '请选择状态', trigger: 'change' }
    ],
    auth_code_expire: [
        { type: 'number', min: 60, max: 600, message: '授权码过期时间应为60-600秒', trigger: 'change' }
    ],
    access_token_expire: [
        { type: 'number', min: 300, max: 86400, message: '访问令牌过期时间应为300-86400秒', trigger: 'change' }
    ],
    refresh_token_expire: [
        { type: 'number', min: 3600, max: 31536000, message: '刷新令牌过期时间应为3600-31536000秒', trigger: 'change' }
    ],
    subject_length: [
        { type: 'number', min: 8, max: 64, message: '用户标识长度应为8-64', trigger: 'change' }
    ],
    subject_prefix: [
        { max: 20, message: '用户标识前缀不能超过20字符', trigger: 'blur' }
    ]
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
    try {
        await navigator.clipboard.writeText(text)
        ElMessage.success('已复制到剪贴板')
    } catch (error) {
        ElMessage.error('复制失败，请手动复制')
    }
}

// 重新生成密钥
const handleRegenerateSecret = (field: 'client_secret' | 'access_token_secret' | 'refresh_token_secret' | 'subject_secret') => {
    regenerateSecret(field)
}

// 添加回调地址
const handleAddUri = () => {
    addRedirectUri()
}

// 删除回调地址
const handleRemoveUri = (index: number) => {
    removeRedirectUri(index)
}

// 验证表单
const validate = async (): Promise<boolean> => {
    if (!formRef.value) return false
    try {
        await formRef.value.validate()
        return true
    } catch (error) {
        return false
    }
}

// 获取表单数据
const getFormData = () => {
    const data: any = {
        name: formData.name
    }
    
    if (props.mode === 'create') {
        // 创建模式：必填密钥字段
        data.client_secret = formData.client_secret
        data.access_token_secret = formData.access_token_secret
        data.refresh_token_secret = formData.refresh_token_secret
        data.subject_secret = formData.subject_secret
    } else {
        // 编辑模式：仅在填写时才带上密钥字段（轮换）
        if (formData.client_secret && formData.client_secret.trim() !== '') {
            data.client_secret = formData.client_secret
        }
        if (formData.access_token_secret && formData.access_token_secret.trim() !== '') {
            data.access_token_secret = formData.access_token_secret
        }
        if (formData.refresh_token_secret && formData.refresh_token_secret.trim() !== '') {
            data.refresh_token_secret = formData.refresh_token_secret
        }
        if (formData.subject_secret && formData.subject_secret.trim() !== '') {
            data.subject_secret = formData.subject_secret
        }
    }
    
    // 基本字段
    if (formData.description !== undefined && formData.description !== '') {
        data.description = formData.description
    }
    if (formData.logo !== undefined && formData.logo !== '') {
        data.logo = formData.logo
    }
    if (formData.redirect_uris && formData.redirect_uris.length > 0) {
        const validUris = formData.redirect_uris.filter(uri => uri.trim() !== '')
        if (validUris.length > 0) {
            data.redirect_uris = validUris
        }
    }
    if (formData.grant_types && formData.grant_types.length > 0) {
        data.grant_types = formData.grant_types
    }
    if (formData.scopes && formData.scopes.length > 0) {
        data.scopes = formData.scopes
    }
    if (formData.status !== undefined) {
        data.status = formData.status
    }

    // 配置字段
    if (formData.auth_code_expire !== undefined) {
        data.auth_code_expire = formData.auth_code_expire
    }
    if (formData.access_token_expire !== undefined) {
        data.access_token_expire = formData.access_token_expire
    }
    if (formData.refresh_token_expire !== undefined) {
        data.refresh_token_expire = formData.refresh_token_expire
    }
    if (formData.subject_length !== undefined) {
        data.subject_length = formData.subject_length
    }
    if (formData.subject_prefix !== undefined) {
        data.subject_prefix = formData.subject_prefix
    }
    
    return data
}

// 重置表单
const resetFields = () => {
    formRef.value?.resetFields()
}

// 加载初始数据
const loadInitialData = () => {
    if (props.initialData) {
        formData.name = props.initialData.name || ''
        formData.description = props.initialData.description || ''
        formData.logo = props.initialData.logo || ''
        formData.redirect_uris = props.initialData.redirect_uris && props.initialData.redirect_uris.length > 0 
            ? [...props.initialData.redirect_uris] 
            : ['']
        formData.grant_types = props.initialData.grant_types ? [...props.initialData.grant_types] : []
        formData.scopes = props.initialData.scopes ? [...props.initialData.scopes] : []
        formData.status = props.initialData.status ?? 1

        // 配置字段
        formData.auth_code_expire = props.initialData.auth_code_expire ?? DEFAULT_AUTH_CODE_EXPIRE
        formData.access_token_expire = props.initialData.access_token_expire ?? DEFAULT_ACCESS_TOKEN_EXPIRE
        formData.refresh_token_expire = props.initialData.refresh_token_expire ?? DEFAULT_REFRESH_TOKEN_EXPIRE
        formData.subject_length = props.initialData.subject_length ?? DEFAULT_SUBJECT_LENGTH
        formData.subject_prefix = props.initialData.subject_prefix ?? DEFAULT_SUBJECT_PREFIX

        // 编辑模式下清空密钥字段（轮换时填写）
        formData.client_secret = ''
        formData.access_token_secret = ''
        formData.refresh_token_secret = ''
        formData.subject_secret = ''
    }
}

// 初始化
onMounted(() => {
    if (props.mode === 'create') {
        initSecrets()
    } else if (props.mode === 'edit' && props.initialData) {
        loadInitialData()
    }
})

// 监听 initialData 变化
watch(() => props.initialData, (newData) => {
    if (newData && props.mode === 'edit') {
        loadInitialData()
    }
}, { deep: true })

// 暴露方法给父组件
defineExpose({
    validate,
    getFormData,
    resetFields
})
</script>

<style scoped>
.oauth-client-form__secret-wrapper {
    width: 100%;
}

.oauth-client-form__secret-wrapper :deep(.el-input-group__append) {
    padding: 0;
}

.oauth-client-form__secret-wrapper :deep(.el-input-group__append .el-button-group) {
    display: flex;
}

.oauth-client-form__secret-wrapper :deep(.el-input-group__append .el-button) {
    margin: 0;
    border: none;
    border-radius: 0;
}

.oauth-client-form__tip {
    font-size: var(--font-size-xs);
    color: var(--color-text-tertiary);
    margin-top: var(--spacing-xs);
    line-height: 1.5;
}

.oauth-client-form__unit {
    margin-left: var(--spacing-sm);
    font-size: var(--font-size-sm);
    color: var(--color-text-secondary);
}

.oauth-client-form__uris-wrapper {
    width: 100%;
}

.oauth-client-form__uri-item {
    display: flex;
    gap: var(--spacing-sm-lg);
    margin-bottom: var(--spacing-sm-lg);
    align-items: center;
}

.oauth-client-form__uri-item .el-input {
    flex: 1;
}
</style>
