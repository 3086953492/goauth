/**
 * 认证反馈事件总线
 * 用于解耦业务逻辑层与 UI 反馈层，业务层只负责发出语义化事件，UI 层统一处理提示和跳转
 */

/**
 * 认证事件类型
 */
export type AuthEventType = 
  | 'auth:login-required'   // 未登录，需要登录
  | 'auth:login-expired'    // 登录已过期
  | 'auth:forbidden'        // 权限不足

/**
 * 认证事件负载
 */
export interface AuthEventPayload {
  /** 提示消息 */
  message: string
  /** 重定向路径（登录后返回的路径或权限不足时跳转的安全路径） */
  redirectPath?: string
  /** 错误码（可选，用于 403 等场景） */
  errorCode?: string
  /** 错误描述（可选） */
  errorDescription?: string
}

/**
 * 认证事件处理器类型
 */
export type AuthEventHandler = (type: AuthEventType, payload: AuthEventPayload) => void

/**
 * 事件总线实现
 */
class AuthFeedbackBus {
  private handlers: Set<AuthEventHandler> = new Set()

  /**
   * 订阅认证事件
   */
  on(handler: AuthEventHandler): void {
    this.handlers.add(handler)
  }

  /**
   * 取消订阅认证事件
   */
  off(handler: AuthEventHandler): void {
    this.handlers.delete(handler)
  }

  /**
   * 发出认证事件
   */
  emit(type: AuthEventType, payload: AuthEventPayload): void {
    this.handlers.forEach(handler => {
      try {
        handler(type, payload)
      } catch (error) {
        console.error('[AuthFeedbackBus] Handler error:', error)
      }
    })
  }
}

// 导出单例实例
const authFeedbackBus = new AuthFeedbackBus()

/**
 * 发出认证事件（便捷方法）
 */
export function emitAuthEvent(type: AuthEventType, payload: AuthEventPayload): void {
  authFeedbackBus.emit(type, payload)
}

/**
 * 订阅认证事件（便捷方法）
 */
export function onAuthEvent(handler: AuthEventHandler): void {
  authFeedbackBus.on(handler)
}

/**
 * 取消订阅认证事件（便捷方法）
 */
export function offAuthEvent(handler: AuthEventHandler): void {
  authFeedbackBus.off(handler)
}

export default authFeedbackBus

