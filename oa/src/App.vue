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
        active-text-color="oklch(0.75 0.15 250)"
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
      <div class="oa-user-info">
        <div class="user-avatar">
          <el-icon :size="18"><User /></el-icon>
        </div>
        <div class="user-detail">
          <span class="user-name">管理员</span>
          <span class="user-role">系统管理员</span>
        </div>
      </div>
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
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
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
  --el-color-primary: oklch(0.75 0.15 250);
  --el-color-primary-light-3: oklch(0.72 0.13 250);
  --el-color-primary-light-5: oklch(0.68 0.11 250);
  --el-color-primary-light-7: oklch(0.60 0.08 250);
  --el-color-primary-light-9: oklch(0.50 0.06 250);
  --el-color-primary-dark-2: oklch(0.65 0.14 250);
}
html { color-scheme: dark; }

/* 全局 Element Plus 组件覆写 */

/* 按钮 primary 渐变 */
.el-button--primary {
  background: linear-gradient(135deg, oklch(0.75 0.15 250), oklch(0.65 0.18 270)) !important;
  border-color: oklch(0.75 0.15 250) !important;
  color: #fff !important;
  font-weight: 500;
  transition: all 0.2s ease !important;
}
.el-button--primary:hover,
.el-button--primary:focus {
  background: linear-gradient(135deg, oklch(0.70 0.16 250), oklch(0.60 0.18 270)) !important;
  border-color: oklch(0.70 0.16 250) !important;
  box-shadow: 0 4px 15px oklch(0.65 0.15 250 / 0.3) !important;
}

/* 表格表头加深 */
.el-table th.el-table__cell {
  background: oklch(0.15 0.005 260) !important;
  color: #b0b0b0 !important;
  font-weight: 600;
  font-size: 13px;
  text-transform: uppercase;
  letter-spacing: 0.03em;
}
.el-table tr {
  transition: background 0.15s ease;
}
.el-table .el-table__row:hover > td.el-table__cell {
  background: oklch(0.22 0.01 260) !important;
}

/* 弹窗圆角加大 */
.el-dialog {
  border-radius: 12px !important;
  overflow: hidden;
}
.el-dialog .el-dialog__header {
  padding: 20px 24px 16px;
  border-bottom: 1px solid oklch(0.25 0.005 260);
}
.el-dialog .el-dialog__body {
  padding: 20px 24px;
}

/* 输入框聚焦 accent */
.el-input__wrapper:focus-within,
.el-textarea__inner:focus,
.el-input.is-focus .el-input__wrapper {
  box-shadow: 0 0 0 1px oklch(0.75 0.15 250) inset !important;
}

/* Tag 优化 */
.el-tag {
  border-radius: 8px;
  font-weight: 500;
  font-size: 12px;
  padding: 4px 12px;
  letter-spacing: 0.02em;
}
.el-tag--warning {
  background: oklch(0.25 0.08 80 / 0.5) !important;
  color: oklch(0.80 0.12 80) !important;
  border: 1px solid oklch(0.30 0.08 80 / 0.3) !important;
}
.el-tag--success {
  background: oklch(0.20 0.08 150 / 0.5) !important;
  color: oklch(0.75 0.12 150) !important;
  border: 1px solid oklch(0.25 0.08 150 / 0.3) !important;
}
.el-tag--danger {
  background: oklch(0.22 0.10 25 / 0.5) !important;
  color: oklch(0.75 0.15 25) !important;
  border: 1px solid oklch(0.28 0.10 25 / 0.3) !important;
}

/* 弹窗进入动画 */
.el-overlay-dialog {
  transition: opacity 0.25s ease;
}
.el-overlay-dialog .el-dialog {
  animation: dialogFadeIn 0.3s ease-out;
}

@keyframes dialogFadeIn {
  from {
    opacity: 0;
    transform: translateY(-12px) scale(0.97);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(15px); }
  to { opacity: 1; transform: translateY(0); }
}

/* --- Layout --- */

.oa-layout {
  height: 100vh;
}

.oa-aside {
  background: linear-gradient(180deg, oklch(0.12 0.005 260), oklch(0.10 0.005 260));
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
  background: linear-gradient(135deg, oklch(0.20 0.04 260), oklch(0.15 0.06 280));
  position: relative;
  overflow: hidden;
}
.oa-logo::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent, oklch(0.75 0.15 250 / 0.08), transparent);
  pointer-events: none;
}

.oa-menu {
  border-right: none !important;
  flex: 1;
  padding-top: 10px;
}

.oa-menu .el-menu-item {
  margin: 2px 8px;
  border-radius: 6px;
  position: relative;
  transition: all 0.15s ease;
}

.oa-menu .el-menu-item:hover {
  background: oklch(0.20 0.005 260) !important;
  color: #fff !important;
  transform: translateX(2px);
}

.oa-menu .el-menu-item.is-active {
  background: oklch(0.22 0.01 260) !important;
}

/* 左侧竖线 active 指示器 */
.oa-menu .el-menu-item.is-active::before {
  content: '';
  position: absolute;
  left: -8px;
  top: 50%;
  transform: translateY(-50%);
  width: 2px;
  height: 60%;
  background: oklch(0.75 0.15 250);
  border-radius: 0 2px 2px 0;
  box-shadow: 0 0 8px oklch(0.75 0.15 250 / 0.5);
}

/* 用户信息区域 */
.oa-user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 20px;
  border-top: 1px solid oklch(0.25 0.005 260);
  background: oklch(0.08 0.005 260);
}
.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, oklch(0.75 0.15 250), oklch(0.65 0.18 270));
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}
.user-detail {
  display: flex;
  flex-direction: column;
}
.user-name {
  color: #e0e0e0;
  font-size: 14px;
  font-weight: 600;
}
.user-role {
  color: #707070;
  font-size: 11px;
}

.oa-main {
  background: oklch(0.13 0.005 260);
  padding: 24px;
  overflow-y: auto;
}
</style>
