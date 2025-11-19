<template>
    <div class="users-page">
        <Navbar />
        <div class="users-page__container">
            <el-card class="users-page__card">
                <template #header>
                    <div class="users-page__header">
                        <h2 class="users-page__title">用户列表</h2>
                        <div class="users-page__filters">
                            <el-input v-model="filters.nickname" placeholder="搜索昵称" clearable class="users-page__filter-input"
                                @input="handleFilterChange" />
                            <el-select v-model="filters.status" placeholder="状态筛选" clearable class="users-page__filter-select"
                                @change="handleFilterChange">
                                <el-option label="正常" :value="1" />
                                <el-option label="禁用" :value="0" />
                            </el-select>
                            <el-select v-model="filters.role" placeholder="角色筛选" clearable class="users-page__filter-select"
                                @change="handleFilterChange">
                                <el-option label="管理员" value="admin" />
                                <el-option label="普通用户" value="user" />
                            </el-select>
                        </div>
                    </div>
                </template>

                <el-table v-loading="loading" :data="userList" stripe style="width: 100%">
                    <el-table-column prop="id" label="ID" width="80" />
                    <el-table-column label="头像" width="100">
                        <template #default="{ row }">
                            <el-avatar :size="avatarSize" :src="row.avatar" :icon="Avatar" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="nickname" label="昵称" min-width="150" />
                    <el-table-column label="状态" width="100">
                        <template #default="{ row }">
                            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="large">
                                {{ row.status === 1 ? '正常' : '禁用' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column label="角色" width="120">
                        <template #default="{ row }">
                            <el-tag :type="row.role === 'admin' ? 'warning' : 'info'" size="large">
                                {{ row.role === 'admin' ? '管理员' : '普通用户' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="160">
                        <template #default="{ row }">
                            <el-button type="primary" link @click="handleViewUser(row.id)">
                                查看
                            </el-button>
                            <el-button v-if="authStore.user?.id !== row.id" type="danger" link
                                @click="handleDeleteUser(row.id)">
                                删除
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>

                <div class="users-page__pagination">
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
import { useRouter } from 'vue-router'
import Navbar from '@/components/Navbar.vue'
import { useUserList } from '@/composables/useUserList'
import { useAuthStore } from '@/stores/useAuthStore'
import { Avatar } from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()

// 头像尺寸（对应 --icon-size-medium）
const avatarSize = 50

// 使用 composable 管理业务逻辑
const {
    loading,
    userList,
    filters,
    pagination,
    fetchUserList,
    handleFilterChange,
    handlePageChange,
    handleSizeChange,
    handleDeleteUser
} = useUserList()

// 查看用户详情
const handleViewUser = (userId: number) => {
    router.push(`/profile/${userId}`)
}

onMounted(() => {
    fetchUserList()
})
</script>

<style scoped>
.users-page {
    min-height: 100vh;
    background:
        linear-gradient(135deg, rgba(245, 247, 250, 0.8) 0%, rgba(228, 231, 235, 0.9) 100%),
        repeating-linear-gradient(45deg, transparent, transparent var(--pattern-size-small), rgba(0, 0, 0, 0.02) var(--pattern-size-small), rgba(0, 0, 0, 0.02) var(--pattern-size-large)),
        repeating-linear-gradient(-45deg, transparent, transparent var(--pattern-size-small), rgba(0, 0, 0, 0.01) var(--pattern-size-small), rgba(0, 0, 0, 0.01) var(--pattern-size-large)),
        var(--color-background-light);
}

.users-page__container {
    min-height: 100vh;
    padding: var(--page-padding-top) var(--spacing-lg) var(--spacing-lg);
    max-width: var(--container-max-width-xlarge);
    margin: 0 auto;
}

.users-page__card {
    border-radius: var(--border-radius-card-large);
    box-shadow: var(--shadow-card-layered);
    background: var(--color-card-background);
    border: var(--border-width-thin) solid var(--color-border-white-translucent);
    overflow: hidden;
}

.users-page__card :deep(.el-card__header) {
    padding: var(--spacing-lg) var(--spacing-xl);
    border-bottom: var(--border-width-thin) solid var(--color-border-lighter);
    background: var(--color-background-header);
}

.users-page__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: var(--spacing-md);
}

.users-page__title {
    margin: 0;
    font-size: var(--font-size-title);
    font-weight: 600;
    color: var(--color-text-primary);
}

.users-page__filters {
    display: flex;
    gap: var(--spacing-sm-lg);
}

.users-page__filter-input {
    width: var(--input-width-medium);
}

.users-page__filter-select {
    width: var(--input-width-small);
}

.users-page__card :deep(.el-card__body) {
    padding: var(--spacing-xl);
}

.users-page__card :deep(.el-table) {
    border-radius: var(--border-radius-card);
    overflow: hidden;
}

.users-page__card :deep(.el-table__header-wrapper th) {
    background-color: var(--color-background-table-header);
    color: var(--color-text-primary);
    font-weight: 600;
}

.users-page__card :deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
    background-color: var(--color-background-table-striped);
}

.users-page__pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--spacing-lg);
}

.users-page__pagination :deep(.el-pagination) {
    justify-content: flex-end;
}

/* 响应式设计 */
/* 平板端：对应 --breakpoint-tablet (768px) */
@media (max-width: 768px) {
    .users-page__container {
        padding: var(--page-padding-top) var(--spacing-md) var(--spacing-md);
    }

    .users-page__header {
        flex-direction: column;
        align-items: flex-start;
    }

    .users-page__filters {
        width: 100%;
        flex-direction: column;
    }

    .users-page__filter-input,
    .users-page__filter-select {
        width: 100%;
    }

    .users-page__card :deep(.el-card__body) {
        padding: var(--spacing-lg);
    }

    .users-page__pagination {
        overflow-x: auto;
    }

    .users-page__pagination :deep(.el-pagination) {
        flex-wrap: wrap;
    }
}
</style>
