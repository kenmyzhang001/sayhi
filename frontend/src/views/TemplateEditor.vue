<template>
  <div class="template-editor">
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <span>内容生成</span>
        </div>
      </template>

      <el-form :model="form" label-width="120px">
        <el-form-item label="选择位置">
          <el-checkbox-group v-model="form.selectedPositions">
            <el-checkbox label="a">位置 A</el-checkbox>
            <el-checkbox label="b">位置 B</el-checkbox>
            <el-checkbox label="c">位置 C</el-checkbox>
            <el-checkbox label="d">位置 D</el-checkbox>
          </el-checkbox-group>
          <div class="form-tip">
            <el-text type="info" size="small">
              至少选择一个位置，系统将根据位置配置和话术组生成内容
            </el-text>
          </div>
        </el-form-item>

        <el-form-item label="话术组选择">
          <div class="speech-group-selector">
            <div v-for="pos in ['a', 'b', 'c', 'd']" :key="pos" class="speech-group-item">
              <el-text size="small">位置 {{ pos.toUpperCase() }}:</el-text>
              <el-select
                v-model="form.speechGroups[pos]"
                placeholder="选择话术组（可选）"
                clearable
                style="width: 200px; margin-left: 10px"
              >
                <el-option
                  v-for="group in speechGroups"
                  :key="group.id"
                  :label="group.name"
                  :value="group.name"
                />
              </el-select>
              <el-text v-if="form.speechGroups[pos]" type="success" size="small" style="margin-left: 10px">
                ({{ getGroupSpeechCount(form.speechGroups[pos]) }} 条话术)
              </el-text>
            </div>
          </div>
          <div class="form-tip">
            <el-text type="info" size="small">
              如果为位置选择了话术组，将使用话术组的内容；否则使用位置配置的值
            </el-text>
          </div>
        </el-form-item>

        <el-form-item label="字符编码">
          <el-select v-model="form.encoding" placeholder="请选择编码">
            <el-option label="ASCII" value="ASCII" />
            <el-option label="Zawgyi" value="Zawgyi" />
            <el-option label="Unicode" value="Unicode" />
            <el-option label="其它" value="Other" />
          </el-select>
        </el-form-item>

        <el-form-item label="生成方式">
          <el-radio-group v-model="form.generateMode">
            <el-radio label="sequential">顺序生成</el-radio>
            <el-radio label="random">随机生成</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleGenerate" :loading="loading" :disabled="form.selectedPositions.length === 0">
            生成内容
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="results-card" v-if="results.length > 0">
      <template #header>
        <div class="card-header">
          <span>生成结果</span>
          <div style="display: flex; align-items: center; gap: 15px">
            <el-text type="info" size="small">
              总计: {{ totalCount }} 条 | 
              超出: <el-text type="danger">{{ exceededCount }}</el-text> 条
            </el-text>
            <el-dropdown @command="handleExport" trigger="click">
              <el-button type="primary">
                导出 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="txt">导出为 TXT</el-dropdown-item>
                  <el-dropdown-item command="csv">导出为 CSV</el-dropdown-item>
                  <el-dropdown-item command="excel">导出为 Excel</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </template>

      <el-table :data="results" stripe style="width: 100%">
        <el-table-column prop="content" label="内容" min-width="300">
          <template #default="{ row }">
            <div class="content-cell">
              <el-input
                v-if="editingIndex === row.index"
                v-model="editingContent"
                @blur="handleSaveEdit(row)"
                @keyup.enter="handleSaveEdit(row)"
                autofocus
              />
              <span v-else>{{ row.content }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="charCount" label="字符数" width="100" align="center">
          <template #default="{ row }">
            <el-text :type="row.isExceeded ? 'danger' : 'success'">
              {{ row.charCount }}
            </el-text>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.isExceeded" type="danger">
              超出 {{ row.exceededChars }} 字符
            </el-tag>
            <el-tag v-else type="success">正常</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="handleEdit(row)"
              v-if="editingIndex !== row.index"
            >
              编辑
            </el-button>
            <el-button
              link
              type="success"
              @click="handleSaveEdit(row)"
              v-else
            >
              保存
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { generateTemplate, getAllPositions, getAllSpeechGroups } from '../api/api'

const form = reactive({
  encoding: 'Unicode',
  generateMode: 'sequential',
  selectedPositions: [],
  positions: {
    a: [],
    b: [],
    c: [],
    d: []
  },
  speechGroups: {}
})

const loading = ref(false)
const results = ref([])
const totalCount = ref(0)
const exceededCount = ref(0)
const editingIndex = ref(-1)
const editingContent = ref('')
const speechGroups = ref([])

// 加载位置配置
const loadPositions = async () => {
  try {
    const data = await getAllPositions()
    if (data.positions) {
      form.positions = {
        a: data.positions.a || [],
        b: data.positions.b || [],
        c: data.positions.c || [],
        d: data.positions.d || []
      }
    }
  } catch (error) {
    console.error('加载位置配置失败:', error)
  }
}

// 加载话术组
const loadSpeechGroups = async () => {
  try {
    const data = await getAllSpeechGroups()
    speechGroups.value = data.groups || []
  } catch (error) {
    console.error('加载话术组失败:', error)
  }
}

// 获取话术组的话术数量
const getGroupSpeechCount = (groupName) => {
  const group = speechGroups.value.find(g => g.name === groupName)
  return group ? (group.speeches?.length || 0) : 0
}

// 生成内容
const handleGenerate = async () => {
  if (form.selectedPositions.length === 0) {
    ElMessage.warning('请至少选择一个位置')
    return
  }

  // 验证每个选择的位置是否有配置或话术组
  for (const pos of form.selectedPositions) {
    const hasPositionValue = form.positions[pos] && form.positions[pos].length > 0
    const hasSpeechGroup = form.speechGroups[pos]
    
    if (!hasPositionValue && !hasSpeechGroup) {
      ElMessage.warning(`位置 ${pos.toUpperCase()} 没有配置值或话术组，请先配置`)
      return
    }
  }

  loading.value = true
  try {
    const requestData = {
      encoding: form.encoding,
      generateMode: form.generateMode,
      positions: form.positions,
      selectedPositions: form.selectedPositions,
      speechGroups: form.speechGroups
    }

    const data = await generateTemplate(requestData)
    results.value = data.results.map((item, index) => ({
      ...item,
      index
    }))
    totalCount.value = data.totalCount
    exceededCount.value = data.exceededCount

    if (data.exceededCount > 0) {
      ElMessage.warning(`生成了 ${data.totalCount} 条内容，其中 ${data.exceededCount} 条超出字符限制`)
    } else {
      ElMessage.success(`成功生成 ${data.totalCount} 条内容`)
    }
  } catch (error) {
    ElMessage.error(error.message || '生成失败')
  } finally {
    loading.value = false
  }
}

// 重置表单
const handleReset = () => {
  form.selectedPositions = ['a', 'b', 'c', 'd']
  form.encoding = 'Unicode'
  form.generateMode = 'sequential'
  form.speechGroups = {}
  results.value = []
  totalCount.value = 0
  exceededCount.value = 0
}

// 编辑内容
const handleEdit = (row) => {
  editingIndex.value = row.index
  editingContent.value = row.content
}

// 保存编辑
const handleSaveEdit = (row) => {
  if (editingIndex.value === row.index) {
    row.content = editingContent.value
    // 重新计算字符数
    row.charCount = editingContent.value.length
    row.isExceeded = row.charCount > 70
    row.exceededChars = row.isExceeded ? row.charCount - 70 : 0
    editingIndex.value = -1
    ElMessage.success('保存成功')
  }
}

// 导出功能
const handleExport = (format) => {
  if (results.value.length === 0) {
    ElMessage.warning('没有可导出的数据')
    return
  }

  const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5)
  const filename = `短信内容_${timestamp}`

  switch (format) {
    case 'txt':
      exportToTXT(filename)
      break
    case 'csv':
      exportToCSV(filename)
      break
    case 'excel':
      exportToExcel(filename)
      break
  }
}

// 导出为TXT
const exportToTXT = (filename) => {
  let content = `短信内容生成结果\n`
  content += `生成时间: ${new Date().toLocaleString('zh-CN')}\n`
  content += `总计: ${totalCount.value} 条\n`
  content += `超出限制: ${exceededCount.value} 条\n`
  content += `\n${'='.repeat(50)}\n\n`

  results.value.forEach((item, index) => {
    content += `${index + 1}. ${item.content}\n`
    content += `   字符数: ${item.charCount}`
    if (item.isExceeded) {
      content += ` (超出 ${item.exceededChars} 字符)`
    }
    content += `\n\n`
  })

  downloadFile(content, `${filename}.txt`, 'text/plain;charset=utf-8')
  ElMessage.success('导出成功')
}

// 导出为CSV
const exportToCSV = (filename) => {
  let content = '\uFEFF序号,内容,字符数,是否超出,超出字符数\n'
  
  results.value.forEach((item, index) => {
    const row = [
      index + 1,
      `"${item.content.replace(/"/g, '""')}"`,
      item.charCount,
      item.isExceeded ? '是' : '否',
      item.exceededChars
    ]
    content += row.join(',') + '\n'
  })

  downloadFile(content, `${filename}.csv`, 'text/csv;charset=utf-8')
  ElMessage.success('导出成功')
}

// 导出为Excel (使用CSV格式，Excel可以打开)
const exportToExcel = (filename) => {
  // 创建Excel格式的内容
  let content = '\uFEFF序号\t内容\t字符数\t是否超出\t超出字符数\n'
  
  results.value.forEach((item, index) => {
    const row = [
      index + 1,
      item.content,
      item.charCount,
      item.isExceeded ? '是' : '否',
      item.exceededChars
    ]
    content += row.join('\t') + '\n'
  })

  downloadFile(content, `${filename}.xls`, 'application/vnd.ms-excel;charset=utf-8')
  ElMessage.success('导出成功')
}

// 下载文件
const downloadFile = (content, filename, mimeType) => {
  const blob = new Blob([content], { type: mimeType })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

onMounted(() => {
  loadPositions()
  loadSpeechGroups()
})
</script>

<style scoped>
.template-editor {
  max-width: 1200px;
  margin: 0 auto;
}

.editor-card,
.results-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  margin-top: 5px;
}

.speech-group-selector {
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.speech-group-item {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
}

.content-cell {
  word-break: break-all;
}

:deep(.el-table) {
  font-size: 14px;
}
</style>
