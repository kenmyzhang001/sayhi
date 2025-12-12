<template>
  <div class="position-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>位置值配置</span>
          <el-button type="primary" @click="handleRefresh">刷新</el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="位置 A" name="a">
          <PositionValueEditor position="a" :values="positions.a" @update="handleUpdate" />
        </el-tab-pane>
        <el-tab-pane label="位置 B" name="b">
          <PositionValueEditor position="b" :values="positions.b" @update="handleUpdate" />
        </el-tab-pane>
        <el-tab-pane label="位置 C" name="c">
          <PositionValueEditor position="c" :values="positions.c" @update="handleUpdate" />
        </el-tab-pane>
        <el-tab-pane label="位置 D" name="d">
          <PositionValueEditor position="d" :values="positions.d" @update="handleUpdate" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAllPositions } from '../api/api'
import PositionValueEditor from '../components/PositionValueEditor.vue'

const activeTab = ref('a')
const positions = reactive({
  a: [],
  b: [],
  c: [],
  d: []
})

// 加载所有位置值
const loadPositions = async () => {
  try {
    const data = await getAllPositions()
    if (data.positions) {
      positions.a = data.positions.a || []
      positions.b = data.positions.b || []
      positions.c = data.positions.c || []
      positions.d = data.positions.d || []
    }
  } catch (error) {
    ElMessage.error('加载配置失败: ' + error.message)
  }
}

// 刷新
const handleRefresh = () => {
  loadPositions()
  ElMessage.success('刷新成功')
}

// 标签切换
const handleTabChange = () => {
  // 可以在这里添加切换逻辑
}

// 更新位置值
const handleUpdate = async (position, values) => {
  positions[position] = values
  await loadPositions()
}

onMounted(() => {
  loadPositions()
})
</script>

<style scoped>
.position-config {
  max-width: 1200px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

