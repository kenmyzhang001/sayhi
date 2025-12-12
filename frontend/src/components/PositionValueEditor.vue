<template>
  <div class="position-value-editor">
    <el-form>
      <el-form-item label="添加新值">
        <div class="add-value-form">
          <el-input
            v-model="newValue"
            placeholder="请输入位置值"
            style="width: 300px; margin-right: 10px"
            @keyup.enter="handleAdd"
          />
          <el-button type="primary" @click="handleAdd">添加</el-button>
        </div>
      </el-form-item>
    </el-form>

    <el-divider />

    <div class="values-list">
      <el-empty v-if="displayValues.length === 0" description="暂无配置值" />
      <el-table v-else :data="displayValues" stripe style="width: 100%">
        <el-table-column prop="value" label="值" min-width="200">
          <template #default="{ row }">
            <el-input
              v-if="row.editing"
              v-model="row.editValue"
              @blur="handleSaveEdit(row)"
              @keyup.enter="handleSaveEdit(row)"
              autofocus
            />
            <span v-else>{{ row.value }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button
              v-if="!row.editing"
              link
              type="primary"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              v-else
              link
              type="success"
              @click="handleSaveEdit(row)"
            >
              保存
            </el-button>
            <el-button
              link
              type="danger"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="batch-edit" style="margin-top: 20px">
      <el-form-item label="批量设置">
        <el-input
          v-model="batchValues"
          type="textarea"
          :rows="5"
          placeholder="每行一个值，用于批量设置"
          style="width: 500px"
        />
        <div style="margin-top: 10px">
          <el-button type="primary" @click="handleBatchSet">批量设置</el-button>
          <el-button @click="batchValues = ''">清空</el-button>
        </div>
      </el-form-item>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { addPositionValue, setPositionValues, deletePositionValue } from '../api/api'

const props = defineProps({
  position: {
    type: String,
    required: true
  },
  values: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update'])

const newValue = ref('')
const batchValues = ref('')
const displayValues = ref([])

// 初始化显示值
const initDisplayValues = () => {
  displayValues.value = props.values.map((val, index) => ({
    value: val,
    editValue: val,
    editing: false,
    index
  }))
}

// 监听 values 变化
watch(() => props.values, () => {
  initDisplayValues()
}, { immediate: true, deep: true })

// 添加值
const handleAdd = async () => {
  const valueToAdd = newValue.value.trim()
  if (!valueToAdd) {
    ElMessage.warning('请输入值')
    return
  }

  try {
    await addPositionValue({
      position: props.position,
      value: valueToAdd
    })
    ElMessage.success('添加成功')
    newValue.value = ''
    emit('update', props.position, null) // 触发刷新
  } catch (error) {
    ElMessage.error('添加失败: ' + error.message)
  }
}

// 编辑值
const handleEdit = (row) => {
  row.editing = true
  row.editValue = row.value
}

// 保存编辑
const handleSaveEdit = async (row) => {
  if (!row.editValue.trim()) {
    ElMessage.warning('值不能为空')
    return
  }

  const oldValue = row.value
  const newVal = row.editValue.trim()

  if (oldValue === newVal) {
    row.editing = false
    return
  }

  try {
    // 先删除旧值，再添加新值
    await deletePositionValue(props.position, oldValue)
    await addPositionValue({
      position: props.position,
      value: newVal
    })
    
    row.value = newVal
    row.editing = false
    ElMessage.success('保存成功')
    emit('update', props.position, null) // 触发刷新
  } catch (error) {
    ElMessage.error('保存失败: ' + error.message)
  }
}

// 删除值
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除值 "${row.value}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deletePositionValue(props.position, row.value)
    ElMessage.success('删除成功')
    emit('update', props.position, null) // 触发刷新
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败: ' + error.message)
    }
  }
}

// 批量设置
const handleBatchSet = async () => {
  if (!batchValues.value.trim()) {
    ElMessage.warning('请输入要设置的值')
    return
  }

  const values = batchValues.value
    .split('\n')
    .map(v => v.trim())
    .filter(v => v.length > 0)

  if (values.length === 0) {
    ElMessage.warning('没有有效的值')
    return
  }

  try {
    await setPositionValues(props.position, values)
    ElMessage.success(`批量设置成功，共 ${values.length} 个值`)
    batchValues.value = ''
    emit('update', props.position, values)
  } catch (error) {
    ElMessage.error('批量设置失败: ' + error.message)
  }
}
</script>

<style scoped>
.position-value-editor {
  padding: 20px;
}

.add-value-form {
  display: flex;
  align-items: center;
}

.values-list {
  margin-top: 20px;
}
</style>

