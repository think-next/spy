<template>
  <div class="dashboard">
    <h2 class="page-title">仪表盘</h2>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="8">
        <el-card shadow="never" class="stat-card" style="--delay: 0ms">
          <div class="stat-icon" style="background: linear-gradient(135deg, oklch(0.55 0.15 250 / 0.3), oklch(0.55 0.15 280 / 0.15)); color: oklch(0.75 0.15 250);">
            <el-icon :size="28"><Clock /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ pendingCount }}</div>
            <div class="stat-label">待审批</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never" class="stat-card" style="--delay: 100ms">
          <div class="stat-icon" style="background: linear-gradient(135deg, oklch(0.55 0.15 80 / 0.3), oklch(0.55 0.15 60 / 0.15)); color: oklch(0.80 0.12 80);">
            <el-icon :size="28"><Bell /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ announcements.length }}</div>
            <div class="stat-label">公告数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never" class="stat-card" style="--delay: 200ms">
          <div class="stat-icon" style="background: linear-gradient(135deg, oklch(0.55 0.12 150 / 0.3), oklch(0.55 0.12 170 / 0.15)); color: oklch(0.75 0.12 150);">
            <el-icon :size="28"><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ employees.length }}</div>
            <div class="stat-label">员工数</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 待办事项 -->
    <el-card class="todo-card">
      <template #header>
        <div class="card-header">
          <span>待办事项</span>
          <el-tag type="info" size="small" effect="plain">{{ pendingLeaves.length }} 项待处理</el-tag>
        </div>
      </template>
      <el-table :data="pendingLeaves" style="width: 100%">
        <el-table-column prop="applicant" label="申请人" width="100" />
        <el-table-column prop="department" label="部门" width="100" />
        <el-table-column prop="type" label="类型" width="80" />
        <el-table-column prop="startDate" label="开始日期" width="120" />
        <el-table-column prop="endDate" label="结束日期" width="120" />
        <el-table-column prop="days" label="天数" width="70" />
        <el-table-column prop="reason" label="事由" show-overflow-tooltip />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="scope">
            <el-button size="small" type="success" @click="handleAction(scope.row, 'approved')">通过</el-button>
            <el-button size="small" type="danger" @click="handleAction(scope.row, 'rejected')">驳回</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { announcements, leaveRequests, employees } from '../mock/data.js'

const leaves = ref([...leaveRequests])

const pendingLeaves = computed(() => leaves.value.filter(l => l.status === 'pending'))
const pendingCount = computed(() => pendingLeaves.value.length)

function handleAction(row, status) {
  row.status = status
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

.stat-cards {
  margin-bottom: 24px;
}

.stat-card {
  background: oklch(0.18 0.005 260 / 0.6);
  border: 1px solid oklch(0.25 0.005 260 / 0.5);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  animation: fadeInUp 0.5s ease-out backwards;
  animation-delay: var(--delay, 0ms);
}
.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px oklch(0 0 0 / 0.3), 0 0 0 1px oklch(0.30 0.01 260 / 0.5);
}

.stat-card :deep(.el-card__body) {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
}

.stat-label {
  font-size: 13px;
  color: #909090;
  margin-top: 2px;
}

.todo-card {
  background: oklch(0.18 0.005 260 / 0.6);
  border: 1px solid oklch(0.25 0.005 260 / 0.5);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  animation: fadeInUp 0.5s ease-out backwards;
  animation-delay: 300ms;
}
.todo-card :deep(.el-card__header) {
  border-bottom: 1px solid oklch(0.25 0.005 260);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #e0e0e0;
  font-weight: 600;
}
</style>
