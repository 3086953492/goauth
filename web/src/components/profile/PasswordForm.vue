<template>
  <div class="password-form">
    <div class="form-section">
      <div class="section-header" @click="showPasswordSection = !showPasswordSection">
        <h3 class="section-title">
          <el-icon class="collapse-icon" :class="{ 'is-active': showPasswordSection }">
            <ArrowRight />
          </el-icon>
          修改密码
        </h3>
        <span class="section-hint">(可选，不修改请留空)</span>
      </div>

      <el-collapse-transition>
        <div v-show="showPasswordSection" class="password-fields">
          <el-form-item label="新密码" prop="password">
            <el-input :model-value="modelValue.password" @input="updateField('password', $event)" type="password"
              placeholder="留空表示不修改密码" show-password clearable :prefix-icon="Lock" />
          </el-form-item>

          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input :model-value="modelValue.confirmPassword" @input="updateField('confirmPassword', $event)"
              type="password" placeholder="请再次输入新密码" show-password clearable :prefix-icon="Lock" />
          </el-form-item>
        </div>
      </el-collapse-transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Lock, ArrowRight } from '@element-plus/icons-vue'

interface PasswordData {
  password: string
  confirmPassword: string
}

interface Props {
  modelValue: PasswordData
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:modelValue', value: PasswordData): void
}>()

const showPasswordSection = ref(false)

const updateField = (field: keyof PasswordData, value: string) => {
  emit('update:modelValue', {
    ...props.modelValue,
    [field]: value
  })
}
</script>

<style scoped>
.password-form {
  margin-bottom: 24px;
}

.form-section {
  padding-bottom: 24px;
  border-bottom: 1px solid #ebeef5;
}

.section-header {
  cursor: pointer;
  user-select: none;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-header:hover .section-title {
  color: #409eff;
}

.section-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-hint {
  font-size: 13px;
  color: #909399;
  font-weight: 400;
}

.collapse-icon {
  transition: transform 0.3s ease;
  font-size: 16px;
}

.collapse-icon.is-active {
  transform: rotate(90deg);
}

.password-fields {
  padding-top: 8px;
}
</style>
