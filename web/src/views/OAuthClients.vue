<template>
    <div class="oauth-clients-wrapper">
        <Navbar />
        <div class="oauth-clients-container">
            <el-card class="oauth-clients-card">
                <template #header>
                    <div class="card-header">
                        <h2 class="page-title">OAuth 客户端管理</h2>
                        <div class="header-actions">
                            <div class="filter-bar">
                                <el-input v-model="filters.name" placeholder="搜索客户端名称" clearable style="width: 200px"
                                    @input="handleFilterChange" />
                                <el-select v-model="filters.status" placeholder="状态筛选" clearable style="width: 140px"
                                    @change="handleFilterChange">
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
                            <el-avatar :size="50" :src="row.logo" :icon="Cpu" />
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

                <div class="pagination-wrapper">
                    <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
                        :page-sizes="[10, 20, 50, 100]" :total="pagination.total"
                        layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange"
                        @current-change="handlePageChange" />
                </div>
            </el-card>
        </div>

        <!-- 新建客户端弹窗 -->
        <el-dialog v-model="createDialogVisible" title="新建 OAuth 客户端" width="700px" :close-on-click-modal="false"
            @close="handleCreateDialogClose">
            <OAuthClientForm ref="createFormRef" mode="create" />
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="handleCancelCreate">取消</el-button>
                    <el-button type="primary" :loading="submitLoading" @click="handleSubmitCreate">
                        确认创建
                    </el-button>
                </div>
            </template>
        </el-dialog>

        <!-- 编辑客户端弹窗 -->
        <el-dialog v-model="editDialogVisible" title="编辑 OAuth 客户端" width="700px" :close-on-click-modal="false"
            @close="handleEditDialogClose">
            <OAuthClientForm ref="editFormRef" mode="edit" :initial-data="currentClient" />
            <template #footer>
                <div class="dialog-footer">
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
        const response = await createOAuthClient(formData)
        if (response.success) {
            ElMessage.success('创建 OAuth 客户端成功')
            createDialogVisible.value = false
            // 刷新列表
            await fetchClientList()
        } else {
            ElMessage.error(response.message || '创建 OAuth 客户端失败')
        }
    } catch (error: any) {
        ElMessage.error(error.message || '创建 OAuth 客户端失败')
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
        if (response.success && response.data) {
            currentClient.value = response.data
            editDialogVisible.value = true
        } else {
            ElMessage.error(response.message || '获取 OAuth 客户端详情失败')
        }
    } catch (error: any) {
        ElMessage.error(error.message || '获取 OAuth 客户端详情失败')
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
        const response = await updateOAuthClient(currentClient.value.id, formData)
        if (response.success) {
            ElMessage.success('更新 OAuth 客户端成功')
            editDialogVisible.value = false
            // 刷新列表
            await fetchClientList()
        } else {
            ElMessage.error(response.message || '更新 OAuth 客户端失败')
        }
    } catch (error: any) {
        ElMessage.error(error.message || '更新 OAuth 客户端失败')
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
        const response = await deleteOAuthClient(id)
        if (response.success) {
            ElMessage.success('删除 OAuth 客户端成功')
            // 刷新列表
            await fetchClientList()
        } else {
            ElMessage.error(response.message || '删除 OAuth 客户端失败')
        }
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
.oauth-clients-wrapper {
    min-height: 100vh;
    background:
        linear-gradient(135deg, rgba(245, 247, 250, 0.8) 0%, rgba(228, 231, 235, 0.9) 100%),
        repeating-linear-gradient(45deg, transparent, transparent 35px, rgba(0, 0, 0, 0.02) 35px, rgba(0, 0, 0, 0.02) 70px),
        repeating-linear-gradient(-45deg, transparent, transparent 35px, rgba(0, 0, 0, 0.01) 35px, rgba(0, 0, 0, 0.01) 70px),
        #f8f9fa;
}

.oauth-clients-container {
    min-height: 100vh;
    padding: 84px 20px 20px;
    max-width: 1400px;
    margin: 0 auto;
}

.oauth-clients-card {
    border-radius: 24px;
    box-shadow:
        0 2px 8px rgba(0, 0, 0, 0.04),
        0 8px 24px rgba(0, 0, 0, 0.06),
        0 16px 48px rgba(0, 0, 0, 0.08);
    background: #ffffff;
    border: 1px solid rgba(255, 255, 255, 0.8);
    overflow: hidden;
}

:deep(.el-card__header) {
    padding: 24px 32px;
    border-bottom: 1px solid #ebeef5;
    background: #fafafa;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 16px;
}

.page-title {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
    color: #303133;
}

.header-actions {
    display: flex;
    align-items: center;
    gap: 16px;
}

.filter-bar {
    display: flex;
    gap: 12px;
}

.dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
}

:deep(.el-card__body) {
    padding: 32px;
}

:deep(.el-table) {
    border-radius: 12px;
    overflow: hidden;
}

:deep(.el-table__header-wrapper th) {
    background-color: #f5f7fa;
    color: #303133;
    font-weight: 600;
}

:deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
    background-color: #fafafa;
}

.pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 24px;
}

:deep(.el-pagination) {
    justify-content: flex-end;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .oauth-clients-container {
        padding: 84px 16px 16px;
    }

    .card-header {
        flex-direction: column;
        align-items: flex-start;
    }

    .header-actions {
        width: 100%;
        flex-direction: column;
    }

    .filter-bar {
        width: 100%;
        flex-direction: column;
    }

    .filter-bar .el-input,
    .filter-bar .el-select {
        width: 100% !important;
    }

    :deep(.el-card__body) {
        padding: 20px;
    }

    .pagination-wrapper {
        overflow-x: auto;
    }

    :deep(.el-pagination) {
        flex-wrap: wrap;
    }
}
</style>
