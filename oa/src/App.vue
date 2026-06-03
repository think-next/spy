<template>
  <el-container class="oa-layout">
    <el-aside width="220px" class="oa-aside">
      <div class="oa-logo">
        <el-icon :size="24"><Monitor /></el-icon>
        <span>OA 办公系统</span>
      </div>
      <el-menu
        :default-active="currentRoute"
        class="oa-menu"
        router
        background-color="transparent"
        text-color="rgba(255,255,255,0.7)"
        active-text-color="#409eff"
      >
        <el-menu-item index="/">
          <el-icon><DataAnalysis /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/announcement">
          <el-icon><Bell /></el-icon>
          <span>公告管理</span>
        </el-menu-item>
        <el-menu-item index="/leave">
          <el-icon><Calendar /></el-icon>
          <span>请假审批</span>
        </el-menu-item>
        <el-menu-item index="/contacts">
          <el-icon><User /></el-icon>
          <span>通讯录</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-main class="oa-main">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const currentRoute = computed(() => route.path)
</script>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  background: oklch(0.13 0.005 260);
  color: #e0e0e0;
}

/* 暗色主题 Element Plus 覆写 */
html.dark {
  --el-bg-color: oklch(0.13 0.005 260);
  --el-bg-color-overlay: oklch(0.18 0.005 260);
  --el-bg-color-page: oklch(0.13 0.005 260);
  --el-fill-color: oklch(0.18 0.005 260);
  --el-fill-color-light: oklch(0.20 0.005 260);
  --el-fill-color-lighter: oklch(0.22 0.005 260);
  --el-fill-color-extra-light: oklch(0.24 0.005 260);
  --el-border-color: oklch(0.25 0.005 260);
  --el-border-color-light: oklch(0.28 0.005 260);
  --el-border-color-lighter: oklch(0.30 0.005 260);
  --el-text-color-primary: #e0e0e0;
  --el-text-color-regular: #c0c0c0;
  --el-text-color-secondary: #909090;
  --el-card-bg-color: oklch(0.18 0.005 260);
}
html { color-scheme: dark; }

.oa-layout {
  height: 100vh;
}

.oa-aside {
  background: oklch(0.10 0.005 260);
  border-right: 1px solid oklch(0.25 0.005 260);
  display: flex;
  flex-direction: column;
}

.oa-logo {
  height: 60px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 20px;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  border-bottom: 1px solid oklch(0.25 0.005 260);
}

.oa-menu {
  border-right: none !important;
  flex: 1;
  padding-top: 10px;
}

.oa-menu .el-menu-item {
  margin: 2px 8px;
  border-radius: 6px;
}

.oa-menu .el-menu-item:hover {
  background: oklch(0.20 0.005 260) !important;
}

.oa-menu .el-menu-item.is-active {
  background: oklch(0.22 0.01 260) !important;
}

.oa-main {
  background: oklch(0.13 0.005 260);
  padding: 24px;
  overflow-y: auto;
}
</style>
