<template>
  <div class="password-form">
    <div class="password-form__section">
      <div class="password-form__header" @click="showPasswordSection = !showPasswordSection">
        <h3 class="password-form__title">
          <el-icon class="password-form__icon" :class="{ 'password-form__icon--active': showPasswordSection }">
            <ArrowRight />
          </el-icon>
          修改密码
        </h3>
        <span class="password-form__hint">(可选，不修改请留空)</span>
      </div>

      <el-collapse-transition>
        <div v-show="showPasswordSection" class="password-form__fields">
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
  margin-bottom: var(--spacing-lg);
}

.password-form__section {
  padding-bottom: var(--spacing-lg);
  border-bottom: var(--border-width-thin) solid var(--color-border-lighter);
}

.password-form__header {
  cursor: pointer;
  user-select: none;
  margin-bottom: var(--spacing-sm-lg);
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.password-form__header:hover .password-form__title {
  color: var(--color-primary);
}

.password-form__title {
  margin: 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.password-form__hint {
  font-size: var(--font-size-xxs);
  color: var(--color-text-tertiary);
  font-weight: 400;
}

.password-form__icon {
  transition: transform 0.3s ease;
  font-size: var(--font-size-base);
}

.password-form__icon--active {
  transform: rotate(90deg);
}

.password-form__fields {
  padding-top: var(--spacing-sm);
}
</style>
