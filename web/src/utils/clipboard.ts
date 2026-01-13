import { ElMessageBox } from 'element-plus'
import { showSuccess } from './toast'

export async function copyText(text: string, successMessage = '已复制', errorMessage = '复制失败') {
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = text
      textarea.setAttribute('readonly', 'true')
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      document.body.removeChild(textarea)
    }
    showSuccess(successMessage)
  } catch (error) {
    console.error('复制失败:', error)
    void ElMessageBox.alert(errorMessage, '错误', {
      confirmButtonText: '确定',
      type: 'error',
      closeOnClickModal: false
    })
  }
}
