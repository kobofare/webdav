export interface FileItem {
  name: string
  path: string
  isDir: boolean
  size: number
  modified: string
}

export interface UploadItem {
  file: File
  relativePath: string
}

export type DropEntry = {
  isFile: boolean
  isDirectory: boolean
  name: string
  fullPath?: string
  file?: (success: (file: File) => void, error?: (error: Error) => void) => void
  createReader?: () => {
    readEntries: (success: (entries: DropEntry[]) => void, error?: (error: Error) => void) => void
  }
}
