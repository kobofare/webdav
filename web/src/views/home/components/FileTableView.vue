<script setup lang="ts">
import { Delete, Download, Edit, Share, User, View } from '@element-plus/icons-vue'
import type { FileItem } from '../types'

defineProps<{
  rows: FileItem[]
  loading: boolean
  onRowClick: (...args: any[]) => void
  formatSize: (size: number) => string
  formatTime: (time: string | number) => string
  openDetailDrawer: (mode: 'file' | 'recycle', item: FileItem) => void
  downloadFile: (item: FileItem) => void
  shareFile: (item: FileItem) => void
  openShareUserDialog: (item: FileItem) => void
  renameItem: (item: FileItem) => void
  deleteFile: (item: FileItem) => void
}>()
</script>

<template>
  <el-table
    :data="rows"
    v-loading="loading"
    style="width: 100%"
    height="100%"
    @row-click="onRowClick"
  >
    <el-table-column label="名称" min-width="280">
      <template #default="{ row }">
        <div class="file-name">
          <span class="iconfont" :class="row.isDir ? 'icon-wenjianjia' : 'icon-wenjian1'"></span>
          <span class="name" :title="row.name">{{ row.name }}</span>
        </div>
      </template>
    </el-table-column>
    <el-table-column label="大小" width="120">
      <template #default="{ row }">
        <span class="size-cell">{{ row.isDir ? '-' : formatSize(row.size) }}</span>
      </template>
    </el-table-column>
    <el-table-column label="修改时间" width="180">
      <template #default="{ row }">
        <span class="time-cell">{{ formatTime(row.modified) }}</span>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="240" fixed="right">
      <template #default="{ row }">
        <div class="actions" @click.stop>
          <el-tooltip v-if="row.isDir" content="详情" placement="top">
            <el-button type="primary" link :icon="View" @click="openDetailDrawer('file', row)" />
          </el-tooltip>
          <el-tooltip v-if="!row.isDir" content="下载" placement="top">
            <el-button type="primary" link :icon="Download" @click="downloadFile(row)" />
          </el-tooltip>
          <el-tooltip v-if="!row.isDir" content="分享" placement="top">
            <el-button type="primary" link :icon="Share" @click="shareFile(row)" />
          </el-tooltip>
          <el-tooltip content="共享给用户" placement="top">
            <el-button type="primary" link :icon="User" @click="openShareUserDialog(row)" />
          </el-tooltip>
          <el-tooltip content="重命名" placement="top">
            <el-button type="primary" link :icon="Edit" @click="renameItem(row)" />
          </el-tooltip>
          <el-tooltip content="删除" placement="top">
            <el-button type="danger" link :icon="Delete" @click="deleteFile(row)" />
          </el-tooltip>
        </div>
      </template>
    </el-table-column>
  </el-table>
</template>

<style scoped src="./homeShared.scss"></style>
