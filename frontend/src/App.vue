<template>
  <el-container class="app-container" v-if="isAuthenticated">
    <el-header class="app-header">
      <h1>短信模板生成系统</h1>
      <div class="header-right">
        <el-menu
          mode="horizontal"
          :default-active="activeIndex"
          @select="handleMenuSelect"
          class="app-menu"
        >
          <el-menu-item index="1" @click="$router.push('/')">模板生成</el-menu-item>
          <el-menu-item index="2" @click="$router.push('/config')">位置配置</el-menu-item>
        </el-menu>
        <el-dropdown @command="handleCommand" class="user-dropdown">
          <span class="user-info">
            <el-icon><User /></el-icon>
            {{ username }}
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
  <router-view v-else />
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { User, ArrowDown } from '@element-plus/icons-vue'
import { useAuth } from './store/auth'

const route = useRoute()
const router = useRouter()
const { authState, clearAuth } = useAuth()

const activeIndex = ref('1')

const isAuthenticated = computed(() => authState.isAuthenticated)
const username = computed(() => authState.username)

watch(() => route.path, (path) => {
  if (path === '/config') {
    activeIndex.value = '2'
  } else if (path === '/') {
    activeIndex.value = '1'
  }
}, { immediate: true })

const handleMenuSelect = (index) => {
  activeIndex.value = index
}

const handleCommand = async (command) => {
  if (command === 'logout') {
    try {
      await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      clearAuth()
      router.push('/login')
    } catch {
      // 用户取消
    }
  }
}
</script>

<style scoped>
.app-container {
  min-height: 100vh;
}

.app-header {
  background-color: #409eff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.app-header h1 {
  margin: 0;
  font-size: 24px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-dropdown {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 5px;
  color: white;
  padding: 0 10px;
}

.app-menu {
  background-color: transparent;
  border: none;
}

.app-menu :deep(.el-menu-item) {
  color: white;
}

.app-menu :deep(.el-menu-item.is-active) {
  color: #409eff;
  background-color: white;
}

.el-main {
  padding: 20px;
  background-color: #f5f5f5;
}
</style>

