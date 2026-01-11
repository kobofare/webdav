export interface FileItem {
  name: string
  path: string
  isDir: boolean
  size: number
  modified: string
}

// 规范化路径用于比较（确保统一格式）
function normalizePathForCompare(path: string): string {
  if (path === '/') return '/'
  // 确保以 / 开头和结尾
  let normalized = path.startsWith('/') ? path : '/' + path
  return normalized.endsWith('/') ? normalized : normalized + '/'
}

// 检查两个路径是否表示同一目录
function isSamePath(path1: string, path2: string): boolean {
  return normalizePathForCompare(path1) === normalizePathForCompare(path2)
}

// 解析 WebDAV PROPFIND 响应
export function parsePropfindResponse(xml: string, currentPath: string): FileItem[] {
  const items: FileItem[] = []

  // 按 <d:response> 分割，每个 response 块包含一个资源的信息
  const responseRegex = /<[Dd]:response[^>]*>([\s\S]*?)<\/[Dd]:response>/gi
  const responses = [...xml.matchAll(responseRegex)]

  console.log('PROPFIND: responses count:', responses.length)

  // 规范化当前路径用于比较
  const normalizedCurrentPath = normalizePathForCompare(currentPath)
  console.log('PROPFIND: currentPath:', currentPath, 'normalized:', normalizedCurrentPath)

  for (const match of responses) {
    const responseXml = match[1]

    // 提取 href
    const hrefMatch = /<[Dd]:href>([^<]+)<\/[Dd]:href>/i.exec(responseXml)
    if (!hrefMatch) continue
    const href = decodeURIComponent(hrefMatch[1])
    console.log('PROPFIND: href:', href)

    // 排除根目录自身和当前目录自身
    if (isSamePath(href, normalizedCurrentPath)) {
      console.log('PROPFIND: skip (same path)')
      continue
    }

    // 提取 displayname
    const nameMatch = /<[Dd]:displayname>([^<]*?)<\/[Dd]:displayname>/i.exec(responseXml)
    let name = nameMatch?.[1] || ''

    // 如果 displayname 为空，从 href 提取名称
    if (!name) {
      name = href.split('/').filter(Boolean).pop() || ''
    }

    if (name === '') continue

    // 提取文件大小（目录可能没有这个属性）
    const sizeMatch = /<[Dd]:getcontentlength>([^<]+)<\/[Dd]:getcontentlength>/i.exec(responseXml)
    const size = parseInt(sizeMatch?.[1] || '0')

    // 提取修改时间
    const lastModMatch = /<[Dd]:getlastmodified>([^<]+)<\/[Dd]:getlastmodified>/i.exec(responseXml)
    const lastMod = lastModMatch?.[1] || ''

    items.push({
      name,
      path: href,
      isDir: href.endsWith('/'),
      size,
      modified: lastMod
    })
  }

  return items
}