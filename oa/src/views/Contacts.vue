<template>
  <div class="contacts">
    <h2 class="page-title">通讯录</h2>

    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <span>员工列表</span>
          <el-input
            v-model="keyword"
            placeholder="搜索姓名/部门..."
            style="width: 220px;"
            clearable
            class="search-input"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </template>
      <el-table :data="filteredEmployees" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="department" label="部门" width="120">
          <template #default="scope">
            <el-tag size="small" effect="plain" :class="'dept-tag dept-' + deptClass(scope.row.department)">{{ scope.row.department }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="position" label="职位" width="150" />
        <el-table-column prop="phone" label="手机" width="130" />
        <el-table-column prop="email" label="邮箱" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { employees } from '../mock/data.js'

const keyword = ref('')

const filteredEmployees = computed(() => {
  const kw = keyword.value.toLowerCase().trim()
  if (!kw) return employees
  return employees.filter(e =>
    e.name.toLowerCase().includes(kw) || e.department.toLowerCase().includes(kw)
  )
})

const deptColorMap = {
  '研发部': 'blue',
  '市场部': 'orange',
  '财务部': 'green',
  '人力资源部': 'pink',
  '行政部': 'cyan',
}

function deptClass(dept) {
  return deptColorMap[dept] || 'blue'
}
</script>

<style scoped>
.page-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 24px 0;
  color: #fff;
  position: relative;
  display: inline-block;
  padding-bottom: 8px;
}
.page-title::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 60px;
  height: 3px;
  border-radius: 2px;
  background: linear-gradient(90deg, oklch(0.75 0.15 250), oklch(0.65 0.18 270));
}

.main-card {
  background: oklch(0.18 0.005 260 / 0.6);
  border: 1px solid oklch(0.25 0.005 260 / 0.5);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  animation: fadeInUp 0.5s ease-out backwards;
  animation-delay: 0ms;
}
.main-card :deep(.el-card__header) {
  border-bottom: 1px solid oklch(0.25 0.005 260);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #e0e0e0;
  font-weight: 600;
}

/* 搜索框聚焦 accent 边框 */
.search-input :deep(.el-input__wrapper:focus-within) {
  box-shadow: 0 0 0 1px oklch(0.75 0.15 250) inset, 0 0 12px oklch(0.75 0.15 250 / 0.15) !important;
}

/* 部门标签颜色 */
.dept-tag {
  border-radius: 8px;
  font-weight: 500;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}
.dept-tag:hover {
  transform: scale(1.05);
}
.dept-blue {
  background: oklch(0.55 0.12 250 / 0.2) !important;
  color: oklch(0.75 0.15 250) !important;
  border-color: oklch(0.55 0.12 250 / 0.3) !important;
}
.dept-orange {
  background: oklch(0.55 0.12 60 / 0.2) !important;
  color: oklch(0.80 0.12 60) !important;
  border-color: oklch(0.55 0.12 60 / 0.3) !important;
}
.dept-green {
  background: oklch(0.55 0.10 150 / 0.2) !important;
  color: oklch(0.75 0.12 150) !important;
  border-color: oklch(0.55 0.10 150 / 0.3) !important;
}
.dept-pink {
  background: oklch(0.55 0.12 340 / 0.2) !important;
  color: oklch(0.80 0.12 340) !important;
  border-color: oklch(0.55 0.12 340 / 0.3) !important;
}
.dept-cyan {
  background: oklch(0.55 0.08 200 / 0.2) !important;
  color: oklch(0.75 0.10 200) !important;
  border-color: oklch(0.55 0.08 200 / 0.3) !important;
}
</style>
