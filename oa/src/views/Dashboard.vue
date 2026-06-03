<template>
  <div class="dashboard">
    <h2 class="page-title">仪表盘</h2>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="8">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: #409eff20; color: #409eff;">
            <el-icon :size="28"><Clock /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ pendingCount }}</div>
            <div class="stat-label">待审批</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: #e6a23c20; color: #e6a23c;">
            <el-icon :size="28"><Bell /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ announcements.length }}</div>
            <div class="stat-label">公告数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-icon" style="background: #67c23a20; color: #67c23a;">
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
          <el-tag type="info" size="small">{{ pendingLeaves.length }} 项待处理</el-tag>
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
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 24px 0;
  color: #fff;
}

.stat-cards {
  margin-bottom: 24px;
}

.stat-card {
  background: oklch(0.18 0.005 260);
  border: 1px solid oklch(0.25 0.005 260);
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
  border-radius: 12px;
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
  background: oklch(0.18 0.005 260);
  border: 1px solid oklch(0.25 0.005 260);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #e0e0e0;
  font-weight: 600;
}
</style>
