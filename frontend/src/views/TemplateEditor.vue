<template>
  <div class="template-editor">
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <span>模板编辑</span>
        </div>
      </template>

      <el-form :model="form" label-width="120px">
        <el-form-item label="模板内容">
          <el-input
            v-model="form.template"
            type="textarea"
            :rows="3"
            placeholder="例如: (1)(baidu.com)(2)(3-10)"
            @input="handleTemplateChange"
          />
          <div class="form-tip">
            <el-text type="info" size="small">
              使用括号 () 定义位置，支持固定值和范围值（如 3-10）
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
          <el-button type="primary" @click="handleGenerate" :loading="loading">
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
          <div>
            <el-text type="info" size="small">
              总计: {{ totalCount }} 条 | 
              超出: <el-text type="danger">{{ exceededCount }}</el-text> 条
            </el-text>
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
import { generateTemplate, getAllPositions } from '../api/api'

const form = reactive({
  template: '(1)(baidu.com)(2)(3-10)',
  encoding: 'Unicode',
  generateMode: 'sequential',
  positions: {
    a: [],
    b: [],
    c: [],
    d: []
  }
})

const loading = ref(false)
const results = ref([])
const totalCount = ref(0)
const exceededCount = ref(0)
const editingIndex = ref(-1)
const editingContent = ref('')

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

// 处理模板变化
const handleTemplateChange = () => {
  // 可以在这里添加实时验证
}

// 生成内容
const handleGenerate = async () => {
  if (!form.template.trim()) {
    ElMessage.warning('请输入模板内容')
    return
  }

  loading.value = true
  try {
    const data = await generateTemplate(form)
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
  form.template = ''
  form.encoding = 'Unicode'
  form.generateMode = 'sequential'
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

onMounted(() => {
  loadPositions()
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

.content-cell {
  word-break: break-all;
}

:deep(.el-table) {
  font-size: 14px;
}
</style>

