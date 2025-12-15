<template>
  <div
    class="avatar-upload-card"
    :class="{
      'avatar-upload-card--dragging': isDragging,
      'avatar-upload-card--has-preview': !!internalPreviewUrl,
      'avatar-upload-card--disabled': disabled
    }"
    tabindex="0"
    role="button"
    :aria-disabled="disabled"
    :aria-label="internalPreviewUrl ? '已选择头像，可更换、裁剪或移除' : '点击或拖拽上传头像'"
    @dragenter.prevent="onDragEnter"
    @dragover.prevent="onDragOver"
    @dragleave.prevent="onDragLeave"
    @drop.prevent="onDrop"
    @paste="onPaste"
    @click="handleCardClick"
    @keydown="onKeyDown"
  >
    <!-- 空态/拖拽态 -->
    <div v-if="!internalPreviewUrl" class="avatar-upload-card__dropzone">
      <div class="avatar-upload-card__icon-container">
        <el-icon class="avatar-upload-card__icon">
          <Plus />
        </el-icon>
      </div>
      <span class="avatar-upload-card__hint">
        {{ isDragging ? '松开上传' : '上传头像' }}
      </span>
      <span class="avatar-upload-card__sub-hint">
        点击 · 拖拽 · 粘贴
      </span>
    </div>

    <!-- 预览态 -->
    <div v-else class="avatar-upload-card__preview">
      <!-- 圆形头像 -->
      <div class="avatar-upload-card__avatar-wrapper">
        <img
          :src="internalPreviewUrl"
          alt="头像预览"
          class="avatar-upload-card__avatar"
        />
      </div>

      <!-- chips 操作条 -->
      <div class="avatar-upload-card__chips">
        <button
          type="button"
          class="avatar-upload-card__chip"
          :disabled="disabled"
          title="更换"
          @click.stop="handleReplace"
        >
          <el-icon><Refresh /></el-icon>
        </button>
        <button
          type="button"
          class="avatar-upload-card__chip"
          :disabled="disabled"
          title="裁剪"
          @click.stop="handleCrop"
        >
          <el-icon><Crop /></el-icon>
        </button>
        <button
          type="button"
          class="avatar-upload-card__chip avatar-upload-card__chip--danger"
          :disabled="disabled"
          title="移除"
          @click.stop="handleRemove"
        >
          <el-icon><Close /></el-icon>
        </button>
      </div>
    </div>

    <!-- 隐藏的文件输入 -->
    <input
      ref="fileInputRef"
      type="file"
      :accept="allowedTypesStr"
      class="avatar-upload-card__input"
      @change="onFileInputChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Refresh, Crop, Close } from '@element-plus/icons-vue'

// Props
export interface AvatarUploadCardProps {
  /** 当前选中的文件 (v-model) */
  modelValue?: File | null
  /** 外部传入的预览 URL（优先使用内部生成） */
  previewUrl?: string | null
  /** 是否禁用 */
  disabled?: boolean
  /** 最大文件大小（字节），默认 4MB */
  maxSize?: number
  /** 允许的 MIME 类型 */
  allowedTypes?: string[]
}
const props = withDefaults(defineProps<AvatarUploadCardProps>(), {
  modelValue: null,
  previewUrl: null,
  disabled: false,
  maxSize: 4 * 1024 * 1024,
  allowedTypes: () => ['image/png', 'image/jpeg', 'image/jpg', 'image/webp']
})

// Emits
const emit = defineEmits<{
  (e: 'update:modelValue', file: File | null): void
  /** 用户请求裁剪当前图片 */
  (e: 'crop'): void
}>()

// 内部状态
const fileInputRef = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
/** 内部生成的预览 URL */
const internalPreviewUrlRaw = ref<string | null>(null)

// 如果外部传入了 previewUrl，优先使用，否则使用内部生成的
const internalPreviewUrl = computed(() => props.previewUrl ?? internalPreviewUrlRaw.value)

// 允许的类型字符串用于 <input accept>
const allowedTypesStr = computed(() => props.allowedTypes.join(','))

// 清理 URL
const revokePreviewUrl = () => {
  if (internalPreviewUrlRaw.value) {
    URL.revokeObjectURL(internalPreviewUrlRaw.value)
    internalPreviewUrlRaw.value = null
  }
}

onUnmounted(() => {
  revokePreviewUrl()
})

// 当 modelValue 变化时更新预览
watch(
  () => props.modelValue,
  (newFile) => {
    if (newFile) {
      revokePreviewUrl()
      internalPreviewUrlRaw.value = URL.createObjectURL(newFile)
    } else {
      revokePreviewUrl()
    }
  },
  { immediate: true }
)

// 验证并设置文件
const validateAndSetFile = (file: File): boolean => {
  if (props.disabled) return false

  if (!props.allowedTypes.includes(file.type)) {
    ElMessage.warning('仅支持 PNG、JPG、JPEG、WebP 格式的图片')
    return false
  }
  if (file.size > props.maxSize) {
    const sizeMB = (props.maxSize / 1024 / 1024).toFixed(0)
    ElMessage.warning(`头像文件不能超过 ${sizeMB}MB`)
    return false
  }

  emit('update:modelValue', file)
  return true
}

// 提取文件（从 DataTransfer 或 FileList）
const extractFileFromItems = (items: DataTransferItemList | FileList | null): File | null => {
  if (!items) return null

  // DataTransferItemList
  if ('length' in items && typeof (items as DataTransferItemList)[0]?.getAsFile === 'function') {
    const list = items as DataTransferItemList
    for (let i = 0; i < list.length; i++) {
      const item = list[i]
      if (item && item.type.startsWith('image/')) {
        return item.getAsFile()
      }
    }
  }

  // FileList
  const fileList = items as FileList
  if (fileList.length > 0) {
    return fileList[0] ?? null
  }

  return null
}

// 拖拽事件
const onDragEnter = () => {
  if (!props.disabled) isDragging.value = true
}
const onDragOver = () => {
  if (!props.disabled) isDragging.value = true
}
const onDragLeave = () => {
  isDragging.value = false
}
const onDrop = (e: DragEvent) => {
  isDragging.value = false
  if (props.disabled) return

  const file = extractFileFromItems(e.dataTransfer?.items ?? null) ?? e.dataTransfer?.files[0]
  if (file) {
    validateAndSetFile(file)
  }
}

// 粘贴事件
const onPaste = (e: ClipboardEvent) => {
  if (props.disabled) return

  // 如果焦点在输入框内，不拦截
  const target = e.target as HTMLElement
  if (
    target instanceof HTMLInputElement ||
    target instanceof HTMLTextAreaElement ||
    target.isContentEditable
  ) {
    return
  }

  const items = e.clipboardData?.items
  if (!items) return

  // 收集图片项
  const imageFiles: File[] = []
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (item && item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) imageFiles.push(file)
    }
  }

  if (imageFiles.length === 0) return

  e.preventDefault()

  if (imageFiles.length > 1) {
    ElMessage.info('检测到多张图片，已取第一张')
  }

  const firstFile = imageFiles[0]
  if (firstFile) {
    validateAndSetFile(firstFile)
  }
}

// 文件输入变化
const onFileInputChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    const success = validateAndSetFile(file)
    if (!success) {
      target.value = ''
    }
  }
  // 清空 input 让同一文件可再次选择
  target.value = ''
}

// 键盘事件处理（空态回车/空格触发选择）
const onKeyDown = (e: KeyboardEvent) => {
  if (props.disabled) return
  if (e.key === 'Enter' || e.key === ' ') {
    e.preventDefault()
    if (!internalPreviewUrl.value) {
      fileInputRef.value?.click()
    }
  }
}

// 点击卡片触发文件选择（仅在空态）
const handleCardClick = () => {
  if (props.disabled) return
  if (!internalPreviewUrl.value) {
    fileInputRef.value?.click()
  }
}

// 更换图片
const handleReplace = () => {
  if (props.disabled) return
  fileInputRef.value?.click()
}

// 请求裁剪
const handleCrop = () => {
  if (props.disabled) return
  emit('crop')
}

// 移除图片
const handleRemove = () => {
  if (props.disabled) return
  emit('update:modelValue', null)
}
</script>

<style scoped>
/* ==================== 卡片主体 ==================== */
.avatar-upload-card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: var(--avatar-upload-size);
  min-height: var(--avatar-upload-size);
  padding: var(--spacing-md);
  border-radius: var(--border-radius-xlarge);
  border: var(--avatar-upload-border-width) dashed var(--avatar-upload-border);
  background: var(--avatar-upload-surface);
  cursor: pointer;
  transition:
    border-color 0.2s ease,
    background-color 0.2s ease,
    box-shadow 0.2s ease;
  outline: none;
}

.avatar-upload-card:hover:not(.avatar-upload-card--disabled):not(.avatar-upload-card--has-preview) {
  border-color: var(--avatar-upload-border-hover);
}

.avatar-upload-card:focus-visible:not(.avatar-upload-card--disabled) {
  box-shadow: var(--avatar-upload-focus-ring);
  border-color: var(--avatar-upload-border-hover);
}

.avatar-upload-card--dragging {
  border-color: var(--avatar-upload-border-drag);
  background: var(--avatar-upload-background-drag);
  box-shadow: var(--shadow-layer);
}

.avatar-upload-card--has-preview {
  border-style: solid;
  border-color: var(--avatar-upload-border);
  cursor: default;
  padding-bottom: var(--spacing-sm);
}

.avatar-upload-card--disabled {
  opacity: 0.55;
  cursor: not-allowed;
  pointer-events: none;
}

/* ==================== 空态 ==================== */
.avatar-upload-card__dropzone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  width: 100%;
  height: 100%;
}

.avatar-upload-card__icon-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: var(--avatar-upload-icon-container-size);
  height: var(--avatar-upload-icon-container-size);
  border-radius: 50%;
  background: var(--avatar-upload-icon-bg);
  transition: background-color 0.2s ease;
}

.avatar-upload-card:hover:not(.avatar-upload-card--disabled) .avatar-upload-card__icon-container {
  background: var(--avatar-upload-icon-bg-hover);
}

.avatar-upload-card--dragging .avatar-upload-card__icon-container {
  background: var(--avatar-upload-icon-bg-hover);
}

.avatar-upload-card__icon {
  font-size: var(--font-size-title);
  color: var(--avatar-upload-icon-color);
  transition: color 0.2s ease;
}

.avatar-upload-card:hover:not(.avatar-upload-card--disabled) .avatar-upload-card__icon {
  color: var(--avatar-upload-icon-color-hover);
}

.avatar-upload-card--dragging .avatar-upload-card__icon {
  color: var(--avatar-upload-icon-color-drag);
}

.avatar-upload-card__hint {
  font-size: var(--font-size-sm);
  font-weight: 500;
  color: var(--avatar-upload-hint-color);
  text-align: center;
}

.avatar-upload-card__sub-hint {
  font-size: var(--font-size-xs);
  color: var(--avatar-upload-sub-hint-color);
  text-align: center;
}

/* ==================== 预览态 ==================== */
.avatar-upload-card__preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-sm-lg);
  width: 100%;
}

.avatar-upload-card__avatar-wrapper {
  position: relative;
  width: var(--avatar-upload-avatar-size);
  height: var(--avatar-upload-avatar-size);
  border-radius: 50%;
  box-shadow: var(--avatar-upload-avatar-shadow);
  overflow: hidden;
  flex-shrink: 0;
}

.avatar-upload-card__avatar-wrapper::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: var(--avatar-upload-avatar-border);
  pointer-events: none;
  z-index: 1;
}

.avatar-upload-card__avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

/* ==================== Chips 操作条 ==================== */
.avatar-upload-card__chips {
  display: flex;
  align-items: center;
  gap: var(--avatar-upload-chips-gap);
}

.avatar-upload-card__chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  padding: 0;
  border-radius: 50%;
  border: var(--avatar-upload-chip-border);
  background: var(--avatar-upload-chip-bg);
  color: var(--avatar-upload-chip-text);
  font-size: var(--font-size-sm);
  cursor: pointer;
  transition:
    background-color 0.15s ease,
    color 0.15s ease,
    box-shadow 0.15s ease,
    transform 0.1s ease;
  box-shadow: var(--avatar-upload-chip-shadow);
}

.avatar-upload-card__chip:hover:not(:disabled) {
  background: var(--avatar-upload-chip-bg-hover);
  color: var(--avatar-upload-chip-text-hover);
  transform: translateY(-1px);
}

.avatar-upload-card__chip:active:not(:disabled) {
  transform: translateY(0);
}

.avatar-upload-card__chip:focus-visible {
  outline: none;
  box-shadow: var(--avatar-upload-focus-ring), var(--avatar-upload-chip-shadow);
}

.avatar-upload-card__chip:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.avatar-upload-card__chip--danger {
  color: var(--avatar-upload-chip-danger-text);
}

.avatar-upload-card__chip--danger:hover:not(:disabled) {
  color: var(--avatar-upload-chip-danger-text-hover);
}

/* ==================== 隐藏文件输入 ==================== */
.avatar-upload-card__input {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}
</style>
