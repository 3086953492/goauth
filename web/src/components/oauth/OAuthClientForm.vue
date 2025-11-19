<template>
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px" :disabled="mode === 'view'">
        <el-form-item label="客户端名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入客户端名称（3-20字符）" maxlength="20" show-word-limit clearable />
        </el-form-item>

        <el-form-item v-if="mode === 'create'" label="客户端密钥" prop="client_secret">
            <div class="oauth-client-form__secret-wrapper">
                <el-input v-model="formData.client_secret" placeholder="客户端密钥" readonly>
                    <template #append>
                        <el-button-group>
                            <el-button :icon="CopyDocument" @click="copySecret">复制</el-button>
                            <el-button :icon="RefreshRight" @click="handleRegenerateSecret">
                                重新生成
                            </el-button>
                        </el-button-group>
                    </template>
                </el-input>
            </div>
            <div class="oauth-client-form__tip">请妥善保管客户端密钥，创建后将无法再次查看完整密钥</div>
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
    </el-form>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { CopyDocument, RefreshRight, Delete, Plus } from '@element-plus/icons-vue'
import { useOAuthClientForm } from '@/composables/useOAuthClientForm'
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
    initClientSecret,
    regenerateClientSecret,
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
    ]
}

// 复制密钥
const copySecret = async () => {
    try {
        await navigator.clipboard.writeText(formData.client_secret)
        ElMessage.success('密钥已复制到剪贴板')
    } catch (error) {
        ElMessage.error('复制失败，请手动复制')
    }
}

// 重新生成密钥
const handleRegenerateSecret = () => {
    regenerateClientSecret()
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
        data.client_secret = formData.client_secret
    }
    
    // 只包含已修改的字段
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
    }
}

// 初始化
onMounted(() => {
    if (props.mode === 'create') {
        initClientSecret()
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
