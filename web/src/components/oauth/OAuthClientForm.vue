<template>
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px" :disabled="mode === 'view'">
        <el-form-item label="客户端名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入客户端名称（3-20字符）" maxlength="20" show-word-limit clearable />
        </el-form-item>

        <el-form-item label="客户端密钥" prop="client_secret">
            <div class="secret-input-wrapper">
                <el-input v-model="formData.client_secret" placeholder="客户端密钥" readonly>
                    <template #append>
                        <el-button-group>
                            <el-button :icon="CopyDocument" @click="copySecret">复制</el-button>
                            <el-button v-if="mode === 'create'" :icon="RefreshRight" @click="handleRegenerateSecret">
                                重新生成
                            </el-button>
                        </el-button-group>
                    </template>
                </el-input>
            </div>
            <div class="form-tip">请妥善保管客户端密钥，创建后将无法再次查看完整密钥</div>
        </el-form-item>

        <el-form-item label="应用描述" prop="description">
            <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入应用描述（选填，最多255字符）"
                maxlength="255" show-word-limit />
        </el-form-item>

        <el-form-item label="Logo URL" prop="logo">
            <el-input v-model="formData.logo" placeholder="请输入Logo图片URL（选填）" clearable />
        </el-form-item>

        <el-form-item label="回调地址" prop="redirect_uris" required>
            <div class="redirect-uris-wrapper">
                <div v-for="(_uri, index) in formData.redirect_uris" :key="index" class="redirect-uri-item">
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
            <div class="form-tip">OAuth 授权成功后的回调地址，至少需要一个有效地址</div>
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
import { ref, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { CopyDocument, RefreshRight, Delete, Plus } from '@element-plus/icons-vue'
import { useOAuthClientForm } from '@/composables/useOAuthClientForm'
import { OAUTH_GRANT_TYPES, OAUTH_SCOPES, OAUTH_CLIENT_STATUS } from '@/constants'
import type { OAuthClientFormMode } from '@/types/oauth_client'

interface Props {
    mode?: OAuthClientFormMode
    initialData?: any
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
        { required: true, message: '客户端密钥不能为空', trigger: 'blur' }
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
                const validUris = value.filter((uri: string) => uri.trim() !== '')
                if (validUris.length === 0) {
                    callback(new Error('至少需要一个有效的回调地址'))
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
        { type: 'array', required: true, message: '请至少选择一种授权类型', trigger: 'change' }
    ],
    scopes: [
        { type: 'array', required: true, message: '请至少选择一个权限范围', trigger: 'change' }
    ],
    status: [
        { required: true, message: '请选择状态', trigger: 'change' }
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
    return {
        ...formData,
        redirect_uris: formData.redirect_uris.filter(uri => uri.trim() !== '')
    }
}

// 重置表单
const resetFields = () => {
    formRef.value?.resetFields()
}

// 初始化
onMounted(() => {
    if (props.mode === 'create') {
        initClientSecret()
    }
})

// 暴露方法给父组件
defineExpose({
    validate,
    getFormData,
    resetFields
})
</script>

<style scoped>
.secret-input-wrapper {
    width: 100%;
}

:deep(.el-input-group__append) {
    padding: 0;
}

:deep(.el-input-group__append .el-button-group) {
    display: flex;
}

:deep(.el-input-group__append .el-button) {
    margin: 0;
    border: none;
    border-radius: 0;
}

.form-tip {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
    line-height: 1.5;
}

.redirect-uris-wrapper {
    width: 100%;
}

.redirect-uri-item {
    display: flex;
    gap: 12px;
    margin-bottom: 12px;
    align-items: center;
}

.redirect-uri-item .el-input {
    flex: 1;
}
</style>
