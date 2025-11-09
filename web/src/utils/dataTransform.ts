/**
 * 数据转换和比较工具函数
 */

/**
 * 检测对象中指定字段的变化
 * @param original 原始数据对象
 * @param current 当前数据对象
 * @param fields 需要检测的字段列表
 * @returns 返回包含变化字段的对象
 */
export function detectChanges<T extends Record<string, any>>(
  original: T,
  current: T,
  fields: (keyof T)[]
): Partial<T> {
  const changes: Partial<T> = {}

  for (const field of fields) {
    const originalValue = original[field]
    const currentValue = current[field]

    // 处理空字符串和 undefined 的情况
    const normalizedOriginal = originalValue === '' ? undefined : originalValue
    const normalizedCurrent = currentValue === '' ? undefined : currentValue

    if (normalizedOriginal !== normalizedCurrent) {
      changes[field] = normalizedCurrent
    }
  }

  return changes
}

/**
 * 深度比较两个对象是否相等
 * @param obj1 对象1
 * @param obj2 对象2
 * @returns 是否相等
 */
export function deepEqual(obj1: any, obj2: any): boolean {
  if (obj1 === obj2) return true

  if (typeof obj1 !== 'object' || typeof obj2 !== 'object' || obj1 === null || obj2 === null) {
    return false
  }

  const keys1 = Object.keys(obj1)
  const keys2 = Object.keys(obj2)

  if (keys1.length !== keys2.length) return false

  for (const key of keys1) {
    if (!keys2.includes(key) || !deepEqual(obj1[key], obj2[key])) {
      return false
    }
  }

  return true
}

/**
 * 移除对象中的空值字段（undefined, null, 空字符串）
 * @param obj 源对象
 * @returns 清理后的对象
 */
export function removeEmptyFields<T extends Record<string, any>>(obj: T): Partial<T> {
  const result: Partial<T> = {}

  for (const key in obj) {
    const value = obj[key]
    if (value !== undefined && value !== null && value !== '') {
      result[key] = value
    }
  }

  return result
}

/**
 * 安全地转换 avatar 字段，空字符串转为 undefined
 * @param avatar 头像 URL
 * @returns 转换后的值
 */
export function normalizeAvatar(avatar: string | undefined): string | undefined {
  return avatar && avatar.trim() !== '' ? avatar : undefined
}

