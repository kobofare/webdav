import { ElMessage } from 'element-plus'

const MESSAGE_OFFSET = 16

export function showSuccess(message: string): void {
  ElMessage({ message, type: 'success', offset: MESSAGE_OFFSET })
}

export function showInfo(message: string): void {
  ElMessage({ message, type: 'info', offset: MESSAGE_OFFSET })
}
