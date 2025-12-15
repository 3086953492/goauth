<template>
  <el-dialog v-model="dialogVisible" title="裁剪头像" width="420px" :close-on-click-modal="false"
    :close-on-press-escape="false" class="avatar-cropper-dialog" @closed="onDialogClosed">
    <div class="avatar-cropper-dialog__body">
      <div class="avatar-cropper-dialog__cropper-wrapper">
        <Cropper ref="cropperRef" class="avatar-cropper-dialog__cropper" :src="imageSrc"
          :stencil-component="CircleStencil" :stencil-props="{ aspectRatio: 1 }" :canvas="canvasOptions"
          background-class="avatar-cropper-dialog__background" image-restriction="stencil" />
      </div>

      <!-- 控制按钮 -->
      <div class="avatar-cropper-dialog__controls">
        <el-button-group>
          <el-button @click="rotate(-90)" title="逆时针旋转">
            <el-icon>
              <RefreshLeft />
            </el-icon>
          </el-button>
          <el-button @click="rotate(90)" title="顺时针旋转">
            <el-icon>
              <RefreshRight />
            </el-icon>
          </el-button>
        </el-button-group>

        <el-button-group>
          <el-button @click="zoom(0.9)" title="缩小">
            <el-icon>
              <ZoomOut />
            </el-icon>
          </el-button>
          <el-button @click="zoom(1.1)" title="放大">
            <el-icon>
              <ZoomIn />
            </el-icon>
          </el-button>
        </el-button-group>

        <el-button @click="reset" title="重置">
          <el-icon>
            <Refresh />
          </el-icon>
        </el-button>
      </div>
    </div>

    <template #footer>
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" :loading="exporting" @click="handleConfirm">
        {{ exporting ? '处理中...' : '确定' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Cropper, CircleStencil } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'
import { RefreshLeft, RefreshRight, ZoomOut, ZoomIn, Refresh } from '@element-plus/icons-vue'

// Props
export interface AvatarCropperDialogProps {
  /** 控制弹窗开关 (v-model) */
  modelValue: boolean
  /** 待裁剪的文件 */
  file: File | null
  /** 输出文件名（不含扩展名），默认 avatar */
  outputName?: string
  /** 输出格式 */
  outputType?: 'image/webp' | 'image/png' | 'image/jpeg'
  /** 输出质量 (0-1)，仅对 webp/jpeg 有效 */
  outputQuality?: number
  /** 最大输出尺寸 */
  maxOutputSize?: number
}

const props = withDefaults(defineProps<AvatarCropperDialogProps>(), {
  file: null,
  outputName: 'avatar',
  outputType: 'image/webp',
  outputQuality: 0.9,
  maxOutputSize: 512
})

// Emits
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm', file: File): void
}>()

// 弹窗显隐双向绑定
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// 图片数据 URL
const imageSrc = ref<string>('')
// 当前 file 改变时读取为 data URL
watch(
  () => props.file,
  async (newFile) => {
    if (newFile) {
      imageSrc.value = await fileToDataUrl(newFile)
    } else {
      imageSrc.value = ''
    }
  },
  { immediate: true }
)

// 裁剪器 ref
const cropperRef = ref<InstanceType<typeof Cropper> | null>(null)

// canvas 导出配置
const canvasOptions = computed(() => ({
  maxWidth: props.maxOutputSize,
  maxHeight: props.maxOutputSize
}))

// 文件转 data URL
const fileToDataUrl = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

// 旋转
const rotate = (angle: number) => {
  cropperRef.value?.rotate(angle)
}

// 缩放
const zoom = (factor: number) => {
  cropperRef.value?.zoom(factor)
}

// 重置
const reset = () => {
  cropperRef.value?.reset()
}

// 导出状态
const exporting = ref(false)

// 确认裁剪
const handleConfirm = async () => {
  if (!cropperRef.value) return

  exporting.value = true
  try {
    const { canvas } = cropperRef.value.getResult()
    if (!canvas) {
      exporting.value = false
      return
    }

    // 导出 blob
    const blob = await new Promise<Blob | null>((resolve) => {
      canvas.toBlob(resolve, props.outputType, props.outputQuality)
    })

    if (!blob) {
      exporting.value = false
      return
    }

    // 构造 File
    const ext = props.outputType.split('/')[1] ?? 'webp'
    const fileName = `${props.outputName}.${ext}`
    const croppedFile = new File([blob], fileName, { type: props.outputType })

    emit('confirm', croppedFile)
    dialogVisible.value = false
  } finally {
    exporting.value = false
  }
}

// 取消
const handleCancel = () => {
  dialogVisible.value = false
}

// 弹窗关闭后清理
const onDialogClosed = () => {
  imageSrc.value = ''
}
</script>

<style scoped>
.avatar-cropper-dialog__body {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.avatar-cropper-dialog__cropper-wrapper {
  width: 100%;
  aspect-ratio: 1 / 1;
  border-radius: var(--border-radius-large);
  overflow: hidden;
  background: var(--color-background-lighter);
}

.avatar-cropper-dialog__cropper {
  width: 100%;
  height: 100%;
}

.avatar-cropper-dialog__background {
  background: var(--color-background-lighter);
}

.avatar-cropper-dialog__controls {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}
</style>
