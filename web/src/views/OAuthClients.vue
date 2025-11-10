<template>
    <div class="oauth-clients-wrapper">
        <Navbar />
        <div class="oauth-clients-container">
            <el-card class="oauth-clients-card">
                <template #header>
                    <div class="card-header">
                        <h2 class="page-title">OAuth 客户端管理</h2>
                        <div class="filter-bar">
                            <el-input v-model="filters.name" placeholder="搜索客户端名称" clearable style="width: 200px"
                                @input="handleFilterChange" />
                            <el-select v-model="filters.status" placeholder="状态筛选" clearable style="width: 140px"
                                @change="handleFilterChange">
                                <el-option label="正常" :value="1" />
                                <el-option label="禁用" :value="0" />
                            </el-select>
                        </div>
                    </div>
                </template>

                <el-table v-loading="loading" :data="clientList" stripe style="width: 100%">
                    <el-table-column prop="id" label="ID" width="80" />
                    <el-table-column label="Logo" width="100">
                        <template #default="{ row }">
                            <el-avatar :size="50" :src="row.logo || defaultLogo" />
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
                </el-table>

                <div class="pagination-wrapper">
                    <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
                        :page-sizes="[10, 20, 50, 100]" :total="pagination.total"
                        layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange"
                        @current-change="handlePageChange" />
                </div>
            </el-card>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import Navbar from '@/components/Navbar.vue'
import { useOAuthClientList } from '@/composables/useOAuthClientList'

const defaultLogo = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

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

.filter-bar {
    display: flex;
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
