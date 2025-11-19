<template>
    <div class="oauth-clients-page">
        <Navbar />
        <div class="oauth-clients-page__container">
            <el-card class="oauth-clients-page__card">
                <template #header>
                    <div class="oauth-clients-page__header">
                        <h2 class="oauth-clients-page__title">OAuth 客户端管理</h2>
                        <div class="oauth-clients-page__actions">
                            <div class="oauth-clients-page__filters">
                                <el-input v-model="filters.name" placeholder="搜索客户端名称" clearable
                                    class="oauth-clients-page__filter-input" @input="handleFilterChange" />
                                <el-select v-model="filters.status" placeholder="状态筛选" clearable
                                    class="oauth-clients-page__filter-select" @change="handleFilterChange">
                                    <el-option label="正常" :value="1" />
                                    <el-option label="禁用" :value="0" />
                                </el-select>
                            </div>
                            <el-button type="primary" :icon="Plus" @click="handleCreateClient">
                                新建客户端
                            </el-button>
                        </div>
                    </div>
                </template>

                <el-table v-loading="loading" :data="clientList" stripe style="width: 100%">
                    <el-table-column prop="id" label="ID" width="80" />
                    <el-table-column label="Logo" width="100">
                        <template #default="{ row }">
                            <el-avatar :size="avatarSize" :src="row.logo" :icon="Cpu" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="name" label="客户端名称" min-width="200" />
                    <el-table-column label="状态" width="120">
                        <template #default="{ row }">
                            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="large">
                                {{ row.status === 1 ? '正常' : '禁用' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="180">
                        <template #default="{ row }">
                            <el-button type="primary" link @click="handleEditClient(row.id)">
                                编辑
                            </el-button>
                            <el-button type="danger" link @click="handleDeleteClient(row.id)">
                                删除
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>

                <div class="oauth-clients-page__pagination">
                    <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
                        :page-sizes="[10, 20, 50, 100]" :total="pagination.total"
                        layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange"
                        @current-change="handlePageChange" />
                </div>
            </el-card>
        </div>

        <!-- 新建客户端弹窗 -->
        <el-dialog v-model="createDialogVisible" title="新建 OAuth 客户端" :width="dialogWidth" :close-on-click-modal="false"
            @close="handleCreateDialogClose">
            <OAuthClientForm ref="createFormRef" mode="create" />
            <template #footer>
                <div class="oauth-clients-page__dialog-footer">
                    <el-button @click="handleCancelCreate">取消</el-button>
                    <el-button type="primary" :loading="submitLoading" @click="handleSubmitCreate">
                        确认创建
                    </el-button>
                </div>
            </template>
        </el-dialog>

        <!-- 编辑客户端弹窗 -->
        <el-dialog v-model="editDialogVisible" title="编辑 OAuth 客户端" :width="dialogWidth" :close-on-click-modal="false"
            @close="handleEditDialogClose">
            <OAuthClientForm ref="editFormRef" mode="edit" :initial-data="currentClient" />
            <template #footer>
                <div class="oauth-clients-page__dialog-footer">
                    <el-button @click="handleCancelEdit">取消</el-button>
                    <el-button type="primary" :loading="submitLoading" @click="handleSubmitEdit">
                        确认更新
                    </el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Cpu } from '@element-plus/icons-vue'
import Navbar from '@/components/Navbar.vue'
import OAuthClientForm from '@/components/oauth/OAuthClientForm.vue'
import { useOAuthClientList } from '@/composables/useOAuthClientList'
import { createOAuthClient, getOAuthClient, updateOAuthClient, deleteOAuthClient } from '@/api/oauth_client'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { OAuthClientDetailResponse } from '@/types/oauth_client'

// 头像尺寸（对应 --icon-size-medium）
const avatarSize = 50

// 对话框宽度（对应 --dialog-width-medium）
const dialogWidth = '700px'

// 使用 composable 管理业务逻辑
const {
    loading,
    clientList,
    filters,
    pagination,
    fetchClientList,
    handleFilterChange,
    handlePageChange,
    handleSizeChange
} = useOAuthClientList()

// 弹窗状态
const createDialogVisible = ref(false)
const editDialogVisible = ref(false)
const submitLoading = ref(false)
const createFormRef = ref<InstanceType<typeof OAuthClientForm>>()
const editFormRef = ref<InstanceType<typeof OAuthClientForm>>()
const currentClient = ref<OAuthClientDetailResponse>()

// 打开新建客户端弹窗
const handleCreateClient = () => {
    createDialogVisible.value = true
}

// 提交创建
const handleSubmitCreate = async () => {
    if (!createFormRef.value) return

    // 验证表单
    const isValid = await createFormRef.value.validate()
    if (!isValid) {
        return
    }

    // 获取表单数据并提交
    const formData = createFormRef.value.getFormData()
    submitLoading.value = true

    try {
        await createOAuthClient(formData)
        ElMessage.success('创建 OAuth 客户端成功')
        createDialogVisible.value = false
        // 刷新列表
        await fetchClientList()
    } catch (error: any) {
        // 错误已在拦截器中统一提示，这里不再重复提示
        console.error('创建 OAuth 客户端失败:', error)
    } finally {
        submitLoading.value = false
    }
}

// 取消创建
const handleCancelCreate = () => {
    createDialogVisible.value = false
}

// 创建弹窗关闭时重置表单
const handleCreateDialogClose = () => {
    createFormRef.value?.resetFields()
}

// 打开编辑客户端弹窗
const handleEditClient = async (id: number) => {
    try {
        const response = await getOAuthClient(id)
        currentClient.value = response.data
        editDialogVisible.value = true
    } catch (error: any) {
        // 错误已在拦截器中统一提示
        console.error('获取 OAuth 客户端详情失败:', error)
    }
}

// 提交编辑
const handleSubmitEdit = async () => {
    if (!editFormRef.value || !currentClient.value) return

    // 验证表单
    const isValid = await editFormRef.value.validate()
    if (!isValid) {
        return
    }

    // 获取表单数据并提交
    const formData = editFormRef.value.getFormData()
    submitLoading.value = true

    try {
        await updateOAuthClient(currentClient.value.id, formData)
        ElMessage.success('更新 OAuth 客户端成功')
        editDialogVisible.value = false
        // 刷新列表
        await fetchClientList()
    } catch (error: any) {
        // 错误已在拦截器中统一提示，这里不再重复提示
        console.error('更新 OAuth 客户端失败:', error)
    } finally {
        submitLoading.value = false
    }
}

// 取消编辑
const handleCancelEdit = () => {
    editDialogVisible.value = false
}

// 编辑弹窗关闭时重置数据
const handleEditDialogClose = () => {
    currentClient.value = undefined
}

// 删除客户端
const handleDeleteClient = async (id: number) => {
    try {
        await ElMessageBox.confirm(
            '确定要删除该 OAuth 客户端吗？删除后将无法恢复。',
            '删除确认',
            {
                confirmButtonText: '确定删除',
                cancelButtonText: '取消',
                type: 'warning',
                confirmButtonClass: 'el-button--danger'
            }
        )

        // 用户确认删除
        await deleteOAuthClient(id)
        ElMessage.success('删除 OAuth 客户端成功')
        // 刷新列表
        await fetchClientList()
    } catch (error: any) {
        // 用户取消操作或删除失败
        if (error !== 'cancel' && error !== 'close') {
            ElMessage.error(error.message || '删除 OAuth 客户端失败')
        }
    }
}

onMounted(() => {
    fetchClientList()
})
</script>

<style scoped>
.oauth-clients-page {
    min-height: 100vh;
    background:
        linear-gradient(135deg, rgba(245, 247, 250, 0.8) 0%, rgba(228, 231, 235, 0.9) 100%),
        repeating-linear-gradient(45deg, transparent, transparent var(--pattern-size-small), rgba(0, 0, 0, 0.02) var(--pattern-size-small), rgba(0, 0, 0, 0.02) var(--pattern-size-large)),
        repeating-linear-gradient(-45deg, transparent, transparent var(--pattern-size-small), rgba(0, 0, 0, 0.01) var(--pattern-size-small), rgba(0, 0, 0, 0.01) var(--pattern-size-large)),
        var(--color-background-light);
}

.oauth-clients-page__container {
    min-height: 100vh;
    padding: var(--page-padding-top) var(--spacing-lg) var(--spacing-lg);
    max-width: var(--container-max-width-xlarge);
    margin: 0 auto;
}

.oauth-clients-page__card {
    border-radius: var(--border-radius-card-large);
    box-shadow: var(--shadow-card-layered);
    background: var(--color-card-background);
    border: var(--border-width-thin) solid var(--color-border-white-translucent);
    overflow: hidden;
}

.oauth-clients-page__card :deep(.el-card__header) {
    padding: var(--spacing-lg) var(--spacing-xl);
    border-bottom: var(--border-width-thin) solid var(--color-border-lighter);
    background: var(--color-background-header);
}

.oauth-clients-page__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: var(--spacing-md);
}

.oauth-clients-page__title {
    margin: 0;
    font-size: var(--font-size-title);
    font-weight: 600;
    color: var(--color-text-primary);
}

.oauth-clients-page__actions {
    display: flex;
    align-items: center;
    gap: var(--spacing-md);
}

.oauth-clients-page__filters {
    display: flex;
    gap: var(--spacing-sm-lg);
}

.oauth-clients-page__filter-input {
    width: var(--input-width-medium);
}

.oauth-clients-page__filter-select {
    width: var(--input-width-small);
}

.oauth-clients-page__dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-sm-lg);
}

.oauth-clients-page__card :deep(.el-card__body) {
    padding: var(--spacing-xl);
}

.oauth-clients-page__card :deep(.el-table) {
    border-radius: var(--border-radius-card);
    overflow: hidden;
}

.oauth-clients-page__card :deep(.el-table__header-wrapper th) {
    background-color: var(--color-background-table-header);
    color: var(--color-text-primary);
    font-weight: 600;
}

.oauth-clients-page__card :deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
    background-color: var(--color-background-table-striped);
}

.oauth-clients-page__pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--spacing-lg);
}

.oauth-clients-page__pagination :deep(.el-pagination) {
    justify-content: flex-end;
}

/* 响应式设计 */
/* 平板端：对应 --breakpoint-tablet (768px) */
@media (max-width: 768px) {
    .oauth-clients-page__container {
        padding: var(--page-padding-top) var(--spacing-md) var(--spacing-md);
    }

    .oauth-clients-page__header {
        flex-direction: column;
        align-items: flex-start;
    }

    .oauth-clients-page__actions {
        width: 100%;
        flex-direction: column;
    }

    .oauth-clients-page__filters {
        width: 100%;
        flex-direction: column;
    }

    .oauth-clients-page__filter-input,
    .oauth-clients-page__filter-select {
        width: 100%;
    }

    .oauth-clients-page__card :deep(.el-card__body) {
        padding: var(--spacing-lg);
    }

    .oauth-clients-page__pagination {
        overflow-x: auto;
    }

    .oauth-clients-page__pagination :deep(.el-pagination) {
        flex-wrap: wrap;
    }
}
</style>
