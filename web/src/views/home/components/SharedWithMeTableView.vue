<script setup lang="ts">
import { computed } from 'vue'
import { Delete, Download, Edit, View } from '@element-plus/icons-vue'
import type { DirectShareItem } from '@/api'
import type { FileItem } from '../types'

const props = defineProps<{
  sharedActive: DirectShareItem | null
  sharedWithMeList: DirectShareItem[]
  sharedEntries: FileItem[]
  loading: boolean
  onRowClick: (...args: any[]) => void
  formatTime: (time: string | number) => string
  formatSize: (size: number) => string
  shortenAddress: (address?: string) => string
  sharedCanRead: boolean
  sharedCanUpdate: boolean
  sharedCanDelete: boolean
  openShareDetail: (mode: 'share' | 'directShare' | 'receivedShare', item: DirectShareItem) => void
  downloadSharedRoot: (item: DirectShareItem) => void
  openSharedEntryDetail: (item: FileItem) => void
  downloadSharedFile: (item: FileItem) => void
  renameSharedItem: (item: FileItem) => void
  deleteSharedItem: (item: FileItem) => void
}>()

const tableRows = computed(() => (props.sharedActive ? props.sharedEntries : props.sharedWithMeList))
</script>

<template>
  <el-table
    :data="tableRows"
    v-loading="loading"
    style="width: 100%"
    height="100%"
    @row-click="onRowClick"
  >
    <template v-if="!sharedActive">
      <el-table-column label="名称" min-width="200">
        <template #default="{ row }">
          <div class="file-name">
            <span class="iconfont" :class="row.isDir ? 'icon-wenjianjia' : 'icon-wenjian1'"></span>
            <span class="name" :title="row.name">{{ row.name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="分享人" min-width="160">
        <template #default="{ row }">
          <span>{{ row.ownerName || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="源钱包" min-width="200">
        <template #default="{ row }">
          <span class="mono">{{ shortenAddress(row.ownerWallet) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="过期时间" width="180">
        <template #default="{ row }">
          <span class="time-cell">{{ row.expiresAt ? formatTime(row.expiresAt) : '永不过期' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          <span class="time-cell">{{ row.createdAt ? formatTime(row.createdAt) : '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <div class="actions" @click.stop>
            <el-tooltip v-if="row.isDir" content="详情" placement="top">
              <el-button type="primary" link :icon="View" @click="openShareDetail('receivedShare', row)" />
            </el-tooltip>
            <el-tooltip v-else-if="row.permissions && row.permissions.includes('read')" content="下载" placement="top">
              <el-button type="primary" link :icon="Download" @click="downloadSharedRoot(row)" />
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
    </template>
    <template v-else>
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
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <div class="actions" @click.stop>
            <el-tooltip v-if="row.isDir" content="详情" placement="top">
              <el-button type="primary" link :icon="View" @click="openSharedEntryDetail(row)" />
            </el-tooltip>
            <el-tooltip v-if="!row.isDir && sharedCanRead" content="下载" placement="top">
              <el-button type="primary" link :icon="Download" @click="downloadSharedFile(row)" />
            </el-tooltip>
            <el-tooltip v-if="sharedCanUpdate" content="重命名" placement="top">
              <el-button type="primary" link :icon="Edit" @click="renameSharedItem(row)" />
            </el-tooltip>
            <el-tooltip v-if="sharedCanDelete" content="删除" placement="top">
              <el-button type="danger" link :icon="Delete" @click="deleteSharedItem(row)" />
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
    </template>
  </el-table>
</template>

<style scoped src="./homeShared.scss"></style>
