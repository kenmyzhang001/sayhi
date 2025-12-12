<template>
  <div class="speech-groups">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>话术组管理</span>
          <el-button type="primary" @click="handleCreate">新建话术组</el-button>
        </div>
      </template>

      <el-table :data="groups" stripe style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="话术组名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="话术数量" width="120" align="center">
          <template #default="{ row }">
            <el-tag>{{ row.speeches?.length || 0 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="话术内容" min-width="300" show-overflow-tooltip>
          <template #default="{ row }">
            <el-text type="info" size="small">
              {{ row.speeches?.join(', ') || '无' }}
            </el-text>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="话术组名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入话术组名称"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="2"
            placeholder="请输入描述（可选）"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="话术内容" prop="speeches">
          <div class="speeches-editor">
            <el-input
              v-model="speechInput"
              placeholder="输入话术后按回车添加"
              @keyup.enter="handleAddSpeech"
              style="margin-bottom: 10px"
            />
            <div class="speeches-list">
              <el-tag
                v-for="(speech, index) in form.speeches"
                :key="index"
                closable
                @close="handleRemoveSpeech(index)"
                style="margin-right: 8px; margin-bottom: 8px"
              >
                {{ speech }}
              </el-tag>
            </div>
            <el-text type="info" size="small">
              共 {{ form.speeches.length }} 条话术
            </el-text>
          </div>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAllSpeechGroups,
  createSpeechGroup,
  updateSpeechGroup,
  deleteSpeechGroup
} from '../api/api'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新建话术组')
const formRef = ref(null)
const speechInput = ref('')
const groups = ref([])
const editingId = ref(null)

const form = reactive({
  name: '',
  description: '',
  speeches: []
})

const rules = {
  name: [
    { required: true, message: '请输入话术组名称', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  speeches: [
    { required: true, message: '至少添加一条话术', trigger: 'change' },
    {
      validator: (rule, value, callback) => {
        if (!value || value.length === 0) {
          callback(new Error('至少添加一条话术'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
}

// 加载话术组列表
const loadGroups = async () => {
  loading.value = true
  try {
    const data = await getAllSpeechGroups()
    groups.value = data.groups || []
  } catch (error) {
    ElMessage.error('加载失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 创建话术组
const handleCreate = () => {
  dialogTitle.value = '新建话术组'
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑话术组
const handleEdit = (row) => {
  dialogTitle.value = '编辑话术组'
  editingId.value = row.id
  form.name = row.name
  form.description = row.description || ''
  form.speeches = [...(row.speeches || [])]
  dialogVisible.value = true
}

// 删除话术组
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除话术组 "${row.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteSpeechGroup(row.id)
    ElMessage.success('删除成功')
    loadGroups()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败: ' + error.message)
    }
  }
}

// 添加话术
const handleAddSpeech = () => {
  const speech = speechInput.value.trim()
  if (!speech) {
    return
  }

  if (form.speeches.includes(speech)) {
    ElMessage.warning('该话术已存在')
    return
  }

  form.speeches.push(speech)
  speechInput.value = ''
  formRef.value?.validateField('speeches')
}

// 删除话术
const handleRemoveSpeech = (index) => {
  form.speeches.splice(index, 1)
  formRef.value?.validateField('speeches')
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (editingId.value) {
        // 更新
        await updateSpeechGroup(editingId.value, {
          name: form.name,
          description: form.description,
          speeches: form.speeches
        })
        ElMessage.success('更新成功')
      } else {
        // 创建
        await createSpeechGroup({
          name: form.name,
          description: form.description,
          speeches: form.speeches
        })
        ElMessage.success('创建成功')
      }

      dialogVisible.value = false
      loadGroups()
    } catch (error) {
      ElMessage.error(error.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

// 关闭对话框
const handleDialogClose = () => {
  resetForm()
  formRef.value?.clearValidate()
}

// 重置表单
const resetForm = () => {
  form.name = ''
  form.description = ''
  form.speeches = []
  speechInput.value = ''
}

onMounted(() => {
  loadGroups()
})
</script>

<style scoped>
.speech-groups {
  max-width: 1400px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.speeches-editor {
  width: 100%;
}

.speeches-list {
  min-height: 60px;
  padding: 10px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  margin-bottom: 10px;
}
</style>

